package othttp

import (
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/yuewokeji/opentracing-go-plugins"
	"net/http"
)

func WrapClient(c *http.Client, options ...Option) *http.Client {
	if c == nil {
		c = http.DefaultClient
	}

	copied := *c
	copied.Transport = WrapRoundTripper(copied.Transport, options...)

	return &copied
}

func WrapRoundTripper(r http.RoundTripper, options ...Option) http.RoundTripper {
	if r == nil {
		r = http.DefaultTransport
	}

	r2 := &roundTripper{
		r:                 r,
		requestIgnoreFunc: NoopRequestIgnore,
		requestNameFunc: URLPathRequestNameFunc,
	}

	for _, option := range options {
		option(r2)
	}

	return r2
}

type roundTripper struct {
	r http.RoundTripper

	tracer opentracing.Tracer

	requestIgnoreFunc RequestIgnoreFunc
	requestNameFunc   RequestNameFunc
}

func (r *roundTripper) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	if r.requestIgnoreFunc(req) {
		return r.r.RoundTrip(req)
	}

	span := opentracing.SpanFromContext(req.Context())
	if span == nil {
		span = otplugins.NonNilTracer(r.tracer).StartSpan(ClientSpanName(r.requestNameFunc(req)))
	} else {
		span = otplugins.NonNilTracer(r.tracer).StartSpan(ClientSpanName(r.requestNameFunc(req)), opentracing.ChildOf(span.Context()))
	}
	defer span.Finish()

	otplugins.NonNilTracer(r.tracer).Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header))

	resp, err = r.r.RoundTrip(req)
	if err != nil {
		ext.LogError(span, err)
	} else {
		ext.HTTPStatusCode.Set(span, uint16(resp.StatusCode))
	}

	ext.SpanKindRPCClient.Set(span)
	ext.HTTPMethod.Set(span, req.Method)
	ext.HTTPUrl.Set(span, req.URL.String())
	return
}

type Option func(*roundTripper)

func WithTracer(tracer opentracing.Tracer) Option {
	if tracer == nil {
		panic("nil tracer")
	}
	return func(r *roundTripper) {
		r.tracer = tracer
	}
}

func WithRequestIgnoreFunc(ri RequestIgnoreFunc) Option {
	if ri == nil {
		ri = NoopRequestIgnore
	}

	return func(r *roundTripper) {
		r.requestIgnoreFunc = ri
	}
}

func WithRequestNameFunc(rn RequestNameFunc) Option {
	if rn == nil {
		panic("nil rn")
	}

	return func(r *roundTripper) {
		r.requestNameFunc = rn
	}
}
