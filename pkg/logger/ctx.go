package logger

import "context"

type key struct{}

func (l *Logger) SaveToContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, key{}, l)
}

func GetFromContext(ctx context.Context) *Logger {
	l, ok := ctx.Value(key{}).(*Logger)
	if !ok {
		return instance()
	}
	return l
}
