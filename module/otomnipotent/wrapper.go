package otomnipotent

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/yuewokeji/opentracing-go-plugins"
)

// 万金油，当module里的其它包不能满足要求时，可以尝试使用wrapper
// wrapper对一个闭包进行封装，并将执行结果上报到链路监控服务

func Wrap(ctx context.Context, spanName string, doer doer) *wrapper {
	return newWrapper(ctx, spanName, doer)
}

type doer func(ctx context.Context) (interface{}, error)

func newWrapper(ctx context.Context, spanName string, doer doer) *wrapper {
	return &wrapper{
		ctx:      ctx,
		spanName: spanName,
		doer:     doer,
	}
}

type wrapper struct {
	ctx      context.Context
	spanName string
	doer     doer
}

func (w *wrapper) Do() (interface{}, error) {
	span := otplugins.NoneNilChildSpanFromContext(w.spanName, w.ctx)
	defer span.Finish()

	w.ctx = opentracing.ContextWithSpan(w.ctx, span)
	v, err := w.doer(w.ctx)
	if err != nil {
		ext.LogError(span, err)
	}

	return v, err
}
