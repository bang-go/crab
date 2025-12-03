package logz

import (
	"errors"
	"os"
	"sync"
	"syscall"

	"github.com/bang-go/opt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	logger              *zap.Logger
	m                   sync.Mutex
	defaultCallSkip     = 1
	defaultLevelEnabler = zap.DebugLevel
)

const (
	DefaultConfigKindDev = iota
	DefaultConfigKindProd
	LogOutByStdout = iota
	LogOutByFile
	DebugLevel    = zap.DebugLevel
	InfoLevel     = zap.InfoLevel
	WarnLevel     = zap.WarnLevel
	ErrorLevel    = zap.ErrorLevel
	DPanicLevel   = zap.DPanicLevel
	PanicLevel    = zap.PanicLevel
	FatalLevel    = zap.FatalLevel
	EncodeJson    = "json"
	EncodeConsole = "console"
)

type Level = zapcore.Level
type Logger = zap.Logger
type FileConfig = lumberjack.Logger

var (
	Skip       = zap.Skip
	Binary     = zap.Binary
	Bool       = zap.Bool
	Boolp      = zap.Boolp
	ByteString = zap.ByteString
	Float64    = zap.Float64
	Float64p   = zap.Float64p
	Float32    = zap.Float32
	Float32p   = zap.Float32p
	String     = zap.String
	Stringp    = zap.Stringp
	Uint       = zap.Uint
	Uintp      = zap.Uintp
	Uint8      = zap.Uint8
	Uint8p     = zap.Uint8p
	Uint32     = zap.Uint32
	Uint32p    = zap.Uint32p
	Uint64     = zap.Uint64
	Uint64p    = zap.Uint64p
	Int        = zap.Int
	Intp       = zap.Intp
	Int8       = zap.Int8
	Int8p      = zap.Int8p
	Int32      = zap.Int32
	Int32p     = zap.Int32p
	Int64      = zap.Int64
	Int64p     = zap.Int64p
	Duration   = zap.Duration
	Durationp  = zap.Durationp
	Any        = zap.Any
)

func New(opts ...opt.Option[Options]) *Logger {
	o := &Options{
		callerSkip:   defaultCallSkip,
		levelEnabler: defaultLevelEnabler,
	}
	opt.Each(o, opts...)
	if o.logOutType == 0 {
		o.logStdout = true
		o.logOutType |= 1
	}
	if o.zapEncoder == nil {
		o.zapEncoder = NewDefaultDevEncoder()
	}
	var writeSyncers []zapcore.WriteSyncer
	if o.logStdout {
		writeSyncers = append(writeSyncers, zapcore.AddSync(os.Stdout))
	}
	if o.logFileConfig != nil {
		writeSyncers = append(writeSyncers, zapcore.AddSync(o.logFileConfig))
	}
	zapOptions := append(o.zapOption, zap.AddCaller(), zap.AddCallerSkip(o.callerSkip))

	core := zapcore.NewCore(o.zapEncoder, zapcore.NewMultiWriteSyncer(writeSyncers...), o.levelEnabler)
	return zap.New(core, zapOptions...)
}

func Build(opts ...opt.Option[Options]) {
	logger = New(opts...)
}

func SetLogger(l *zap.Logger) {
	logger = l
}

func defaultLogger() *zap.Logger {
	if logger == nil {
		m.Lock()
		if logger == nil {
			Build()
		}
		m.Unlock()
	}
	return logger
}

func GetLogger() *zap.Logger {
	return defaultLogger()
}

func Clone() *zap.Logger {
	l := defaultLogger()
	c := *l
	return &c
}

func NewDefaultProdEncoder() zapcore.Encoder {
	return zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	})
}

func NewDefaultDevEncoder() zapcore.Encoder {
	return zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		TimeKey:        "T",
		LevelKey:       "L",
		NameKey:        "N",
		CallerKey:      "C",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "M",
		StacktraceKey:  "S",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	})
}

func Debug(msg string, fields ...zap.Field) {
	defaultLogger().Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	defaultLogger().Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	defaultLogger().Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	defaultLogger().Error(msg, fields...)
}

func DPanic(msg string, fields ...zap.Field) {
	defaultLogger().DPanic(msg, fields...)
}

func Panic(msg string, fields ...zap.Field) {
	defaultLogger().Panic(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	defaultLogger().Fatal(msg, fields...)
}

func Sync() error {
	err := defaultLogger().Sync()
	// NOTE: we use syscall.EBADF to check if the error is specifically related to a bad file descriptor,
	// which should be the case for if the stderr is a TTY.
	if err != nil && (!errors.Is(err, syscall.EBADF) && !errors.Is(err, syscall.ENOTTY)) {
		return err
	}
	return nil
}
