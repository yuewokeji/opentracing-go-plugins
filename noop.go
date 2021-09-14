package otplugins

import (
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type NoopSpan struct{}
type NoopSpanContext struct{}

var (
	defaultNoopSpanContext opentracing.SpanContext = NoopSpanContext{}
	defaultNoopSpan        opentracing.Span        = NoopSpan{}
	defaultNoopTracer      opentracing.Tracer      = opentracing.NoopTracer{}
)

const (
	emptyString = ""
)

// NoopSpanContext:
func (n NoopSpanContext) ForeachBaggageItem(handler func(k, v string) bool) {}

// NoopSpan:
func (n NoopSpan) Context() opentracing.SpanContext                       { return defaultNoopSpanContext }
func (n NoopSpan) SetBaggageItem(key, val string) opentracing.Span        { return n }
func (n NoopSpan) BaggageItem(key string) string                          { return emptyString }
func (n NoopSpan) SetTag(key string, value interface{}) opentracing.Span  { return n }
func (n NoopSpan) LogFields(fields ...log.Field)                          {}
func (n NoopSpan) LogKV(keyVals ...interface{})                           {}
func (n NoopSpan) Finish()                                                {}
func (n NoopSpan) FinishWithOptions(opts opentracing.FinishOptions)       {}
func (n NoopSpan) SetOperationName(operationName string) opentracing.Span { return n }
func (n NoopSpan) Tracer() opentracing.Tracer                             { return defaultNoopTracer }
func (n NoopSpan) LogEvent(event string)                                  {}
func (n NoopSpan) LogEventWithPayload(event string, payload interface{})  {}
func (n NoopSpan) Log(data opentracing.LogData)                           {}
