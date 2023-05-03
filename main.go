package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"slog-example/logger"
	"slog-example/middleware"
	"time"

	"github.com/go-chi/chi"
	"github.com/gofrs/uuid/v5"
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
	router.Post("/user", PostUserHandler())

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

type User struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Password Secret    `json:"password"`
}

type Secret string

func (s Secret) String() string {
	return "********"
}

func PostUserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l, ok := logger.TraceLoggerFrom(r.Context())
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, "logger not found")
			return
		}

		user := User{
			ID:       uuid.Must(uuid.NewV4()),
			Name:     "kmtym1998",
			Password: "password",
		}

		fmt.Println(user)
		fmt.Println(user.Password)
		l.Debug(fmt.Errorf("failed to create user: %+v", user).Error())
	}
}
