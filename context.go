package otplugins

import (
	"github.com/opentracing/opentracing-go"
	"golang.org/x/net/context"
)

func PropagateContextOnlyValue(old context.Context, keys ...interface{}) context.Context {
	ctx := context.Background()
	if len(keys) > 0 {
		for _, key := range keys {
			ctx = context.WithValue(ctx, key, old.Value(key))
		}
	}

	span := opentracing.SpanFromContext(old)
	ctx = opentracing.ContextWithSpan(ctx, span)
	return ctx
}
