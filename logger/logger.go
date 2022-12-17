package logger

import (
	"os"

	"golang.org/x/exp/slog"
)

type Logger struct {
	logger  *slog.Logger
	onError func(l *Logger, msg string, err error, arg ...any)
}

type Opts struct {
	Level       slog.Level
	ReplaceAttr func(groups []string, a slog.Attr) slog.Attr
	OnError     func(l *Logger, msg string, err error, arg ...any)
}

func New(opts Opts) *Logger {
	return &Logger{
		logger: slog.New(
			slog.HandlerOptions{
				Level:       opts.Level,
				ReplaceAttr: opts.ReplaceAttr,
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
	l.logger.Error(msg, err, arg...)

	go func() {
		// エラーログ出力後なにかやりたい時 (sentry に送るとか) は OnError() を呼び元から渡す
		l.onError(l, msg, err, arg...)
	}()
}
