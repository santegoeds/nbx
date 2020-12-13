package api_v1

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
)

type GetOrdersRequest struct {
	client    *Client
	accountID string
}

func NewGetOrdersRequest(client *Client, accountID string) *GetOrdersRequest {
	return &GetOrdersRequest{
		client:    client,
		accountID: accountID,
	}
}

func (r *GetOrdersRequest) Do(ctx context.Context) ([]Order, error) {
	endpoint := r.client.Endpoint + "/accounts/" + r.accountID + "/orders"
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
	orders := make([]Order, 0)
	if err = dec.Decode(&orders); err != nil {
		return nil, err
	}
	return orders, nil
}
