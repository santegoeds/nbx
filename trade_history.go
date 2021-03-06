package nbx

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type TradeHistoryRequest struct {
	client      *Client
	market      string
	nextPage string
}

func NewTradeHistoryRequest(client *Client, market string) *TradeHistoryRequest {
	return &TradeHistoryRequest{
		client: client,
		market: market,
	}
}

func (r *TradeHistoryRequest) HasNextPage() bool {
	return r.nextPage != ""
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
	CreatedAt  time.Time     `json:"createdAt"`
	MakerOrder HistoricOrder `json:"makerOrder"`
	Price      float64       `json:"price,string"`
	Quantity   float64       `json:"quantity,string"`
	TakerOrder HistoricOrder `json:"takerOrder"`
}

func (r *TradeHistoryRequest) Do(ctx context.Context) ([]HistoricTrade, error) {
	var endpoint string
	if r.nextPage != "" {
		endpoint = r.nextPage
	} else {
		endpoint = r.client.Endpoint + "/markets/" + r.market + "/trades"
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		log.Printf("Failed to create request for endpoint %s\n", endpoint)
		return nil, err
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
	r.nextPage = rsp.Header.Get("x-next-page-url")

	dec := json.NewDecoder(rsp.Body)
	trades := make([]HistoricTrade, 0, 100)
	if err = dec.Decode(&trades); err != nil {
		return nil, err
	}
	return trades, nil
}
