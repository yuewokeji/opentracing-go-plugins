package othttp

import (
	"github.com/opentracing/opentracing-go"
	"net/http"
)

type ClientIPFunc func(r *http.Request) string

func RequestWithContext(r *http.Request, s opentracing.Span) *http.Request {
	return r.WithContext(opentracing.ContextWithSpan(r.Context(), s))
}

func SpanContextFormRequest(r *http.Request, tracer opentracing.Tracer) (opentracing.SpanContext, error) {
	return tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
}

func ClientSpanName(req string) string {
	const prefix = "http.client:"
	return prefix + req
}

func ServerSpanName(serverName, req string) string {
	const prefix = "http.server:"
	if serverName != "" {
		return prefix + serverName + req
	}
	return prefix + req
}

func IsStatusError(code int) bool {
	if code >= 500 && code < 599 {
		return true
	}
	return false
}

type RequestIgnoreFunc func(r *http.Request) bool

func NoopRequestIgnore(r *http.Request) bool {
	return false
}

type RequestNameFunc func(r *http.Request) string

func URLPathRequestNameFunc(r *http.Request) string {
	if r.URL.Path == "" {
		return "/"
	}
	return r.URL.Path
}
