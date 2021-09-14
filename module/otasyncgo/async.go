package otasyncgo

import (
	"context"
	"errors"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/yuewokeji/opentracing-go-plugins"
	"os"
	"runtime/debug"
	"time"
)

// 本包针对goroutine的创建进行包装，并接入链路追踪
// 支持panic时，自动创建新的goroutine，并分为只执行一次、执行N次和常驻goroutine 3种方式

// Executor是goroutine的真正执行者
// ctx：一个自定义的上下文，并附加链路追踪的信息
// span：你可以通过span记录更多的信息
//
// 注意：ctx与下面的traceCtx是不同的
type Executor func(ctx context.Context) error

// 创建一个goroutine，执行fn()
// 如果发生panic，在retry次数内，会sleep一个interval周期后，重新开启新的goroutine
//
// traceCtx：与链路追踪相关的上下文，这个上下文仅用来链接到caller，随后会被销毁
//           如果你需要使用context的特性，请使用WithContext()
//           当然，你可以传递traceCtx，即WithContext(traceCtx)
// serviceName：服务名，对应链路追踪的span name
func GoWithRecover(traceCtx context.Context, serviceName string, fn Executor, options ...Option) {
	config := processConfig(options...)
	if config.retry < 0 {
		panic("retry must a positive number")
	}

	goWithRecover(traceCtx, false, serviceName, config, fn)
}

// 创建一个goroutine，执行fn()，只执行一次
func GoWithRecoverOnce(traceCtx context.Context, serviceName string, fn Executor, options ...Option) {
	config := processConfig(options...)
	config.retry = 0

	goWithRecover(traceCtx, false, serviceName, config, fn)
}

// 创建一个goroutine，执行fn()
// 如果发生panic，会sleep一个interval周期后，重新开启新的goroutine，没有次数限制
// 任何一次创建goroutine，链路追踪都会生成一个新的traceID
func GoWithRecoverResident(traceCtx context.Context, serviceName string, fn Executor, options ...Option) {
	config := processConfig(options...)
	config.retry = -1

	// 长驻的goroutine，不应该传递traceCtx，不然会导致链路无限增长
	goWithRecover(context.Background(), true, serviceName, config, fn)
}

// isBlock：是否阻断链路追踪，isBlock=true时是常驻的goroutine
func goWithRecover(traceCtx context.Context, isBlock bool, serviceName string, config *config, fn Executor) {
	go func() {
		span := opentracing.SpanFromContext(traceCtx)
		if span == nil {
			span = otplugins.NonNilTracer(config.tracer).StartSpan(spanName(serviceName))
		} else {
			span = otplugins.NonNilTracer(config.tracer).StartSpan(spanName(serviceName), opentracing.ChildOf(span.Context()))
		}
		defer span.Finish()

		span.SetTag("retry.interval", config.interval)
		span.SetTag("retry", config.retry)

		defer func() {
			if v := recover(); v != nil {
				err := errors.New(fmt.Sprintf("[Recover] panic: %v", v))

				stack := string(debug.Stack())
				os.Stderr.WriteString(err.Error() + stack)
				otplugins.Log(otplugins.LogLevelError, err.Error(), "\n", stack)

				ext.LogError(span, err, otplugins.WithLogFieldStack(stack))

				if config.interval < 0 {
					return
				}
				if config.retry == 0 {
					return
				} else if config.retry > 0 {
					config.retry--
				}
				time.Sleep(config.interval)

				if isBlock {
					goWithRecover(context.Background(), isBlock, serviceName, config, fn)
				} else {
					config.ctx = opentracing.ContextWithSpan(traceCtx, span)
					goWithRecover(config.ctx, isBlock, serviceName, config, fn)
				}
			}
		}()

		config.ctx = opentracing.ContextWithSpan(config.ctx, span)
		err := fn(config.ctx)
		if err != nil {
			ext.LogError(span, err)
		}
	}()
}

func processConfig(options ...Option) *config {
	c := defaultConfig().Clone()
	for _, o := range options {
		o(c)
	}
	return c
}

func defaultConfig() *config {
	return &config{
		interval: time.Millisecond * 100,
		retry:    0,
		ctx:      context.Background(),
		tracer: 	nil,
	}
}

type config struct {
	// 每新重新创建goroutine的时间间隔
	interval time.Duration

	// 重新创建goroutine的最大次数
	// > 0：执行指定的次数
	// = 0：只执行一次
	// =-1：永不退出
	retry int

	// ctx会传递给Executor，与traceCtx相反，代表真正的上下文内容
	ctx context.Context

	tracer opentracing.Tracer
}

func (c *config) Clone() *config {
	copied := *c
	return &copied
}

type Option func(*config)

// 设置时间间隔
func WithInterval(interval time.Duration) Option {
	if interval == 0 {
		panic("interval is zero")
	}
	return func(c *config) {
		c.interval = interval
	}
}

// 设置重试次数
func WithRetry(retry int) Option {
	return func(c *config) {
		c.retry = retry
	}
}

// 指定一个上下文
func WithContext(ctx context.Context) Option {
	return func(c *config) {
		c.ctx = ctx
	}
}

func WithTracer(tracer opentracing.Tracer) Option {
	if nil == tracer {
		panic("nil tracer")
	}
	return func(c *config) {
		c.tracer = tracer
	}
}

func spanName(serviceName string) string {
	return "async.go:" + serviceName
}
