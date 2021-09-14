package otgrpc

import (
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"strings"
)

var (
	ComponentTag = opentracing.Tag{Key: string(ext.Component), Value: "grpc " + grpc.Version}
)

type MDReaderWriter struct {
	metadata.MD
}

func (m MDReaderWriter) ForeachKey(handler func(key, val string) error) error {
	for k, vs := range m.MD {
		for _, v := range vs {
			if err := handler(k, v); err != nil {
				return err
			}
		}
	}
	return nil
}

func (m MDReaderWriter) Set(key, val string) {
	key = strings.ToLower(key)
	m.MD[key] = append(m.MD[key], val)
}
