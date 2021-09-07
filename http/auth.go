package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	"github.com/togglr-io/togglr"
	"github.com/togglr-io/togglr/uid"
	"go.uber.org/zap"
)

// the number of chars to skip when reading bearer tokens
const bearerWidth = 7

// The Claims that should be embedded in auth tokens.
type Claims struct {
	jwt.StandardClaims
	AccountID uid.UID `json:"accountId,omitempty"`
}

// Valid implements the jwt.Claims interface and determines whether or not the given claims are valid.
func (c Claims) Valid() error {
	if _, err := uid.FromString(c.Subject); err != nil {
		return fmt.Errorf("invalid subject: %w", err)
	}

	return c.StandardClaims.Valid()
}

// WithClaims takes a context and returns a new context containing the given Claims.
func WithClaims(ctx context.Context, claims Claims) context.Context {
	return context.WithValue(ctx, ctxKey{name: "claims"}, claims)
}

// GetClaims retrieves an embedded Claims struct from the given context.
func GetClaims(ctx context.Context) Claims {
	val, _ := ctx.Value(ctxKey{name: "claims"}).(Claims)
	return val
}

// A TokenFn generates a token containing claims based on the given User.
type TokenFn func(user togglr.User) (string, error)

// makeTokenFn creates a TokenFn that uses HS512 to sign with the given key.
func makeTokenFn(secret []byte) TokenFn {
	return func(user togglr.User) (string, error) {
		claims := Claims{
			StandardClaims: jwt.StandardClaims{
				Subject: user.ID.String(),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
		tokenString, err := token.SignedString(secret)
		if err != nil {
			return "", err
		}

		return tokenString, nil
	}
}

func parseToken(secret []byte, tokenString string) (Claims, error) {
	var claims Claims
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(tok *jwt.Token) (interface{}, error) {
		if _, ok := tok.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", tok.Header["alg"])
		}

		return secret, nil
	})

	if err != nil {
		return Claims{}, fmt.Errorf("failed to parse token: %w", err)
	}

	if token.Valid {
		return claims, nil
	}

	return Claims{}, errors.New("invalid token")
}

func requireAuth(secret []byte, log *zap.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// check for cookie first
			var token string
			for _, cookie := range r.Cookies() {
				if cookie.Name == "authToken" {
					token = cookie.Value
				}
			}

			// if token is still empty, check the auth header
			if token == "" {
				auth := r.Header.Get("Authorization")
				if len(auth) > bearerWidth {
					token = auth[bearerWidth:]
				}
			}

			// if token is _still_ empty, auth check is failed
			if token == "" {
				log.Error("could not find token")
				badRequest(w, "authorization missing")
				return
			}

			claims, err := parseToken(secret, token)
			if err != nil {
				log.Error("failed to parse token", zap.Error(err))
				unauthorized(w, "authentication failure")
				return
			}

			ctx := WithClaims(r.Context(), claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
