# OpenTracing-Go-Plugins

---

The plugins of [opentracing-go](https://github.com/opentracing/opentracing-go).

## Installation

---

```bash
go get -u github.com/yuewokeji/opentracing-go-plugins
```

## Configuration

---

### Initialize a tracer

Create a tracer such as [jaeger](https://github.com/uber/jaeger-client-go).

```go
package main

import (
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/transport"
	"io"
)

func initJaeger(service, url string) (opentracing.Tracer, io.Closer) {
	sender := transport.NewHTTPTransport(url)
	reporter := jaeger.NewRemoteReporter(sender, jaeger.ReporterOptions.Logger(jaeger.StdLogger))

	// samples 100% of traces
	tracer, closer := jaeger.NewTracer(service, jaeger.NewConstSampler(true), reporter)
	return tracer, closer
}

```
### Initialize the global tracer

Let's initialize the global tracer, that's because the function `opentracing.GlobalTracer()` returns a no-op tracer by default.

```go
func initGlobalTracer() io.Closer {
	// the closer can be used in shutdown hooks
	tracer, closer := initJaeger("hello-world", "https://your-reporter-url")

	opentracing.SetGlobalTracer(tracer)
	return closer
}

```

## Plugin Summary

---

1. [goroutine](module/otasyncgo/README.md)
1. [gin](module/otgin/README.md)
1. [goredis](module/otgoredisv8/README.md)
1. [gorm](module/otgorm/README.md)
1. [grpc](module/otgrpc/README.md)
1. [http client](module/othttp/README.md)
1. [omnipotent](module/otomnipotent/README.md)
