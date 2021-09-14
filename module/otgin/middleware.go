package otgin

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/yuewokeji/opentracing-go-plugins"
	"github.com/yuewokeji/opentracing-go-plugins/module/othttp"
	"net/http"
	"os"
	"runtime/debug"
)

var componentName = fmt.Sprintf("gin %s", gin.Version)

func NewMiddleware(engine *gin.Engine, options ...Option) gin.HandlerFunc {
	m := &middleware{
		engine:            engine,
		requestIgnoreFunc: othttp.NoopRequestIgnore,
		clientIPFunc:      othttp.GetRemoteAddr,
		recoverFunc:       nil,
	}

	for _, o := range options {
		o(m)
	}
	return m.handle
}

type middleware struct {
	engine *gin.Engine

	tracer opentracing.Tracer

	serverNamePrefix   string
	requestIgnoreFunc  othttp.RequestIgnoreFunc
	clientIPFunc       othttp.ClientIPFunc
	recoverFunc        gin.HandlerFunc
	saveTracingContext bool
}

func (m *middleware) handle(ctx *gin.Context) {
	if m.requestIgnoreFunc(ctx.Request) {
		if m.recoverFunc != nil {
			m.recoverFunc(ctx)
		}

		ctx.Next()
		return
	}

	spanCtx, _ := othttp.SpanContextFormRequest(ctx.Request, otplugins.NonNilTracer(m.tracer))
	span := otplugins.NonNilTracer(m.tracer).StartSpan(othttp.ServerSpanName(m.serverNamePrefix, ctx.Request.URL.Path), opentracing.ChildOf(spanCtx))
	defer span.Finish()

	defer func() {
		ext.PeerHostIPv4.SetString(span, m.clientIPFunc(ctx.Request))
		ext.SpanKindRPCServer.Set(span)
		ext.Component.Set(span, componentName)
		ext.HTTPMethod.Set(span, ctx.Request.Method)
		ext.HTTPUrl.Set(span, ctx.Request.URL.String())
		ext.HTTPStatusCode.Set(span, uint16(ctx.Writer.Status()))

		if v := recover(); v != nil {
			if ctx.Writer.Written() {
				ctx.Abort()
			} else {
				ctx.AbortWithStatus(http.StatusInternalServerError)
			}

			err := errors.New(fmt.Sprintf("[Recover] panic: %v", v))
			stack := string(debug.Stack())
			os.Stderr.WriteString(err.Error() + "\n" + stack)
			otplugins.Log(otplugins.LogLevelError, err.Error(), "\n", stack)

			ext.LogError(span, err, otplugins.WithLogFieldStack(stack))
		} else if othttp.IsStatusError(ctx.Writer.Status()) {
			if len(ctx.Errors) > 0 {
				for _, e := range ctx.Errors {
					ext.LogError(span, e)
				}
			} else {
				ext.LogError(span, nil)
			}
		}
	}()

	ctx.Request = othttp.RequestWithContext(ctx.Request, span)
	if m.saveTracingContext {
		SaveTracingContext(ctx, opentracing.ContextWithSpan(context.Background(), span))
	}
	ctx.Next()
}

type Option func(*middleware)

func WithRequestIgnore(r othttp.RequestIgnoreFunc) Option {
	if r == nil {
		r = othttp.NoopRequestIgnore
	}

	return func(m *middleware) {
		m.requestIgnoreFunc = r
	}
}

func WithTracer(tracer opentracing.Tracer) Option {
	return func(m *middleware) {
		if nil == tracer {
			panic("nil tracer")
		}
		m.tracer = tracer
	}
}

func WithServerNamePrefix(s string) Option {
	return func(m *middleware) {
		m.serverNamePrefix = s
	}
}

func WithClientIP(fun othttp.ClientIPFunc) Option {
	return func(m *middleware) {
		m.clientIPFunc = fun
	}
}

func WithRecover(r gin.HandlerFunc) Option {
	return func(m *middleware) {
		m.recoverFunc = r
	}
}

func WithSaveTracingContext() Option {
	return func(m *middleware) {
		m.saveTracingContext = true
	}
}
