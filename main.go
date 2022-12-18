package main

import (
	"errors"
	"io"
	"log"
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

		OnError: func(l *logger.Logger, msg string, err error, arg ...any) {
			traceIDContext, ok := l.LoggerContext("traceID")
			if !ok {
				log.Println(msg)
				return
			}

			log.Printf("%s のリクエストでエラーが起きたよ\n", traceIDContext.Value)
		},
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}

	projectID := os.Getenv("PROJECT_ID")

	l.Debug("serving...")
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
	router.Use(middleware.LoggerInjector(l, projectID))

	// ルーティング
	router.Get("/", GetHealthCheckHandler())
	router.Get("/healthcheck", GetHealthCheckHandler())

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

func GetHealthCheckHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l, ok := logger.TraceLoggerFrom(r.Context())
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, "logger not found")
			return
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
