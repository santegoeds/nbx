package nbx_test

import (
	"context"
	"github.com/santegoeds/nbx"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateOrder(t *testing.T) {
	accountID := os.Getenv("ACCOUNT_ID")
	key := os.Getenv("KEY")
	secret := os.Getenv("SECRET")
	passphrase := os.Getenv("PASSPHRASE")
	market := "BTC-NOK"

	client := nbx.NewClient()
	err := client.Authenticate(context.TODO(), accountID, key, secret, passphrase, nbx.Minute)
	require.NoError(t, err)

	require.True(t, t.Run("create market buy order", func(t *testing.T) {
		orderID, err := client.MarketBuy(context.TODO(), market, 0.000001, 0.01)
		require.NoError(t, err)
		require.NotEmpty(t, orderID)
	}))

	require.True(t, t.Run("create market sell order", func(t *testing.T) {
		orderID, err := client.MarketSell(context.TODO(), market, 0.000001)
		require.NoError(t, err)
		require.NotEmpty(t, orderID)
	}))

	require.True(t, t.Run("create limit buy order", func(t *testing.T) {
		orderID, err := client.LimitBuy(context.TODO(), market, 10.0, 1)
		require.NoError(t, err)
		require.NotEmpty(t, orderID)
	}))

	require.True(t, t.Run("create limit sell order", func(t *testing.T) {
		orderbook, err := client.Orderbook(context.TODO(), market)
		require.NoError(t, err)

		price := orderbook.Sells[len(orderbook.Sells)-1].Price * 2
		orderID, err := client.LimitSell(context.TODO(), market, price, 0.000001)
		require.NoError(t, err)
		require.NotEmpty(t, orderID)
	}))
}
