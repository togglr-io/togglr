package token

import (
	"fmt"

	"github.com/golang-jwt/jwt/v4"
	"github.com/togglr-io/togglr/uid"
)

type Tokenizer struct {
	secret []byte
}

type Claims struct {
	jwt.StandardClaims
	AccountID uid.UID `json:"accountId"`
}

func (c Claims) Valid() error {
	if _, err := uid.FromString(c.Subject); err != nil {
		return fmt.Errorf("invalid subject: %w", err)
	}

	return c.StandardClaims.Valid()
}

func (t Tokenizer) Token(userID, accountID uid.UID) (string, error) {
	claims := Claims{
		StandardClaims: jwt.StandardClaims{
			Subject: userID.String(),
		},
		AccountID: accountID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	tokenString, err := token.SignedString(t.secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
