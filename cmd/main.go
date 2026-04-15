package main

import (
	"minigo/config"

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

	err := app.Listen(config.GetConfigInstance().Server.Port)
	if err != nil {
		return
	}
}

func setRoute() {

}
