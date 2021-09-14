# otgorm

OpenTracing plugin for [gorm](https://github.com/go-gorm/gorm).

## Usage

```go
package main

import (
	"github.com/yuewokeji/opentracing-go-plugins/module/otgorm"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := ""
	orm, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	otgorm.RegisterCallbacks(orm)

	//do something
}

```