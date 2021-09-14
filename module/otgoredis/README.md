# otgoredis

---

OpenTracing plugin for [goredis](https://github.com/go-redis/redis).

## Usage

---

```go
package main

import (
	"github.com/go-redis/redis/v8"
	"github.com/yuewokeji/opentracing-go-plugins/module/otgoredis"
)

func main() {
	r := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1",
	})
	r.AddHook(otgoredis.NewHook())

	//do something
}

```