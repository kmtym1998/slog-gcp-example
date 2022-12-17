package middleware

import (
	"fmt"
	"net/http"

	"slog-example/logger"

	"golang.org/x/exp/slog"
)

// traceID 付きでロギングできるロガーを context に詰める
func TraceLoggerInjector(l *logger.Logger, projectID string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			header := []byte(r.Header.Get("X-Cloud-Trace-Context"))
			newLogger := injectTraceID(l, header, projectID)

			r = r.WithContext(logger.TraceLoggerWith(r.Context(), newLogger))

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}

// traceID 付きでロギングできるロガーを生成する
func injectTraceID(l *logger.Logger, raw []byte, projectID string) *logger.Logger {
	traceID, ok := ExtractTraceID(raw)
	if !ok {
		return l
	}

	return l.With(
		slog.String("logging.googleapis.com/trace", fmt.Sprintf("projects/%s/traces/%s", projectID, string(traceID))),
	)
}
