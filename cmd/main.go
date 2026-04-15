package main

import (
	"github.com/kataras/iris/v12"
)

func main() {
	app := iris.New()

	// 启用最佳压缩效果
	app.Use(iris.Compression)

	app.Get("/", func(ctx iris.Context) {
		ctx.JSON(iris.Map{
			"message": "Hello from Iris!",
		})
	})

	go func() {
		app.Listen(":8080")
	}()
}

func setRoute() {

}
