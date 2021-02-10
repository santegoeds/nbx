package api_v1

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
)

type AssetRequest struct {
	client    *Client
	accountID string
	asset     string
}

func NewAssetRequest(client *Client, asset string) *AssetRequest {
	return &AssetRequest{
		client:    client,
		accountID: client.AccountID,
		asset:     asset,
	}
}

func (r *AssetRequest) Do(ctx context.Context) (*Asset, error) {
	endpoint := r.client.Endpoint + "/accounts/" + r.accountID + "/assets/" + r.asset
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

	asset := Asset{}
	if err = dec.Decode(&asset); err != nil {
		return nil, err
	}
	return &asset, nil
}
