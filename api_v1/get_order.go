package api_v1

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type Order struct {
	ID        string         `json:"id"`
	Events    OrderEvents    `json:"events"`
	Fills     []OrderFill    `json:"fills"`
	Market    string         `json:"market"`
	Quantity  float64        `json:"quantity,string"`
	Side      string         `json:"side"`
	Execution OrderExecution `json:"executions"`
}

type OrderEvents struct {
	CreatedAt  time.Time  `json:"createdAt"`
	OpenedAt   *time.Time `json:"openedAt"`
	ClosedAt   *time.Time `json:"closedAt"`
	RejectedAt *time.Time `json:"rejectedAt"`
}

type OrderFill struct {
	Quantity  float64   `json:"quantity,string"`
	Price     float64   `json:"price,string"`
	Fee       float64   `json:"fee,string"`
	CreatedAt time.Time `json:"createdAt"`
}

type OrderExecution struct {
	Type        string `json:"type"`
	Price       string `json:"price,string"`
	TimeInForce struct {
		Type string `json:"type"`
	} `json:"timeInForce"`
}

type GetOrderRequest struct {
	client    *Client
	accountID string
	orderID   string
}

func NewGetOrderRequest(client *Client, accountID, orderID string) *GetOrderRequest {
	return &GetOrderRequest{
		client:    client,
		accountID: accountID,
		orderID:   orderID,
	}
}

func (r *GetOrderRequest) Do(ctx context.Context) (*Order, error) {
	endpoint := r.client.Endpoint + "/accounts/" + r.accountID + "/orders/" + r.orderID
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		log.Printf("Failed to create request for endpoint %s\n", endpoint)
		return nil, err
	}
	rsp, err := r.client.do(req)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	if err = errorFromResponse(rsp); err != nil {
		return nil, err
	}
	dec := json.NewDecoder(rsp.Body)
	order := &Order{}
	if err = dec.Decode(order); err != nil {
		return nil, err
	}
	return order, nil
}
