package otplugins

import (
	"github.com/opentracing/opentracing-go/log"
)

func WithLogFieldStack(stack string) log.Field {
	return log.String("stack", stack)
}
