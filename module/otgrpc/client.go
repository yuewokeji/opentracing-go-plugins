package otgrpc

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/yuewokeji/opentracing-go-plugins"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func NewUnaryClientInterceptor(options ...ClientOption) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, resp interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		co := &clientOptions{}
		for _, o := range options {
			o(co)
		}

		var spanCtx opentracing.SpanContext
		span := opentracing.SpanFromContext(ctx)
		if nil != span {
			spanCtx = span.Context()
		}
		span = otplugins.NonNilTracer(co.tracer).StartSpan(clientSpanName(method),
			opentracing.ChildOf(spanCtx),
			ComponentTag,
			ext.SpanKindRPCClient)
		defer span.Finish()

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		} else {
			md = md.Copy()
		}

		otplugins.NonNilTracer(co.tracer).Inject(span.Context(), opentracing.TextMap, MDReaderWriter{md})
		ctx = metadata.NewOutgoingContext(ctx, md)
		err := invoker(ctx, method, req, resp, cc, opts...)
		if err != nil {
			ext.LogError(span, err)
		}
		return err
	}
}

type clientOptions struct {
	tracer opentracing.Tracer
}

type ClientOption func(*clientOptions)

func WithClientTracer(tracer opentracing.Tracer) ClientOption {
	if nil == tracer {
		panic("nil tracer")
	}
	return func(options *clientOptions) {
		options.tracer = tracer
	}
}

func clientSpanName(method string) string {
	return "grpc.client:" + method
}
