package logger

import (
	"context"
	"log/slog"
	"os"
)

type ctxKey string

const (
	logCtx ctxKey = "logger"
)

func Debug(ctx context.Context, msg string, args ...any) {
	log := FromCtx(ctx)
	log.DebugContext(ctx, msg, args...)
}

func Info(ctx context.Context, msg string, args ...any) {
	log := FromCtx(ctx)
	log.InfoContext(ctx, msg, args...)
}

func Warn(ctx context.Context, msg string, args ...any) {
	log := FromCtx(ctx)
	log.WarnContext(ctx, msg, args...)
}

func Error(ctx context.Context, msg string, args ...any) {
	log := FromCtx(ctx)
	log.ErrorContext(ctx, msg, args...)
}

func WithAttrs(ctx context.Context, args ...any) context.Context {
	log := FromCtx(ctx)
	log = log.With(args...)
	return WithCtx(ctx, log)
}

func FromCtx(ctx context.Context) *slog.Logger {
	log, ok := ctx.Value(logCtx).(*slog.Logger)
	if !ok {
		log = slog.Default()
	}
	return log
}

func WithCtx(ctx context.Context, log *slog.Logger) context.Context {
	return context.WithValue(ctx, logCtx, log)
}

func init() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))
	slog.SetDefault(log)
}
