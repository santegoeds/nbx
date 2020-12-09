package api_v1

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type TradeHistoryRequest struct {
	client      *Client
	marketID    string
	pagingState string
}

func NewTradeHistoryRequest(client *Client, marketID string) *TradeHistoryRequest {
	return &TradeHistoryRequest{
		client:   client,
		marketID: marketID,
	}
}

func (r *TradeHistoryRequest) HasNextPage() bool {
	return r.pagingState != ""
}

type HistoricTradeAccount struct {
	ID string `json:"id"`
}

type HistoricOrder struct {
	Account HistoricTradeAccount `json:"account"`
	ID      string               `json:"id"`
	Side    string               `json:"side"`
}

type HistoricTrade struct {
	CreatedAt  time.Time     `json:"createdAt,string"`
	MakerOrder HistoricOrder `json:"makerOrder"`
	Price      float64       `json:"price,string"`
	Quantity   float64       `json:"quantity,string"`
	TakerOrder HistoricOrder `json:"takerOrder"`
}

func (r *TradeHistoryRequest) Do(ctx context.Context) ([]HistoricTrade, error) {
	endpoint := r.client.Endpoint + "/markets/" + r.marketID + "/trades"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		log.Printf("Failed to create request for endpoint %s\n", endpoint)
		return nil, err
	}
	if r.pagingState != "" {
		req.Header.Set("x-paging-state", r.pagingState)
	}

	rsp, err := r.client.HttpClient.Do(req)
	if err != nil {
		log.Printf("Request for TradeHistory endpoint %s failed\n", endpoint)
		return nil, err
	}
	defer rsp.Body.Close()

	if err = errorFromResponse(rsp); err != nil {
		return nil, err
	}

	// Update the paging state so that the same request can be used to get the next page.
	r.pagingState = rsp.Header.Get("x-paging-state")

	dec := json.NewDecoder(rsp.Body)
	trades := make([]HistoricTrade, 0, 100)
	if err = dec.Decode(&trades); err != nil {
		return nil, err
	}
	return trades, nil
}
