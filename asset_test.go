package nbx_test

import (
	"context"
	"github.com/santegoeds/nbx"
	"os"
	"testing"

	"github.com/kr/pretty"
	"github.com/stretchr/testify/require"
)

func TestGetAsset(t *testing.T) {
	accountID := os.Getenv("ACCOUNT_ID")
	key := os.Getenv("KEY")
	secret := os.Getenv("SECRET")
	passphrase := os.Getenv("PASSPHRASE")

	client := nbx.NewClient()
	err := client.Authenticate(context.TODO(), accountID, key, secret, passphrase, nbx.Minute)
	require.NoError(t, err)

	asset, err := client.Asset(context.TODO(), "ETH")
	require.NoError(t, err)
	require.NotNil(t, asset)

	pretty.Println(asset)

	require.GreaterOrEqual(t, asset.Balance.Available, 0.0)
	require.GreaterOrEqual(t, asset.Balance.Total, asset.Balance.Available)
	for _, freeze := range asset.Freezes {
		require.NotEmpty(t, freeze.ID)
		require.Greater(t, freeze.Amount, 0)
	}
}
