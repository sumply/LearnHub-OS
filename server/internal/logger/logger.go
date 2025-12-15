package logger

import (
	"context"
	"fmt"
)

type ctxKey string

const (
	tracedKey ctxKey = "traced fields"
	loggerKey ctxKey = "logger"
)

type TraceField struct {
	Key   string
	Value any
}

type Logger interface {
	Debug(string)
	Info(string)
	Warn(string)
	Error(string)

	With(...TraceField) Logger
}

type Mock struct {
	ctx context.Context
}

func NewMock() *Mock {
	return &Mock{
		ctx: context.Background(),
	}
}

func (m *Mock) Debug(msg string) {
	m.print("DEBUG", msg)
}
func (m *Mock) Info(msg string) {
	m.print("INFO", msg)
}
func (m *Mock) Warn(msg string) {
	m.print("WARN", msg)
}
func (m *Mock) Error(msg string) {
	m.print("ERROR", msg)
}

func (m *Mock) With(fields ...TraceField) Logger {
	f, ok := getTracedFields(m.ctx)
	if !ok {
		f = make(map[string]any)
	}
	new := copyTracedFields(f, fields...)
	return &Mock{
		ctx: withTracedFields(m.ctx, new),
	}
}

func (m *Mock) print(layer string, msg string) {
	f, ok := getTracedFields(m.ctx)
	if !ok {
		f = make(map[string]any)
	}
	fmt.Printf("%s: %v; message: %s\n", layer, f, msg)
}

func getTracedFields(ctx context.Context) (map[string]any, bool) {
	f, ok := ctx.Value(tracedKey).(map[string]any)
	return f, ok
}

func withTracedFields(ctx context.Context, fields map[string]any) context.Context {
	return context.WithValue(ctx, tracedKey, fields)
}

func copyTracedFields(to map[string]any, from ...TraceField) map[string]any {
	new := make(map[string]any)
	for key, val := range to {
		new[key] = val
	}
	for _, f := range from {
		new[f.Key] = f.Value
	}
	return new
}

var newFunc func() Logger

func SetNewFunc(new func() Logger) {
	newFunc = new
}

func FromCtx(ctx context.Context) Logger {
	f, ok := ctx.Value(loggerKey).(Logger)
	if !ok {
		return newFunc()
	}
	return f
}

func WithLoggerCtx(ctx context.Context, l Logger) context.Context {
	return context.WithValue(ctx, loggerKey, l)
}
