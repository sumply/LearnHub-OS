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
	log := fromCtx(ctx)
	log.DebugContext(ctx, msg, args...)
}

func Info(ctx context.Context, msg string, args ...any) {
	log := fromCtx(ctx)
	log.InfoContext(ctx, msg, args...)
}

func Warn(ctx context.Context, msg string, args ...any) {
	log := fromCtx(ctx)
	log.WarnContext(ctx, msg, args...)
}

func Error(ctx context.Context, msg string, args ...any) {
	log := fromCtx(ctx)
	log.ErrorContext(ctx, msg, args...)
}

func WithAttrs(ctx context.Context, args ...any) context.Context {
	log := fromCtx(ctx)
	log = log.With(args...)
	return withCtx(ctx, log)
}

func fromCtx(ctx context.Context) *slog.Logger {
	log, ok := ctx.Value(logCtx).(*slog.Logger)
	if !ok {
		log = slog.Default()
	}
	return log
}

func withCtx(ctx context.Context, log *slog.Logger) context.Context {
	return context.WithValue(ctx, logCtx, log)
}

func init() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))
	slog.SetDefault(log)
}
