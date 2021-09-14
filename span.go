package otplugins

import (
	"context"
	"github.com/opentracing/opentracing-go"
)

func NoneNilSpanFromContext(ctx context.Context) opentracing.Span {
	span := opentracing.SpanFromContext(ctx)
	if span == nil {
		span = NoopSpan{}
	}
	return span
}

func NoneNilChildSpanFromContext(spanName string, ctx context.Context) opentracing.Span {
	span := opentracing.SpanFromContext(ctx)
	if span != nil {
		span = span.Tracer().StartSpan(spanName, opentracing.ChildOf(span.Context()))
	} else {
		span = NoopSpan{}
	}
	return span
}
