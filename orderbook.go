package nbx

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type OrderbookRequest struct {
	client *Client
	market string
	Side   string
}

type OrderbookOrder struct {
	ID       string  `json:"id"`
	Price    float64 `json:"price,string"`
	Quantity float64 `json:"quantity,string"`
	Side     string  `json:"side"`
}

type Orderbook struct {
	Buys  []OrderbookOrder
	Sells []OrderbookOrder
}

func NewOrderbookRequest(client *Client, market string) *OrderbookRequest {
	return &OrderbookRequest{
		client: client,
		market: market,
	}
}

func (r *OrderbookRequest) Do(ctx context.Context) (*Orderbook, error) {
	endpoint := r.client.Endpoint + "/markets/" + r.market + "/orders"
	if r.Side != "" {
		endpoint += "?side=" + r.Side
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		log.Printf("Failed to create request for endpoint %s\n", endpoint)
		return nil, err
	}

	rsp, err := r.client.HttpClient.Do(req)
	if err != nil {
		log.Printf("Request for Orderbook endpoint %s failed\n", endpoint)
		return nil, err
	}
	defer rsp.Body.Close()

	if err = errorFromResponse(rsp); err != nil {
		return nil, err
	}

	dec := json.NewDecoder(rsp.Body)
	rawOrders := make([]OrderbookOrder, 0, 100)
	if err = dec.Decode(&rawOrders); err != nil {
		return nil, err
	}

	book := &Orderbook{}
	switch strings.ToUpper(r.Side) {
	case "BUY":
		book.Buys = rawOrders
	case "SELL":
		book.Sells = rawOrders
	default:
		for _, order := range rawOrders {
			if order.Side == "SELL" {
				book.Sells = append(book.Sells, order)
			} else {
				book.Buys = append(book.Buys, order)
			}
		}
	}
	return book, nil
}
