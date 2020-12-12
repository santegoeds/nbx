package api_v1_test

import (
	"context"
	"os"
	"testing"

	"github.com/santegoeds/nbx/api_v1"
	"github.com/stretchr/testify/require"
)

func TestCreateOrder(t *testing.T) {
	accountID := os.Getenv("ACCOUNT_ID")
	key := os.Getenv("KEY")
	secret := os.Getenv("SECRET")
	passphrase := os.Getenv("PASSPHRASE")
	marketID := "BTC-NOK"

	client := api_v1.NewClient()
	err := client.Authenticate(context.TODO(), accountID, key, secret, passphrase, api_v1.Minute)
	require.NoError(t, err)

	require.True(t, t.Run("create market buy order", func(t *testing.T) {
		orderID, err := client.MarketBuy(context.TODO(), accountID, marketID, 0.000001, 0.01)
		require.NoError(t, err)
		require.NotEmpty(t, orderID)
	}))

	require.True(t, t.Run("create market sell order", func(t *testing.T) {
		orderID, err := client.MarketSell(context.TODO(), accountID, marketID, 0.000001)
		require.NoError(t, err)
		require.NotEmpty(t, orderID)
	}))

	require.True(t, t.Run("create limit buy order", func(t *testing.T) {
		orderID, err := client.LimitBuy(context.TODO(), accountID, marketID, 10.0, 1)
		require.NoError(t, err)
		require.NotEmpty(t, orderID)
	}))

	require.True(t, t.Run("create limit sell order", func(t *testing.T) {
		orderbook, err := client.Orderbook(context.TODO(), marketID)
		require.NoError(t, err)

		price := orderbook.Sells[len(orderbook.Sells)-1].Price * 2
		orderID, err := client.LimitSell(context.TODO(), accountID, marketID, price, 0.000001)
		require.NoError(t, err)
		require.NotEmpty(t, orderID)
	}))
}
