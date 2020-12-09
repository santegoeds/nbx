package api_v1_test

import (
	"context"
	"os"
	"testing"

	"github.com/santegoeds/nbx/api_v1"

	"github.com/stretchr/testify/require"
)

func TestAuthenticate(t *testing.T) {
	accountID := os.Getenv("ACCOUNT_ID")
	key := os.Getenv("KEY")
	secret := os.Getenv("SECRET")
	passphrase := os.Getenv("PASSPHRASE")

	client := api_v1.NewClient()
	err := client.Authenticate(context.TODO(), accountID, key, secret, passphrase, api_v1.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, client.Token)
}
