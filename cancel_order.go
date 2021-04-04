package nbx

import (
	"context"
	"log"
	"net/http"
)

type CancelOrderRequest struct {
	client    *Client
	accountID string
	orderID   string
}

func NewCancelOrder(client *Client, accountID, orderID string) *CancelOrderRequest {
	return &CancelOrderRequest{
		client:    client,
		accountID: accountID,
		orderID:   orderID,
	}
}

func (r *CancelOrderRequest) Do(ctx context.Context) error {
	endpoint := r.client.Endpoint + "/accounts/" + r.accountID + "/orders/" + r.orderID
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, endpoint, nil)
	if err != nil {
		log.Printf("Failed to create request for endpoint %s\n", endpoint)
		return err
	}
	rsp, err := r.client.do(req)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()
	return errorFromResponse(rsp)
}
