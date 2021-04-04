package nbx_test

import (
	"context"
	"github.com/santegoeds/nbx"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTradeHistory(t *testing.T) {
	client := nbx.NewClient()
	market := "BTC-NOK"

	require.True(t, t.Run("first page without pagination", func(t *testing.T) {
		trades, err := client.TradeHistory(context.TODO(), market)
		require.NoError(t, err)
		require.NotEmpty(t, trades)

		createdAt := trades[0].CreatedAt
		for _, trade := range trades {
			afterOrEqual := createdAt.After(trade.CreatedAt) || createdAt.Equal(trade.CreatedAt)
			require.True(t, afterOrEqual)
			require.Greater(t, trade.Price, 0.0)
			require.Greater(t, trade.Quantity, 0.0)
			require.NotEmpty(t, trade.MakerOrder.ID)
			require.NotEmpty(t, trade.MakerOrder.Account.ID)
			require.NotEmpty(t, trade.TakerOrder.ID)
			require.NotEmpty(t, trade.TakerOrder.Account.ID)

			sides := []string{trade.MakerOrder.Side, trade.TakerOrder.Side}
			require.Contains(t, sides, "BUY")
			require.Contains(t, sides, "SELL")
		}
	}))

	require.True(t, t.Run("request with pagination", func(t *testing.T) {
		req := nbx.NewTradeHistoryRequest(client, market)
		firstPage, err := req.Do(context.TODO())
		require.NoError(t, err)
		require.NotEmpty(t, firstPage)

		secondPage, err := req.Do(context.TODO())
		require.NoError(t, err)
		require.NotEmpty(t, secondPage)

		createdAt := firstPage[len(firstPage)-1].CreatedAt
		for _, trade := range secondPage {
			afterOrEqual := createdAt.After(trade.CreatedAt) || createdAt.Equal(trade.CreatedAt)
			require.True(t, afterOrEqual)
		}
	}))
}
