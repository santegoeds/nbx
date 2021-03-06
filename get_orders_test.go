package nbx_test

import (
	"context"
	"github.com/santegoeds/nbx"
	"os"
	"testing"
	"time"

	"github.com/kr/pretty"
	"github.com/stretchr/testify/require"
)

func TestGetOrders(t *testing.T) {
	accountID := os.Getenv("ACCOUNT_ID")
	key := os.Getenv("KEY")
	secret := os.Getenv("SECRET")
	passphrase := os.Getenv("PASSPHRASE")

	client := nbx.NewClient()
	err := client.Authenticate(context.TODO(), accountID, key, secret, passphrase, nbx.Minute)
	require.NoError(t, err)

	orders, err := client.Orders(context.TODO())
	require.NoError(t, err)
	require.NotEmpty(t, orders)

	pretty.Println(orders)

	now := time.Now()
	for _, order := range orders {
		require.NotEmpty(t, order.ID)
		require.NotEmpty(t, order.Market)
		require.Greater(t, order.Quantity, 0.0)
		require.Contains(t, []string{"BUY", "SELL"}, order.Side)
		require.True(t, order.Events.CreatedAt.Before(now) || order.Events.CreatedAt.Equal(now))
	}
}
