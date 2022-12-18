package middleware

import (
	"fmt"
	"net/http"
	"regexp"

	"slog-example/logger"

	"golang.org/x/exp/slog"
)

// logger にリクエストの情報を加える
func LoggerInjector(l *logger.Logger, projectID string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			raw := []byte(r.Header.Get("X-Cloud-Trace-Context"))
			traceID, _ := extractTraceID(raw)

			newLogger := l.With(
				slog.String("logging.googleapis.com/trace", fmt.Sprintf("projects/%s/traces/%s", projectID, traceID)),
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

func extractTraceID(raw []byte) (string, bool) {
	if (len(raw)) == 0 {
		return "", false
	}

	// NOTE: https://cloud.google.com/trace/docs/setup?hl=ja#force-trace
	matches := regexp.MustCompile(`([a-f\d]+)/([a-f\d]+)`).FindAllSubmatch(raw, -1)
	if len(matches) != 1 {
		return "", false
	}

	sub := matches[0]
	if len(sub) != 3 {
		return "", false
	}

	return string(sub[1]), true
}
