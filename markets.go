package nbx

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/santegoeds/nbx/errors"
)

type Market struct {
	MinimumQuantity float64 `json:"min-qty,string"`
	MinimumValue    float64 `json:"min-val,string"`
	MakerFee        float64 `json:"maker-fee,string"`
	TakerFee        float64 `json:"taker-fee,string"`
	MinimumTickSize float64 `json:"min-tick,string"`
}

type MarketsRequest struct {
	client *Client
}

func NewMarketsRequest(client *Client) *MarketsRequest {
	return &MarketsRequest{
		client: client,
	}
}

func (r *MarketsRequest) Do(ctx context.Context) (map[string]Market, error) {
	if len(r.client.markets) == 0 {
		if r.client.Token == "" {
			return nil, fmt.Errorf("client is not authenticated: %w", errors.ErrBadRequest)
		}
		claims, err := decodeClaims(r.client.Token)
		if err != nil {
			return nil, err
		}
		r.client.markets = claims.Markets
	}
	return r.client.markets, nil
}

type claims struct {
	Markets map[string]Market
}

func decodeClaims(token string) (*claims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("%w: invalid JWT token %s", errors.ErrBadRequest, token)
	}
	claimsStr := parts[1]
	paddingLen := 4 - (len(claimsStr) % 4)
	if paddingLen > 0 && paddingLen < 4 {
		claimsStr += strings.Repeat("=", paddingLen)
	}
	data, err := base64.StdEncoding.DecodeString(claimsStr)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid JWT token %s", errors.ErrBadRequest, token)
	}

	claims := claims{}
	if err := json.Unmarshal(data, &claims); err != nil {
		return nil, fmt.Errorf("%w: invalid JWT token %s", errors.ErrBadRequest, token)
	}

	return &claims, nil
}
