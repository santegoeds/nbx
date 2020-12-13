package api_v1

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/santegoeds/nbx/errors"
)

type Client struct {
	Endpoint   string
	HttpClient *http.Client
	Token      string
	AccountID  string
}

func NewClient() *Client {
	client := &Client{
		Endpoint:   "https://api.nbx.com",
		HttpClient: http.DefaultClient,
	}
	return client
}

func (c *Client) Authenticate(
	ctx context.Context,
	accountID string,
	keyID string,
	secret string,
	passphrase string,
	lifetime Lifetime,
) error {
	req := NewAuthenticateRequest(c, accountID, keyID, secret, passphrase, lifetime)

	return req.Do(ctx)
}

func (c *Client) Orderbook(ctx context.Context, marketID string) (*Orderbook, error) {
	req := NewOrderbookRequest(c, marketID)
	return req.Do(ctx)
}

func (c *Client) TradeHistory(ctx context.Context, marketID string) ([]HistoricTrade, error) {
	req := NewTradeHistoryRequest(c, marketID)
	return req.Do(ctx)
}

func (c *Client) LimitBuy(
	ctx context.Context,
	marketID string,
	price float64,
	quantity float64,
) (string, error) {
	req := NewCreateOrderRequest(c, c.AccountID, marketID)
	return req.SetLimitBuy(quantity, price).Do(ctx)
}

func (c *Client) LimitSell(
	ctx context.Context,
	marketID string,
	price float64,
	quantity float64,
) (string, error) {
	req := NewCreateOrderRequest(c, c.AccountID, marketID)
	return req.SetLimitSell(quantity, price).Do(ctx)
}

func (c *Client) MarketBuy(
	ctx context.Context,
	marketID string,
	quantity float64,
	amount float64,
) (string, error) {
	req := NewCreateOrderRequest(c, c.AccountID, marketID)
	return req.SetMarketBuy(quantity, amount).Do(ctx)
}

func (c *Client) MarketSell(
	ctx context.Context,
	marketID string,
	quantity float64,
) (string, error) {
	req := NewCreateOrderRequest(c, c.AccountID, marketID)
	return req.SetMarketSell(quantity).Do(ctx)
}

func (c *Client) GetOrder(ctx context.Context, orderID string) (*Order, error) {
	req := NewGetOrderRequest(c, c.AccountID, orderID)
	return req.Do(ctx)
}

func (c *Client) CancelOrder(ctx context.Context, orderID string) error {
	req := NewCancelOrder(c, c.AccountID, orderID)
	return req.Do(ctx)
}

func (c *Client) GetOrders(ctx context.Context) ([]Order, error) {
	req := NewGetOrdersRequest(c, c.AccountID)
	return req.Do(ctx)
}

func (c *Client) do(req *http.Request) (*http.Response, error) {
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}
	return c.HttpClient.Do(req)
}

func errorFromResponse(rsp *http.Response) error {
	switch rsp.StatusCode / 100 {
	case 5:
		body, _ := ioutil.ReadAll(rsp.Body)
		log.Printf("%s - Server Error - %s\n", rsp.Request.URL, body)
		return errors.ErrServer

	case 4:
		body, _ := ioutil.ReadAll(rsp.Body)
		log.Printf("%s - Bad Request - %s\n", rsp.Request.URL, body)
		return errors.ErrBadRequest
	}
	return nil
}
