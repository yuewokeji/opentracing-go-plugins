# othttp

---

OpenTracing plugin for the Go standard library **net/http**.

## Usage

---

### Client

```go
package main

import (
	"github.com/yuewokeji/opentracing-go-plugins/module/othttp"
	"net/http"
	"strings"
)

func main() {
	client := othttp.WrapClient(&http.Client{})
	req, err := http.NewRequest("GET", "https://github.com", strings.NewReader(""))
	if err != nil {
		panic(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	// do something
}

```

### Server

Implement later.
