package main

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/core/router"
)

func main() {
	app := iris.New()

	// 启用最佳压缩效果
	app.Use(iris.Compression)

	// 依赖注入
	app.ConfigureContainer(func(container *router.APIContainer) {

	})

	app.Get("/", func(ctx iris.Context) {
		ctx.JSON(iris.Map{
			"message": "Hello from Iris!",
		})
	})

	app.Listen(":8080")
}
