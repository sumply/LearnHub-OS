package logger

import (
	"context"
	"encoding/json"
	"maps"
	"os"
	"reflect"
	"strings"
	"time"
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

type Map map[string]any

func (m *Map) Join(src ...Map) {
	for _, joinable := range src {
		maps.Copy((*m), joinable)
	}
}

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

type Stub struct{}

func (m *Stub) Debug(string) {
}
func (m *Stub) Info(string) {
}
func (m *Stub) Warn(string) {
}
func (m *Stub) Error(string) {
}

func (m *Stub) With(...TraceField) Logger {
	return &Stub{}
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
	f["timestamp"] = time.Now().UTC()
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
	maps.Copy(new, to)
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

func Mask(value string) (mask string) {
	return strings.Repeat("*", len(value))
}

func extractValue(v reflect.Value, t reflect.Type) any {
	if v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
		t = v.Type()
	}
	if t.Kind() != reflect.Struct || t == reflect.TypeOf(time.Time{}) {
		return v.Interface()
	}
	m := make(map[string]any)
	for i := range v.NumField() {
		fVal := v.Field(i)
		fType := t.Field(i)
		if !fType.IsExported() {
			continue
		}
		var printable any
		switch fType.Tag.Get("log") {
		case "hide":
			printable = "[HIDED]"
		case "mask":
			if fVal.Kind() != reflect.String {
				printable = "[MASKED]"
			} else {
				printable = strings.Repeat("*", fVal.Len())
			}
		default:
			printable = extractValue(fVal, fType.Type)
		}
		m[fType.Name] = printable
	}
	return m
}

func TraceFieldFromAny(a any) TraceField {
	m := MapFromAny(a)
	var new TraceField
	for key, val := range m {
		new.Key = key
		new.Value = val
	}
	return new
}

func MapFromAny(a any) Map {
	m := make(Map)
	t := reflect.TypeOf(a)
	v := reflect.ValueOf(a)
	if t.Kind() == reflect.Pointer {
		if v.IsNil() {
			m[t.Name()] = v.Interface()
			return m
		}
		v = v.Elem()
		t = v.Type()
	}
	if t.Kind() != reflect.Struct {
		m[t.Name()] = v.Interface()
		return m
	}

	m[t.Name()] = extractValue(v, t)
	return m
}
