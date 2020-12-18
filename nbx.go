package nbx

import (
	"github.com/santegoeds/nbx/api_v1"
)

const (
	Minute = api_v1.Minute
	Hour   = api_v1.Hour
	Day    = api_v1.Day
	Week   = api_v1.Week
)

type Client = api_v1.Client
type Lifetime = api_v1.Lifetime

var (
	NewClient              = api_v1.NewClient
	NewTradeHistoryRequest = api_v1.NewTradeHistoryRequest
)
