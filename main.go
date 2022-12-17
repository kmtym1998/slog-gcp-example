package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"slog-example/logger"
	"slog-example/middleware"
	"time"

	"github.com/go-chi/chi"
	"golang.org/x/exp/slog"
)

func main() {
	l := logger.New(logger.Opts{
		Level: slog.LevelDebug,

		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// NOTE: GCP の severity: https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry#logseverity
			switch {
			case a.Key == "level" && a.Value.String() == "WARN":
				return slog.String("severity", "WARNING")
			case a.Key == "level":
				return slog.String("severity", a.Value.String())
			case a.Key == "err":
				return slog.String("errorMessage", a.Value.String())
			}

			return a
		},

		OnError: func(l *logger.Logger, msg string, err error, arg ...any) {
			fmt.Println("Error が起きたよ")
		},
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}

	projectID := os.Getenv("PROJECT_ID")

	fmt.Println("serving...")
	if err := ListenAndServe(port, projectID, l); err != nil {
		panic(err)
	}
}

func ListenAndServe(
	port string,
	projectID string,
	l *logger.Logger,
) error {
	router := chi.NewRouter()

	// ミドルウェア
	router.Use(middleware.TraceLoggerInjector(l, projectID))

	// ルーティング
	router.Get("/healthcheck", GetHealthCheck(l))

	server := &http.Server{
		Addr:              ":" + port,
		ReadHeaderTimeout: 1 * time.Second,
		ReadTimeout:       3 * time.Second,
		WriteTimeout:      3 * time.Second,
		IdleTimeout:       1 * time.Second,
		Handler:           router,
	}

	return server.ListenAndServe()
}

func GetHealthCheck(l *logger.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l.Info(time.Now().Format(time.RFC3339))

		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "OK")
	}
}
