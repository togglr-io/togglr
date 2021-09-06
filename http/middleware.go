package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ctxKey struct {
	name string
}

func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, ctxKey{name: "requestID"}, requestID)
}

func GetRequestID(ctx context.Context) string {
	val, _ := ctx.Value(ctxKey{name: "requestID"}).(string)
	return val
}

type Middleware func(http.Handler) http.Handler

func Telemetry(logger *zap.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer logger.Sync()
			start := time.Now()
			requestID := uuid.New().String()
			newCtx := WithRequestID(r.Context(), requestID)
			r = r.WithContext(newCtx)
			log := logger.With(
				zap.String("path", r.URL.String()),
				zap.String("method", r.Method),
			)

			log.Info("request started",
				zap.String("requestID", requestID),
				zap.Time("time", start),
			)

			next.ServeHTTP(w, r)

			end := time.Now()
			log.Info("request finished",
				zap.String("requestID", requestID),
				zap.Time("time", end),
				zap.String("latency", fmt.Sprintf("%dms", end.Sub(start).Milliseconds())),
			)
		})
	}
}
