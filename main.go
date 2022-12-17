package main

import (
	"errors"
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
			case a.Key == "msg":
				return slog.String("message", a.Value.String())
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
	router.Get("/healthcheck", GetHealthCheck())

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

func GetHealthCheck() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l, ok := logger.TraceLoggerFrom(r.Context())
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, "logger not found")
		}

		err := errors.New("ErrorErrorErrorErrorError")
		l.Debug("DebugDebugDebugDebugDebug")
		l.Info("InfoInfoInfoInfoInfo")
		l.Warning("WarningWarningWarningWarningWarning")
		l.Error(err.Error(), err)

		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "OK")
	}
}
