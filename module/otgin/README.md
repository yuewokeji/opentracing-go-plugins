# otgin

OpenTracing plugin for [gin](https://github.com/gin-gonic/gin) web framework.

---

## Usage

---

```go
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/yuewokeji/opentracing-go-plugins/module/otgin"
)

func main() {
	server := gin.New()
	server.Use(otgin.NewMiddleware(server))

	//do something
}

```