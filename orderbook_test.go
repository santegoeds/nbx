package nbx_test

import (
	"context"
	"github.com/santegoeds/nbx"
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOrderbook(t *testing.T) {
	client := nbx.NewClient()
	market := "BTC-NOK"

	require.True(t, t.Run("should return a full orderbook", func(t *testing.T) {
		orderbook, err := client.Orderbook(context.TODO(), market)
		require.NoError(t, err)
		require.NotEmpty(t, orderbook.Buys)
		require.NotEmpty(t, orderbook.Sells)

		// Assert that the orders are sorted by price
		lastPrice := math.Inf(0)
		for _, order := range orderbook.Buys {
			require.Equal(t, "BUY", order.Side)
			require.LessOrEqual(t, order.Price, lastPrice)
			require.Greater(t, order.Price, 0.0)
			require.Greater(t, order.Quantity, 0.0)
		}
		lastPrice = math.Inf(-1)
		for _, order := range orderbook.Sells {
			require.Equal(t, "SELL", order.Side)
			require.Greater(t, order.Price, lastPrice)
			require.Greater(t, order.Price, 0.0)
			require.Greater(t, order.Quantity, 0.0)
		}
	}))

	require.True(t, t.Run("should return an orderbook with only BUY orders", func(t *testing.T) {
		req := nbx.NewOrderbookRequest(client, market)
		req.Side = "BUY"

		orderbook, err := req.Do(context.TODO())
		require.NoError(t, err)
		require.NotEmpty(t, orderbook.Buys)
		require.Empty(t, orderbook.Sells)

		lastPrice := math.Inf(0)
		for _, order := range orderbook.Buys {
			require.Equal(t, "BUY", order.Side)
			require.LessOrEqual(t, order.Price, lastPrice)
		}
	}))

	require.True(t, t.Run("should return an orderbook with only SELL orders", func(t *testing.T) {
		req := nbx.NewOrderbookRequest(client, market)
		req.Side = "SELL"

		orderbook, err := req.Do(context.TODO())
		require.NoError(t, err)
		require.NotEmpty(t, orderbook.Sells)
		require.Empty(t, orderbook.Buys)

		lastPrice := math.Inf(-1)
		for _, order := range orderbook.Sells {
			require.Equal(t, "SELL", order.Side)
			require.GreaterOrEqual(t, order.Price, lastPrice)
		}
	}))
}
