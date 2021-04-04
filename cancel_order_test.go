package nbx_test

import (
	"context"
	"github.com/santegoeds/nbx"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCancelOrder(t *testing.T) {
	t.Skipf("Skipping while account is not funded")

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

	err = client.CancelOrder(context.TODO(), orderID)
	require.NoError(t, err)
}
