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
	LevelDebug = slog.LevelDebug
	LevelInfo  = slog.LevelInfo
	LevelWarn  = slog.LevelWarn
	LevelError = slog.LevelError
)
const (
	LogEncodeText = iota
	LogEncodeJson
)

var (
	logger     *slog.Logger
	m          sync.Mutex
	logLevel   *slog.LevelVar
	callerSkip = 3 // 默认跳过层数：包装函数 -> logx.XXX -> slog.Handler
)

type Logger = slog.Logger
type Level = slog.Level

func init() {
	logLevel = new(slog.LevelVar)
	logLevel.Set(slog.LevelInfo)
}

func New(opts ...opt.Option[Options]) *Logger {
	o := &Options{logLevel: logLevel, source: true, callerSkip: callerSkip}
	opt.Each(o, opts...)
	var w io.Writer = os.Stderr
	if o.logOutType == LogOutByFile {
		w = o.logFileConfig
	}
	hOpts := &slog.HandlerOptions{AddSource: o.source, Level: o.logLevel}
	var h slog.Handler
	if o.logEncodeType == LogEncodeJson {
		h = slog.NewJSONHandler(w, hOpts)
	} else {
		h = slog.NewTextHandler(w, hOpts)
	}
	// 如果需要自定义 caller skip，使用 CallerHandler 包装
	if o.callerSkip != 0 && o.source {
		h = NewCallerHandler(h, o.callerSkip)
	}
	return slog.New(h)
}
func Build(opts ...opt.Option[Options]) {
	logger = New(opts...)
}

// SetCallerSkip 设置全局的 caller skip 层数
func SetCallerSkip(skip int) {
	m.Lock()
	defer m.Unlock()
	callerSkip = skip
}

// GetCallerSkip 获取当前的 caller skip 层数
func GetCallerSkip() int {
	return callerSkip
}

func SetLogger(l *slog.Logger) {
	logger = l
}
func SetLoggerLevel(level slog.Level) {
	logLevel.Set(level)
}
func Clone() *slog.Logger {
	l := defaultLogger()
	c := *l
	return &c
}

func GetLogger() *slog.Logger {
	return defaultLogger()
}

func defaultLogger() *slog.Logger {
	if logger == nil {
		m.Lock()
		if logger == nil {
			// 默认使用 Info 级别，输出到 stderr
			Build(WithLevel(LevelInfo), WithEncodeText())
		}
		m.Unlock()
	}
	return logger
}

func Debug(msg string, args ...any) {
	DebugContext(context.Background(), msg, args...)
}
func DebugContext(ctx context.Context, msg string, args ...any) {
	// 使用全局函数时，需要额外 +2 层 skip
	// 因为调用链：用户代码 -> Debug -> DebugContext -> logger.DebugContext
	if h, ok := defaultLogger().Handler().(*CallerHandler); ok {
		// 创建一个临时的 logger，额外 +2 skip
		tempHandler := NewCallerHandler(h.GetHandler(), h.GetSkip()+2)
		slog.New(tempHandler).DebugContext(ctx, msg, args...)
	} else {
		defaultLogger().DebugContext(ctx, msg, args...)
	}
}

func Info(msg string, args ...any) {
	InfoContext(context.Background(), msg, args...)
}
func InfoContext(ctx context.Context, msg string, args ...any) {
	if h, ok := defaultLogger().Handler().(*CallerHandler); ok {
		tempHandler := NewCallerHandler(h.GetHandler(), h.GetSkip()+2)
		slog.New(tempHandler).InfoContext(ctx, msg, args...)
	} else {
		defaultLogger().InfoContext(ctx, msg, args...)
	}
}

func Warn(msg string, args ...any) {
	WarnContext(context.Background(), msg, args...)
}
func WarnContext(ctx context.Context, msg string, args ...any) {
	if h, ok := defaultLogger().Handler().(*CallerHandler); ok {
		tempHandler := NewCallerHandler(h.GetHandler(), h.GetSkip()+2)
		slog.New(tempHandler).WarnContext(ctx, msg, args...)
	} else {
		defaultLogger().WarnContext(ctx, msg, args...)
	}
}
func Error(msg string, args ...any) {
	ErrorContext(context.Background(), msg, args...)
}
func ErrorContext(ctx context.Context, msg string, args ...any) {
	if h, ok := defaultLogger().Handler().(*CallerHandler); ok {
		tempHandler := NewCallerHandler(h.GetHandler(), h.GetSkip()+2)
		slog.New(tempHandler).ErrorContext(ctx, msg, args...)
	} else {
		defaultLogger().ErrorContext(ctx, msg, args...)
	}
}
