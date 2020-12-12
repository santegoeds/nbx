package api_v1

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/santegoeds/nbx/errors"
)

const (
	orderFlagGTC    = "GOOD_TIL_CANCELED"
	orderFlagIOC    = "IMMEDIATE_OR_CANCEL"
	orderSideBuy    = "BUY"
	orderSideSell   = "SELL"
	orderTypeLimit  = "LIMIT"
	orderTypeMarket = "MARKET"
)

type CreateOrderRequest struct {
	client          *Client
	accountID       string
	marketID        string
	side            string
	quantity        string
	orderType       string
	price           string
	freezeAmount    string
	timeInForceType string
}

func NewCreateOrderRequest(client *Client, accountID, marketID string) *CreateOrderRequest {
	return &CreateOrderRequest{
		client:    client,
		accountID: accountID,
		marketID:  marketID,
	}
}

func (r *CreateOrderRequest) SetLimitBuy(quantity float64, price float64) *CreateOrderRequest {
	r.side = orderSideBuy
	r.orderType = orderTypeLimit
	r.timeInForceType = orderFlagGTC
	r.price = formatFloat(price)
	r.quantity = formatFloat(quantity)
	return r
}

func (r *CreateOrderRequest) SetLimitSell(quantity float64, price float64) *CreateOrderRequest {
	r.side = orderSideSell
	r.orderType = orderTypeLimit
	r.timeInForceType = orderFlagGTC
	r.price = formatFloat(price)
	r.quantity = formatFloat(quantity)
	return r
}

func (r *CreateOrderRequest) SetMarketBuy(quantity float64, amount float64) *CreateOrderRequest {
	r.side = orderSideBuy
	r.orderType = orderTypeMarket
	r.timeInForceType = orderFlagIOC
	r.quantity = formatFloat(quantity)
	r.freezeAmount = formatFloat(amount)
	return r
}

func (r *CreateOrderRequest) SetMarketSell(quantity float64) *CreateOrderRequest {
	r.side = orderSideSell
	r.orderType = orderTypeMarket
	r.timeInForceType = orderFlagIOC
	r.quantity = formatFloat(quantity)
	return r
}

func (r *CreateOrderRequest) Do(ctx context.Context) (string, error) {
	endpoint := r.client.Endpoint + "/accounts/" + r.accountID + "/orders"

	body := creatOrderPayload{
		Market:   r.marketID,
		Quantity: r.quantity,
		Side:     r.side,
		Execution: payloadExecution{
			Type: r.orderType,
		},
	}
	switch r.orderType {
	case "MARKET":
		body.Execution.TimeInForce.Type = "IMMEDIATE_OR_CANCEL"
		if r.side == "BUY" {
			body.Execution.Freeze = &payloadFreeze{
				Type:  "AMOUNT",
				Value: r.freezeAmount,
			}
		}
	case "LIMIT":
		body.Execution.TimeInForce.Type = "GOOD_TIL_CANCELED"
		body.Execution.Price = &r.price
	}

	data, err := json.Marshal(body)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewBuffer(data))
	if err != nil {
		log.Printf("Failed to create request for endpoint %s\n", endpoint)
		return "", err
	}
	rsp, err := r.client.do(req)
	if err != nil {
		return "", err
	}
	defer rsp.Body.Close()

	if err = errorFromResponse(rsp); err != nil {
		return "", err
	}
	location := rsp.Header.Get("location")
	if location == "" {
		return "", errors.ErrServer
	}
	var parts []string
	if location != "" {
		parts = strings.Split(location, "/")
		if len(parts) == 0 {
			return "", errors.ErrServer
		}
	}
	return parts[len(parts)-1], nil
}

type creatOrderPayload struct {
	Market    string           `json:"market"`
	Quantity  string           `json:"quantity"`
	Side      string           `json:"side"`
	Execution payloadExecution `json:"execution"`
}

type payloadExecution struct {
	Type        string             `json:"type"`
	TimeInForce payloadTimeInForce `json:"timeInForce"`
	Price       *string            `json:"price,omitemtpy"`
	Freeze      *payloadFreeze     `json:"freeze,omitempty"`
}

type payloadTimeInForce struct {
	Type string `json:"type"`
}

type payloadFreeze struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

func formatFloat(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}
