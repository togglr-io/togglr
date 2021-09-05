package hmac

import (
	"crypto/hmac"
	"crypto/sha512"
	"errors"

	"github.com/togglr-io/togglr/env"
)

// our hmac hashes are 512 bits in length, so we're grabbing the byte length
const hmacHashLength = 512 / 8

// A Signer implements the togglr.Signer interface using HMACSHA512 to generate signatures
type Signer struct {
	secret []byte
}

// NewSigner creates a new Signer with the given key
func NewSigner(secret string) Signer {
	return Signer{
		secret: []byte(secret),
	}
}

// FromEnv creates a new Signer using a secret found at some environment key
func FromEnv(key string) Signer {
	return NewSigner(env.GetString(key, "TOGGLR_SIGNING_KEY"))
}

func (s Signer) getSignature(data []byte) ([]byte, error) {
	mac := hmac.New(sha512.New, s.secret)
	if _, err := mac.Write(data); err != nil {
		return nil, err
	}

	return mac.Sum(nil), nil
}

// Sign some data using HMACSHA512 and the configured key
func (s Signer) Sign(data []byte) ([]byte, error) {
	sig, err := s.getSignature(data)
	if err != nil {
		return nil, err
	}

	return append(data, sig...), nil
}

// Validate some signed data, returning the original payload if validation succeeds
func (s Signer) Validate(signed []byte) ([]byte, error) {
	sigStart := len(signed) - hmacHashLength
	data := signed[:sigStart]
	signature := signed[sigStart:]

	// sign the source data and compare it to the signature given
	expected, err := s.getSignature(data)
	if err != nil {
		return nil, err
	}

	for i, b := range signature {
		// if the signature given doesn't match the new signature,
		// the payload is invalid
		if b != expected[i] {
			return nil, errors.New("invalid signature")
		}
	}

	return signed, nil
}
