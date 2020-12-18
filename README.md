# NBX SDK for Go

NBX SDK for Go is a library to access the REST API of Norwegian Block Exchange using the Go programming language.

Norwegian Block Exchange is a pioneering, forward-thinking, client-focused cryptocurrency exchange, custodian and payment system.  To find out more see https://nbx.com.

This SDK is not authorized or endorsed by Norwegian Block Exchange.  Note that this SDK this is an alpha release.  The code is likely to contain bugs, so use at your own risk.

## Example

```go
package main

import (
	"context"
	"os"

	"github.com/kr/pretty"
	"github.com/santegoeds/nbx"
)

const market = "BTC-NOK"

func main() {
	ctx := context.Background()

	// Create a new client
	client := nbx.NewClient()

	// Fetch historic trades (no authentication needed)
	trades, _ := client.TradeHistory(ctx, market)
	pretty.Println(trades)

	// Paginate through all historic trades (no authentication needed)
	req := nbx.NewTradeHistoryRequest(client, market)
	trades, _ = req.Do(ctx)
	for req.HasNextPage() {
		trades, _ = req.Do(ctx)
	}

	// Fetch the orderbook (no authentication needed)
	orderbook, _ := client.Orderbook(ctx, market)
	pretty.Println(orderbook)

	// Read authentication details from environment
	accountID := os.Getenv("ACCOUNT_ID")
	keyID := os.Getenv("KEY")
	secret := os.Getenv("SECRET")
	passphrase := os.Getenv("PASSPHRASE")

	// Authenticate and request a token with a lifetime of one minute.
	_ = client.Authenticate(ctx, accountID, keyID, secret, passphrase, nbx.Minute)

	// List Market information
	markets, _ := client.Markets(ctx)
	pretty.Println(markets)

	// Create a market buy order to buy 0.00001 BTC for no more than 10 NOK
	orderID, _ := client.MarketBuy(ctx, market, 0.00001, 10.0)

	// Get the order details
	order, _ := client.GetOrder(ctx, orderID)
	pretty.Println(order)

	// Create a limit order to sell 0.00001 BTC at a price of 200,000 NOK
	orderID, _ = client.LimitSell(ctx, market, 200_000.0, 0.00001)

	// Get all orders for the account
	orders, _ := client.Orders(ctx)
	pretty.Println(orders)

	// Cancel limit order
	_ = client.CancelOrder(ctx, orderID)
}
```

## License

NBX SDK for Go is released under the MIT license
