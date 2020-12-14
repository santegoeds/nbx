package nbx

import (
	"github.com/santegoeds/nbx/api_v1"
)

type Client = api_v1.Client
type Lifetime = api_v1.Lifetime
type AuthenticateRequest = api_v1.AuthenticateRequest
type CancelOrderRequest = api_v1.CancelOrderRequest
type CreateOrderRequest = api_v1.CreateOrderRequest
type GetOrderRequest = api_v1.GetOrderRequest
type GetOrdersRequest = api_v1.GetOrdersRequest
type OrderbookRequest = api_v1.OrderbookRequest
type TradeHistoryRequest = api_v1.TradeHistoryRequest

const (
	Minute = api_v1.Minute
	Hour   = api_v1.Hour
	Day    = api_v1.Day
	Week   = api_v1.Week
)
