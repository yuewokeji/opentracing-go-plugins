package otgrpc

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/yuewokeji/opentracing-go-plugins"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"net"
	"os"
	"runtime/debug"
)

func NewUnaryServerInterceptor(options ...ServerOption) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		so := &serverOptions{}
		for _, o := range options {
			o(so)
		}

		if so.requestIgnorerFunc != nil && so.requestIgnorerFunc(info) {
			defer func() {
				if v := recover(); v != nil {
					err = status.Errorf(codes.Internal, "%s", v)
					stack := string(debug.Stack())
					os.Stderr.WriteString(err.Error() + "\n" + stack)
					otplugins.Log(otplugins.LogLevelError, err.Error(), "\n", stack)
				}
			}()

			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		}

		spanCtx, _ := otplugins.NonNilTracer(so.tracer).Extract(opentracing.TextMap, MDReaderWriter{md})
		span := opentracing.StartSpan(
			serverSpanName(info),
			ComponentTag,
			ext.RPCServerOption(spanCtx),
			ext.SpanKindRPCServer)
		defer span.Finish()

		defer func() {
			ext.PeerAddress.Set(span, GetPeerAddr(ctx))

			if v := recover(); v != nil {
				err = status.Errorf(codes.Internal, "%s", v)
				stack := string(debug.Stack())
				os.Stderr.WriteString(err.Error() + "\n" + stack)
				otplugins.Log(otplugins.LogLevelError, err.Error(), "\n", stack)

				if span != nil {
					ext.LogError(span, err, otplugins.WithLogFieldStack(stack))
				}
			}
		}()

		ctx = opentracing.ContextWithSpan(ctx, span)
		resp, err = handler(ctx, req)

		if err != nil {
			ext.LogError(span, err)
		}

		return
	}
}

type serverOptions struct {
	tracer             opentracing.Tracer
	requestIgnorerFunc RequestIgnorerFunc
}

type ServerOption func(*serverOptions)

func WithServerTracer(tracer opentracing.Tracer) ServerOption {
	if tracer == nil {
		panic("nil tracer")
	}
	return func(options *serverOptions) {
		options.tracer = tracer
	}
}

type RequestIgnorerFunc func(*grpc.UnaryServerInfo) bool

func WithServerRequestIgnorer(r RequestIgnorerFunc) ServerOption {
	return func(options *serverOptions) {
		options.requestIgnorerFunc = r
	}
}

func serverSpanName(info *grpc.UnaryServerInfo) string {
	return "grpc.server:" + info.FullMethod
}

func GetPeerAddr(ctx context.Context) string {
	var addr string
	if pr, ok := peer.FromContext(ctx); ok {
		if tcpAddr, ok := pr.Addr.(*net.TCPAddr); ok {
			addr = tcpAddr.IP.String()
		} else {
			addr = pr.Addr.String()
		}
	}
	return addr
}
