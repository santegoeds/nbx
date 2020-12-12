package api_v1

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/santegoeds/nbx/errors"
)

type Lifetime int

const (
	Minute Lifetime = 1
	Hour   Lifetime = Minute * 60
	Day    Lifetime = Hour * 24
	Week   Lifetime = Day * 7
)

type AuthenticateRequest struct {
	client     *Client
	AccountID  string
	KeyID      string
	Secret     string
	Passphrase string
	Lifetime   Lifetime
}

func NewAuthenticateRequest(c *Client, accountID, keyID, secret, passphrase string, lifetime Lifetime) *AuthenticateRequest {
	return &AuthenticateRequest{
		client:     c,
		AccountID:  accountID,
		KeyID:      keyID,
		Secret:     secret,
		Passphrase: passphrase,
		Lifetime:   lifetime,
	}
}

func (r *AuthenticateRequest) Do(ctx context.Context) error {
	path := "/accounts/" + r.AccountID + "/api_keys/" + r.KeyID + "/tokens"
	body := `{"expiresIn": ` + strconv.Itoa(int(r.Lifetime)) + "}"

	endpoint := r.client.Endpoint + path
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(body))
	if err != nil {
		log.Printf("Failed to create a request for Authenticate endpoint %s\n", endpoint)
		return err
	}
	if err = r.sign(req, body); err != nil {
		return err
	}

	rsp, err := r.client.HttpClient.Do(req)
	if err != nil {
		log.Printf("Request to endpoint %s failed\n", endpoint)
		return err
	}
	defer rsp.Body.Close()

	if err = errorFromResponse(rsp); err != nil {
		return err
	}

	var s struct {
		Token string `json:"token"`
	}
	dec := json.NewDecoder(rsp.Body)
	if err := dec.Decode(&s); err != nil {
		return err
	}

	if s.Token == "" {
		log.Println("Server error - empty token")
		return errors.ErrServer
	}
	r.client.AccountID = r.AccountID
	r.client.Token = s.Token
	return nil
}

func (r *AuthenticateRequest) sign(req *http.Request, body string) error {
	nowMillies := strconv.FormatInt(time.Now().UnixNano()/1000_000, 10)
	secret, err := base64.StdEncoding.DecodeString(r.Secret)
	if err != nil {
		return err
	}
	msg := []byte(nowMillies + req.Method + req.URL.Path + body)
	mac := hmac.New(sha256.New, secret)
	mac.Write(msg)
	signature := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	req.Header.Set("Authorization", "NBX-HMAC-SHA256 "+r.Passphrase+":"+signature)
	req.Header.Set("X-NBX-TIMESTAMP", nowMillies)
	return nil
}
