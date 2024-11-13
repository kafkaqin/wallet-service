package logger

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc/metadata"
	"runtime"
	"strings"
	"time"
)

var level = zap.NewAtomicLevelAt(zapcore.DebugLevel)

type Logger struct {
	z      *zap.Logger
	fields []Field
}

type Field struct {
	Key   string
	Value interface{}
}

var globalZapLogger = mustBuildZapLogger()

func NewLogger() *Logger {
	var l = &Logger{}
	l.z = globalZapLogger
	setLogLevelFromEnviron()
	return l
}

func mustBuildZapLogger() *zap.Logger {
	c := zap.NewProductionConfig()

	//设置日志输出路径
	p, err := getOutputPaths()
	if err != nil {
		//nothing
	} else {
		c.OutputPaths = []string{"stdout", p}
		c.ErrorOutputPaths = []string{"stderr"}
	}

	//设置日志指针
	c.Level = level

	//设置日志格式
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "time"
	customTimeEncoder := func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		l := Get()
		enc.AppendString(t.In(l).Format("2006-01-02 15:04:05.000"))
	}
	encoderConfig.EncodeTime = customTimeEncoder
	encoderConfig.EncodeCaller = zapcore.FullCallerEncoder

	c.EncoderConfig = encoderConfig

	z, err := c.Build(zap.AddCallerSkip(2))
	if err != nil {
		panic(err)
	}
	return z
}

func (l *Logger) GetZapLogger() *zap.Logger {
	return l.z
}

// ZapInstance 返回zap实例
func ZapInstance() *zap.Logger {
	return instance().GetZapLogger()
}

func getStack() string {
	var s string
	pc := make([]uintptr, 32)
	n := runtime.Callers(3, pc)
	for i := 0; i < n; i++ {
		f := runtime.FuncForPC(pc[i])
		file, line := f.FileLine(pc[i])
		var funcName string
		arr := strings.Split(f.Name(), "/")
		if len(arr) > 0 {
			funcName = arr[len(arr)-1]
		}
		if i != n-1 {
			s += fmt.Sprintf("%s:%d:%s---", file, line, funcName)
		} else {
			s += fmt.Sprintf("%s:%d:%s", file, line, funcName)
		}
	}
	return s
}

func (l *Logger) WithField(key string, value interface{}) *Logger {
	return l.WithFields([]Field{{key, value}}...)
}

func (l *Logger) WithFields(kvs ...Field) *Logger {
	logger := NewLogger()
	//加入原来的fields
	for _, v := range l.fields {
		logger.fields = append(logger.fields, v)
	}
	//增加新的fields
	for _, v := range kvs {
		logger.fields = append(logger.fields, v)
	}
	return logger
}

func (l *Logger) Debug(ctx context.Context, format string, args ...interface{}) {
	args, fs := fromArgs(args)
	output(ctx, l, zap.DebugLevel, fmt.Sprintf(format, args...), fs...)
}

func (l *Logger) Info(ctx context.Context, format string, args ...interface{}) {
	args, fs := fromArgs(args)
	output(ctx, l, zap.InfoLevel, fmt.Sprintf(format, args...), fs...)
}

func (l *Logger) Warn(ctx context.Context, format string, args ...interface{}) {
	args, fs := fromArgs(args)
	output(ctx, l, zap.WarnLevel, fmt.Sprintf(format, args...), fs...)
}

func (l *Logger) Error(ctx context.Context, format string, args ...interface{}) {
	args, fs := fromArgs(args)
	fs = append(fs, zap.String("stack", getStack()))
	output(ctx, l, zap.ErrorLevel, fmt.Sprintf(format, args...), fs...)
}

func (l *Logger) Panic(ctx context.Context, format string, args ...interface{}) {
	args, _ = fromArgs(args)
	output(ctx, l, zap.PanicLevel, fmt.Sprintf(format, args...))
}

type Level zapcore.Level

var (
	DebugLevel = Level(zapcore.DebugLevel)
	InfoLevel  = Level(zapcore.InfoLevel)
	WarnLevel  = Level(zapcore.WarnLevel)
	ErrorLevel = Level(zapcore.ErrorLevel)
	PanicLevel = Level(zapcore.PanicLevel)
)

