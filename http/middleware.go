package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type RequestID string

const requestIDKey = "toggle/requestID"

func WithRequestID(ctx context.Context, requestID string) context.Context {
	key := RequestID(requestIDKey)
	return context.WithValue(ctx, key, requestID)
}

func GetRequestID(ctx context.Context) RequestID {
	val, _ := ctx.Value(requestIDKey).(RequestID)
	return val
}

type Middleware func(http.Handler) http.Handler

func Telemetry(logger *zap.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			requestID := uuid.New().String()
			newCtx := WithRequestID(r.Context(), requestID)
			r.WithContext(newCtx)
			logger = logger.With(
				zap.String("path", r.URL.String()),
				zap.String("method", r.Method),
			)

			logger.Info("request started",
				zap.String("requestID", requestID),
				zap.Time("time", start),
			)
			logger.Sync()

			next.ServeHTTP(w, r)

			end := time.Now()
			logger.Info("request finished",
				zap.String("requestID", requestID),
				zap.Time("time", end),
				zap.String("latency", fmt.Sprintf("%dms", end.Sub(start).Milliseconds())),
			)
			logger.Sync()
		})
	}
}
