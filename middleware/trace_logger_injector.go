package middleware

import (
	"fmt"
	"net/http"

	"slog-example/logger"

	"golang.org/x/exp/slog"
)

// logger にリクエストの情報を加える
func LoggerInjector(l *logger.Logger, projectID string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			raw := []byte(r.Header.Get("X-Cloud-Trace-Context"))
			traceID, _ := ExtractTraceID(raw)

			newLogger := l.With(
				slog.String("logging.googleapis.com/trace", fmt.Sprintf("projects/%s/traces/%s", projectID, string(traceID))),
				slog.String("path", r.URL.Path),
			)

			newLogger.SetLoggerContexts(
				logger.LoggerContext{Key: "traceID", Value: traceID},
			)

			r = r.WithContext(logger.TraceLoggerWith(r.Context(), newLogger))

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}
