package api_v1_test

import (
	"context"
	"os"
	"testing"

	"github.com/santegoeds/nbx/api_v1"
	"github.com/stretchr/testify/require"
)

func TestCancelOrder(t *testing.T) {
	t.Skipf("Skipping while account is not funded")

	accountID := os.Getenv("ACCOUNT_ID")
	key := os.Getenv("KEY")
	secret := os.Getenv("SECRET")
	passphrase := os.Getenv("PASSPHRASE")
	marketID := "BTC-NOK"

	client := api_v1.NewClient()
	err := client.Authenticate(context.TODO(), accountID, key, secret, passphrase, api_v1.Minute)
	require.NoError(t, err)

	orderID, err := client.MarketBuy(context.TODO(), marketID, 0.000001, 0.01)
	require.NoError(t, err)
	require.NotEmpty(t, orderID)

	err = client.CancelOrder(context.TODO(), orderID)
	require.NoError(t, err)
}
