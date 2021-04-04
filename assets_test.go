package nbx_test

import (
	"context"
	"github.com/santegoeds/nbx"
	"os"
	"testing"

	"github.com/kr/pretty"
	"github.com/stretchr/testify/require"
)

func TestGetAssets(t *testing.T) {
	accountID := os.Getenv("ACCOUNT_ID")
	key := os.Getenv("KEY")
	secret := os.Getenv("SECRET")
	passphrase := os.Getenv("PASSPHRASE")

	client := nbx.NewClient()
	err := client.Authenticate(context.TODO(), accountID, key, secret, passphrase, nbx.Minute)
	require.NoError(t, err)

	assets, err := client.Assets(context.TODO())
	require.NoError(t, err)
	require.NotEmpty(t, assets)

	pretty.Println(assets)

	for _, asset := range assets {
		require.GreaterOrEqual(t, asset.Balance.Available, 0.0)
		require.GreaterOrEqual(t, asset.Balance.Total, asset.Balance.Available)
		for _, freeze := range asset.Freezes {
			require.NotEmpty(t, freeze.ID)
			require.Greater(t, freeze.Amount, 0)
		}
	}
}
