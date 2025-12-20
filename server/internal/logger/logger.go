package logger

import (
	"context"
	"encoding/json"
	"os"
)

type Level uint8

var level Level = DEBUG

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
)

func SetLayer(l Level) {
	level = l
}

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

type Fake struct {
	ctx context.Context
}

func NewFake() *Fake {
	return &Fake{
		ctx: context.Background(),
	}
}

func (m *Fake) Debug(msg string) {
	if level > DEBUG {
		return
	}
	m.print("DEBUG", msg)
}
func (m *Fake) Info(msg string) {
	if level > INFO {
		return
	}
	m.print("INFO", msg)
}
func (m *Fake) Warn(msg string) {
	if level > WARN {
		return
	}
	m.print("WARN", msg)
}
func (m *Fake) Error(msg string) {
	if level > ERROR {
		return
	}
	m.print("ERROR", msg)
}

func (m *Fake) With(fields ...TraceField) Logger {
	f, ok := getTracedFields(m.ctx)
	if !ok {
		f = make(map[string]any)
	}
	new := copyTracedFields(f, fields...)
	return &Fake{
		ctx: withTracedFields(m.ctx, new),
	}
}

func (m *Fake) print(level string, msg string) {
	f, ok := getTracedFields(m.ctx)
	if !ok {
		f = make(map[string]any)
	}
	f["level"] = level
	f["message"] = msg
	e := json.NewEncoder(os.Stdout)
	e.SetIndent("=", "\t")
	e.Encode(f)
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

type NewFunc func() Logger

var newFunc NewFunc

func SetNewFunc(new NewFunc) {
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

func OnDebug(f func()) {
	if level <= DEBUG {
		f()
	}
}
