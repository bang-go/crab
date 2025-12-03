package logx

import (
	"context"
	"io"
	"log/slog"
	"os"
	"sync"

	"github.com/bang-go/opt"
)

const (
	LogOutByStdout = iota
	LogOutByFile
)
const (
	LogEncodeText = iota
	LogEncodeJson
)
const (
	LevelDebug = slog.LevelDebug
	LevelInfo  = slog.LevelInfo
	LevelWarn  = slog.LevelWarn
	LevelError = slog.LevelError
)

var (
	logger   *slog.Logger
	m        sync.Mutex
	logLevel *slog.LevelVar
)

type Logger = slog.Logger
type Level = slog.Level

func init() {
	logLevel = new(slog.LevelVar)
	logLevel.Set(slog.LevelInfo)
}

func New(opts ...opt.Option[Options]) *Logger {
	o := &Options{logLevel: logLevel, source: true}
	opt.Each(o, opts...)
	var w io.Writer
	switch o.logOutType {
	case LogOutByFile:
		w = o.logFileConfig
	default:
		w = os.Stderr
	}
	var h slog.Handler
	hOpts := &slog.HandlerOptions{AddSource: o.source, Level: o.logLevel}
	switch o.logEncodeType {
	case LogEncodeJson:
		h = slog.NewJSONHandler(w, hOpts)
	default:
		h = slog.NewTextHandler(w, hOpts)
	}
	return slog.New(h)
}
func Build(opts ...opt.Option[Options]) {
	logger = New(opts...)
	return
}
func SetLogger(l *slog.Logger) {
	logger = l
}
func SetLoggerLevel(level slog.Level) {
	logLevel.Set(level)
}
func Clone() *slog.Logger {
	c := *logger
	return &c
}
func GetLogger() *slog.Logger {
	return defaultLogger()
}

func defaultLogger() *slog.Logger {
	if logger == nil {
		m.Lock()
		Build()
		m.Unlock()
	}
	return logger
}

func Debug(msg string, args ...any) {
	DebugContext(context.Background(), msg, args...)
}
func DebugContext(ctx context.Context, msg string, args ...any) {
	defaultLogger().DebugContext(ctx, msg, args...)
}

func Info(msg string, args ...any) {
	InfoContext(context.Background(), msg, args...)
}
func InfoContext(ctx context.Context, msg string, args ...any) {
	defaultLogger().InfoContext(ctx, msg, args...)
}

func Warn(msg string, args ...any) {
	WarnContext(context.Background(), msg, args...)
}
func WarnContext(ctx context.Context, msg string, args ...any) {
	defaultLogger().WarnContext(ctx, msg, args...)
}
func Error(msg string, args ...any) {
	ErrorContext(context.Background(), msg, args...)
}
func ErrorContext(ctx context.Context, msg string, args ...any) {
	defaultLogger().ErrorContext(ctx, msg, args...)
}
