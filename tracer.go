package otplugins

import "github.com/opentracing/opentracing-go"

func NonNilTracer(t opentracing.Tracer) opentracing.Tracer {
	if t == nil {
		return opentracing.GlobalTracer()
	}
	return t
}