func ParseAndSetLogLevel(text string) error {
	lvl, err := zapcore.ParseLevel(text)
	if err != nil {
		return err
	}
	level.SetLevel(lvl)
	return nil
}

func SetLogLevel(lvl Level) {
	level.SetLevel(zapcore.Level(lvl))
}

func Debug(ctx context.Context, format string, args ...interface{}) {
	l := GetFromContext(ctx)
	args, fs := fromArgs(args)
	output(ctx, l, zap.DebugLevel, fmt.Sprintf(format, args...), fs...)
}

func Info(ctx context.Context, format string, args ...interface{}) {
	l := GetFromContext(ctx)
	args, fs := fromArgs(args)
	output(ctx, l, zap.InfoLevel, fmt.Sprintf(format, args...), fs...)
}

func Warn(ctx context.Context, format string, args ...interface{}) {
	l := GetFromContext(ctx)
	args, fs := fromArgs(args)
	output(ctx, l, zap.WarnLevel, fmt.Sprintf(format, args...), fs...)
}

func Error(ctx context.Context, format string, args ...interface{}) {
	l := GetFromContext(ctx)
	args, fs := fromArgs(args)
	fs = append(fs, zap.String("stack", getStack()))
	output(ctx, l, zap.ErrorLevel, fmt.Sprintf(format, args...), fs...)
}

// ErrorAlert 记录日志并推送日志到告警群，一些关键的错误可以用在这里
func ErrorAlert(ctx context.Context, format string, args ...interface{}) {
	l := GetFromContext(ctx)
	args, fs := fromArgs(args)
	fs = append(fs, zap.String("stack", getStack()))
	output(ctx, l, zap.ErrorLevel, fmt.Sprintf(format, args...), fs...)
}

func Panic(ctx context.Context, format string, args ...interface{}) {
	l := GetFromContext(ctx)
	args, fs := fromArgs(args)
	output(ctx, l, zap.PanicLevel, fmt.Sprintf(format, args...), fs...)
}

func fromArgs(args []interface{}) ([]interface{}, []zap.Field) {
	var realArgs = make([]interface{}, 0, 8)
	var fs = make([]zap.Field, 0, 8)
	for _, arg := range args {
		f, ok := arg.(zap.Field)
		if !ok {
			realArgs = append(realArgs, arg)
		} else {
			fs = append(fs, f)
		}
	}
	return realArgs, fs
}

func output(ctx context.Context, l *Logger, level zapcore.Level, msg string, fs ...zap.Field) string {
	var fields = make([]zap.Field, 0)
	for _, v := range incomingKeys {
		f, ok := createZapFieldFromIncomingCtx(ctx, v)
		if !ok {
			continue
		}
		fields = append(fields, f)
	}
	for _, v := range outgoingKeys {
		f, ok := createZapFieldFromOutgoingCtx(ctx, v)
		if !ok {
			continue
		}
		fields = append(fields, f)
	}
	for _, v := range l.fields {
		fields = append(fields, zap.Any(v.Key, v.Value))
	}
	for _, f := range fs {
		fields = append(fields, f)
	}
	l.z.Log(level, msg, fields...)

	return msg + fieldsToString(fields)
}

func createZapFieldFromIncomingCtx(ctx context.Context, key string) (zap.Field, bool) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return zap.Field{}, false
	}
	var v = md.Get(key)
	if len(v) == 0 {
		return zap.Field{}, false
	}
	value := strings.Join(v, ",")
	field := zap.Any(key, value)
	return field, true
}

func createZapFieldFromOutgoingCtx(ctx context.Context, key string) (zap.Field, bool) {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		return zap.Field{}, false
	}
	var v = md.Get(key)
	if len(v) == 0 {
		return zap.Field{}, false
	}
	value := strings.Join(v, ",")
	field := zap.Any(key, value)
	return field, true
}

func fieldsToString(fields []zap.Field) string {
	return fmt.Sprint(fields)
}
