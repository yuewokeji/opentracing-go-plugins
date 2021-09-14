# otomnipotent

The omnipotent plugin of OpenTracing.

## Usage

```go
package main

import (
	"context"
	"fmt"
	"github.com/yuewokeji/opentracing-go-plugins/module/otomnipotent"
)

func doAnything(ctx context.Context) (interface{}, error) {
	return "do anything", nil
}

func main() {
	wrapper := otomnipotent.Wrap(context.Background(), "do-anything", doAnything)
	v, err := wrapper.Do()
	if err != nil {
		panic(err)
	}

	if s, ok := v.(string); ok {
		fmt.Println(s)
	}
}

```