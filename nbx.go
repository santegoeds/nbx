package nbx

import (
	"github.com/santegoeds/nbx/api_v1"
)

type Client = api_v1.Client
type Lifetime = api_v1.Lifetime

const (
	Minute = api_v1.Minute
	Hour   = api_v1.Hour
	Day    = api_v1.Day
	Week   = api_v1.Week
)
