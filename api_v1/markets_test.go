package api_v1_test

import (
	"context"
	"os"
	"testing"

	"github.com/kr/pretty"
	"github.com/santegoeds/nbx/api_v1"

	"github.com/stretchr/testify/require"
)

func TestMarkets(t *testing.T) {
	accountID := os.Getenv("ACCOUNT_ID")
	key := os.Getenv("KEY")
	secret := os.Getenv("SECRET")
	passphrase := os.Getenv("PASSPHRASE")

	client := api_v1.NewClient()
	err := client.Authenticate(context.TODO(), accountID, key, secret, passphrase, api_v1.Minute)
	require.NoError(t, err)

	markets, err := client.Markets(context.TODO())
	require.NoError(t, err)

	require.NotEmpty(t, markets)
	t.Log(pretty.Sprint(markets))
	for _, market := range markets {
		require.Greater(t, market.MakerFee, 0.0)
		require.Greater(t, market.MinimumQuantity, 0.0)
		require.Greater(t, market.MinimumTickSize, 0.0)
		require.Greater(t, market.MinimumValue, 0.0)
	}
}
