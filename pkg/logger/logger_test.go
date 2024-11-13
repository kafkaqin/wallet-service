package logger

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"testing"
	"time"
)

func TestNewLogger(t *testing.T) {
	var ctx = context.Background()
	var name = "virgil"
	Error(ctx, "my name is %v", name)
	Warn(ctx, "my name is %v", name)
	Info(ctx, "my name is %v", name)
	Debug(ctx, "my name is %v", name)
}

func TestWithField(t *testing.T) {
	var ctx = context.Background()
	var l = NewLogger()
	l = l.WithFields(Field{"name", "dante"}, Field{"age", 23})
	l.Error(ctx, "now is %v", time.Now())
	l.Warn(ctx, "now is %v", time.Now())
	l.Info(ctx, "now is %v", time.Now())
	l.Debug(ctx, "now is %v", time.Now())

	Error(ctx, "now is %v", time.Now())
	Warn(ctx, "now is %v", time.Now())
	Info(ctx, "now is %v", time.Now())
	Debug(ctx, "now is %v", time.Now())

	l = l.WithField("year", 2023)
	l.Error(ctx, "now is %v", time.Now())
	l.Warn(ctx, "now is %v", time.Now())
	l.Info(ctx, "now is %v", time.Now())
	l.Debug(ctx, "now is %v", time.Now())

	Error(ctx, "now is %v", time.Now())
	Warn(ctx, "now is %v", time.Now())
	Info(ctx, "now is %v", time.Now())
	Debug(ctx, "now is %v", time.Now())
}

func TestWithContext(t *testing.T) {
	var ctx = context.Background()

	var l = NewLogger()
	l = l.WithFields(Field{"name", "dante"}, Field{"age", 23})
	l.Error(ctx, "now is %v", time.Now())
	l.Warn(ctx, "now is %v", time.Now())
	l.Info(ctx, "now is %v", time.Now())
	l.Debug(ctx, "now is %v", time.Now())

	Error(ctx, "now is %v", time.Now())
	Warn(ctx, "now is %v", time.Now())
	Info(ctx, "now is %v", time.Now())
	Debug(ctx, "now is %v", time.Now())

	ctx = l.SaveToContext(ctx)
	Error(ctx, "now is %v", time.Now())
	Warn(ctx, "now is %v", time.Now())
	Info(ctx, "now is %v", time.Now())
	Debug(ctx, "now is %v", time.Now())
}

func TestDuplicateFieldName(t *testing.T) {
	var ctx = context.Background()

	var l = NewLogger()
	l = l.WithFields(Field{"name", "dante"}, Field{"age", 23})
	l.Error(ctx, "now is %v", time.Now())
	l.Warn(ctx, "now is %v", time.Now())
	l.Info(ctx, "now is %v", time.Now())
	l.Debug(ctx, "now is %v", time.Now())

	l = l.WithFields(Field{"name", "virgil"}, Field{"age", 25})
	l.Error(ctx, "now is %v", time.Now())
	l.Warn(ctx, "now is %v", time.Now())
	l.Info(ctx, "now is %v", time.Now())
	l.Debug(ctx, "now is %v", time.Now())
}

func TestGetFromContextAndWithField(t *testing.T) {
	var ctx = context.Background()

	l := GetFromContext(ctx)

	nl := l.WithField("name", "virgil")

	l.Info(ctx, "now is %v", time.Now())
	nl.Info(ctx, "now is %v", time.Now())
	Info(ctx, "now is %v", time.Now())
}

func TestLogLevel(t *testing.T) {
	var ctx = context.Background()

	Error(ctx, "now is %v", time.Now())
	Warn(ctx, "now is %v", time.Now())
	Info(ctx, "now is %v", time.Now())
	Debug(ctx, "now is %v", time.Now())

	SetLogLevel(DebugLevel)
	Error(ctx, "now is %v", time.Now())
	Warn(ctx, "now is %v", time.Now())
	Info(ctx, "now is %v", time.Now())
	Debug(ctx, "now is %v", time.Now())

	SetLogLevel(ErrorLevel)
	l := GetFromContext(ctx)
	l.Error(ctx, "now is %v", time.Now())
	l.Warn(ctx, "now is %v", time.Now())
	l.Info(ctx, "now is %v", time.Now())
	l.Debug(ctx, "now is %v", time.Now())

	SetLogLevel(WarnLevel)
	Error(ctx, "now is %v", time.Now())
	Warn(ctx, "now is %v", time.Now())
	Info(ctx, "now is %v", time.Now())
	Debug(ctx, "now is %v", time.Now())

	err := ParseAndSetLogLevel("debug")
	if err != nil {
		panic(err)
	}
	Error(ctx, "now is %v", time.Now())
	Warn(ctx, "now is %v", time.Now())
	Info(ctx, "now is %v", time.Now())
	Debug(ctx, "now is %v", time.Now())
}

func TestZapField(t *testing.T) {
	var ctx = context.Background()
	var err = errors.New("some error")
	Info(ctx, "hello world, now is %v", time.Now().Format(time.DateTime), zap.String("name", "virgil"), zap.Error(err))
	Error(ctx, "hello world, now is %v", time.Now().Format(time.DateTime), zap.String("name", "virgil"), zap.Error(err))
}

func TestCompareStruct(t *testing.T) {
	type s1 struct {
	}
	type s2 struct {
	}
	var a interface{} = s1{}
	var b interface{} = s2{}
	fmt.Println(a == b)
}

func TestLoggerContextWithFields(t *testing.T) {
	var ctx = context.Background()
	lg := GetFromContext(ctx)
	lg = lg.WithField("test", 1)
	lg.Info(ctx, "print", zap.Any("name", "virgil"))
}
