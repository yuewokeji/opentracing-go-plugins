package otgin

import (
	"context"
	"github.com/gin-gonic/gin"
)

const TracingContextKey = "otgin:tracing:context:key"

func MustTracingContext(c *gin.Context) context.Context {
	ctx, ok := GetTracingContext(c)
	if !ok {
		return context.Background()
	}
	return ctx
}

func GetTracingContext(c *gin.Context) (context.Context, bool) {
	v, ok := c.Get(TracingContextKey)
	if !ok {
		return nil, false
	}
	ctx, ok := v.(context.Context)
	return ctx, ok
}

func SaveTracingContext(c *gin.Context, ctx context.Context) {
	c.Set(TracingContextKey, ctx)
}
