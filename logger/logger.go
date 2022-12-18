package logger

import (
	"context"
	"fmt"
	"os"

	"github.com/pkg/errors"
	"golang.org/x/exp/slog"
)

type Logger struct {
	logger         *slog.Logger
	onError        func(l *Logger, msg string, err error, arg ...any)
	loggerContexts []LoggerContext
}

type Opts struct {
	Level   slog.Level
	OnError func(l *Logger, msg string, err error, arg ...any)
}

type LoggerContext struct {
	Key   string
	Value string
}

func New(opts Opts) *Logger {
	return &Logger{
		logger: slog.New(
			slog.HandlerOptions{
				Level: opts.Level,

				ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
					// NOTE: GCP の severity: https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry#logseverity
					switch {
					case a.Key == slog.MessageKey:
						return slog.String("message", a.Value.String())
					case a.Key == slog.LevelKey && a.Value.String() == slog.LevelWarn.String():
						return slog.String("severity", "WARNING")
					case a.Key == slog.LevelKey:
						return slog.String("severity", a.Value.String())
					case a.Key == slog.ErrorKey:
						return slog.String("errorMessage", a.Value.String())
					}

					return a
				},
			}.NewJSONHandler(os.Stdout),
		),
		onError: opts.OnError,
	}
}

func (l *Logger) With(args ...any) *Logger {
	return &Logger{
		logger:  l.logger.With(args...),
		onError: l.onError,
	}
}

func (l *Logger) LoggerContext(k string) (*LoggerContext, bool) {
	for _, v := range l.loggerContexts {
		if v.Key == k {
			return &v, true
		}
	}

	return nil, false
}

func (l *Logger) SetLoggerContexts(args ...LoggerContext) {
	l.loggerContexts = append(l.loggerContexts, args...)
}

func (l *Logger) Debug(msg string, arg ...any) {
	l.logger.Debug(msg, arg...)
}

func (l *Logger) Info(msg string, arg ...any) {
	l.logger.Info(msg, arg...)
}

func (l *Logger) Warning(msg string, arg ...any) {
	l.logger.Warn(msg, arg...)
}

func (l *Logger) Error(msg string, err error, arg ...any) {
	// stacktrace を表示
	arg = append(arg, slog.String("stack", fmt.Sprintf("%+v", errors.WithStack(err))))
	l.logger.Error(msg, err, arg...)

	go func() {
		// エラーログ出力後なにかやりたい時 (Sentry に送るとか) は OnError() を呼び元から渡す
		l.onError(l, msg, err, arg...)
	}()
}

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
