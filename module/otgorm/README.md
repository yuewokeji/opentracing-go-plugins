# otgorm

OpenTracing plugin for [gorm](https://github.com/go-gorm/gorm).

## Usage

```go
package main

import (
	"context"
	"fmt"
	"github.com/yuewokeji/opentracing-go-plugins/module/otgorm"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var orm *gorm.DB

// return db instace with context
func Orm(ctx context.Context) *gorm.DB {
	db := otgorm.Wrap(ctx, orm)
	return db.WithContext(ctx)
}

func main() {
	var (
		dsn = ""
		err error
	)
	orm, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	otgorm.RegisterCallbacks(orm)

	var count int64 = 0
	Orm(context.Background()).Raw("show tables;").Count(&count)
	if err != nil {
		panic(err)
	}
	fmt.Println(count)

	//do something
}

```