package nbx_test

import (
	"context"
	"github.com/santegoeds/nbx"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestGetOrder(t *testing.T) {
	accountID := os.Getenv("ACCOUNT_ID")
	key := os.Getenv("KEY")
	secret := os.Getenv("SECRET")
	passphrase := os.Getenv("PASSPHRASE")
	market := "BTC-NOK"

	client := nbx.NewClient()
	err := client.Authenticate(context.TODO(), accountID, key, secret, passphrase, nbx.Minute)
	require.NoError(t, err)

	orderID, err := client.MarketBuy(context.TODO(), market, 0.000001, 0.01)
	require.NoError(t, err)
	require.NotEmpty(t, orderID)

	order, err := client.GetOrder(context.TODO(), orderID)
	require.NoError(t, err)

	require.Equal(t, orderID, order.ID)
	require.Equal(t, market, order.Market)
	require.Greater(t, order.Quantity, 0.0)
	require.Equal(t, "BUY", order.Side)
	now := time.Now()
	require.True(t, order.Events.CreatedAt.Before(now) || order.Events.CreatedAt.Equal(now))
}
