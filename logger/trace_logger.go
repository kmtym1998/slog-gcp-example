package logger

import (
	"context"
)

type traceLoggerCtxKey struct{}

// context に logger を詰める
func TraceLoggerWith(ctx context.Context, logger *Logger) context.Context {
	return context.WithValue(ctx, traceLoggerCtxKey{}, logger)
}

// context から Logger を抜き出す
func TraceLoggerFrom(ctx context.Context) (*Logger, bool) {
	traceLogger, ok := ctx.Value(traceLoggerCtxKey{}).(*Logger)

	return traceLogger, ok
}
