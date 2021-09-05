package hmac_test

import (
	"testing"

	"github.com/togglr-io/togglr/hmac"
)

func Test_Signer(t *testing.T) {
	secret := "testing"
	signer := hmac.NewSigner(secret)
	src := []byte("Hello, world!")

	signed, err := signer.Sign(src)
	if err != nil {
		t.Fatal("failed to sign data")
	}

	original, err := signer.Validate(signed)
	if err != nil {
		t.Fatalf("failed to validate signed data: %s", err)
	}

	for idx, b := range src {
		if b != original[idx] {
			t.Fatal("decoded data does not match src")
		}
	}
}
