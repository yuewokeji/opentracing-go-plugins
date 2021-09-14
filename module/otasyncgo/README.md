# otasyncgo

---

OpenTracing plugin for spawns new goroutine.

## Usage

---

```go
package main

import (
	"context"
	"fmt"
	"github.com/yuewokeji/opentracing-go-plugins/module/otasyncgo"
	"time"
)

func main() {
	fn := func(ctx context.Context) error {
		fmt.Println("from another goroutine")
		return nil
	}

	// create a goroutine and call the function fn
	otasyncgo.GoWithRecoverOnce(context.Background(), "another-goroutine", fn)

	time.Sleep(time.Second * 3)
}

```
