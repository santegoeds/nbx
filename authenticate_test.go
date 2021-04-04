package nbx_test

import (
	"context"
	"github.com/santegoeds/nbx"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAuthenticate(t *testing.T) {
	accountID := os.Getenv("ACCOUNT_ID")
	key := os.Getenv("KEY")
	secret := os.Getenv("SECRET")
	passphrase := os.Getenv("PASSPHRASE")

	client := nbx.NewClient()
	err := client.Authenticate(context.TODO(), accountID, key, secret, passphrase, nbx.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, client.Token)
}
