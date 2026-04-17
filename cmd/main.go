package main

import (
	"context"
	"minigo/config"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/sirupsen/logrus"
)

func main() {
	app := buildIris()
	setRoute(app)
	runServer(app)
	gracefulShutdown(app)
}

// 构建Iris实例
func buildIris() *iris.Application {
	app := iris.New()
	// 启用最佳压缩效果
	app.Use(iris.Compression)
	return app
}

// 设置路由
func setRoute(app *iris.Application) {
	group := app.Party("/api")

	group.Get("/hello", func(ctx iris.Context) {
		ctx.JSON(iris.Map{"message": "Hello, World!"})
	})
}

// 启动服务器
func runServer(app *iris.Application) {
	go func() {
		logrus.WithFields(logrus.Fields{"port": config.GetConfigInstance().Server.Port}).Info("Server is starting")
		if err := app.Listen(config.GetConfigInstance().Server.Port); err != nil {
			logrus.WithError(err).Error("Server error")
		}
	}()
}

// 优雅关闭服务器
func gracefulShutdown(app *iris.Application) {
	// 优雅关闭服务器	// 创建通道用于监听系统信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 等待退出信号
	sig := <-quit
	logrus.WithFields(logrus.Fields{"signal": sig}).Info("Server is shutting down")

	// 创建超时上下文，最多等待 30 秒完成正在处理的请求
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 优雅关闭服务器
	if err := app.Shutdown(ctx); err != nil {
		logrus.WithError(err).Error("Server shutdown error")
	}

	logrus.Info("Server has been shutdown")
}
