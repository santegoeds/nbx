package api_v1

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
)

type Balance struct {
	Available float64 `json:"available,string"`
	Total     float64 `json:"total,string"`
}

type Freeze struct {
	ID     string  `json:"id"`
	Amount float64 `json:"amount,string"`
}

type Asset struct {
	Balance Balance
	Freezes []Freeze
}

type AssetsRequest struct {
	client    *Client
	accountID string
}

func NewAssetsRequest(client *Client) *AssetsRequest {
	return &AssetsRequest{
		client:    client,
		accountID: client.AccountID,
	}
}

func (r *AssetsRequest) Do(ctx context.Context) (map[string]Asset, error) {
	endpoint := r.client.Endpoint + "/accounts/" + r.accountID + "/assets"
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

	type rawAsset struct {
		ID string `json:"id"`
		Asset
	}
	rawAssets := make([]rawAsset, 0)
	if err = dec.Decode(&rawAssets); err != nil {
		return nil, err
	}
	assets := make(map[string]Asset, len(rawAssets))
	for _, asset := range rawAssets {
		assets[asset.ID] = asset.Asset
	}
	return assets, nil
}
