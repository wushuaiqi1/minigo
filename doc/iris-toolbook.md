# Iris 框架工具书

> 基于 Iris 官方文档整理，面向 Go Web 开发新手的入门手册。
>
> 目标：读完后能独立搭建一个 Iris Web/API 项目，并理解常见功能的基本用法。

---

## 1. Iris 是什么

Iris 是一个 Go 语言 Web 框架，强调高性能、完整功能和良好的开发体验。它既适合写简单的 REST API，也适合构建带模板渲染、会话管理、中间件、MVC、WebSocket、API 版本控制等能力的中大型 Web 服务。

Iris 官方强调的核心能力包括：

- 高性能路由
- 中间件机制
- 内建 Session
- 依赖注入
- MVC 支持
- WebSocket 支持
- 多模板引擎支持
- API 版本控制
- 丰富的请求读取与响应输出能力
- 较好的可测试性

适合使用 Iris 的场景：

- 快速开发 RESTful API
- 构建后台管理系统
- 服务端渲染网站
- 需要中间件、鉴权、会话、模板、WebSocket 的项目
- 想用更强分层能力组织 Go Web 项目

---

## 2. 安装与环境准备

### 2.1 前置条件

你需要先安装 Go。官方资料显示 Iris 运行需要较新的 Go 版本，实际开发时建议使用当前官方支持的稳定版 Go。

### 2.2 安装 Iris

```bash
go get github.com/kataras/iris/v12@latest
```

如果是新项目，通常这样初始化：

```bash
mkdir myapp
cd myapp
go mod init myapp
go get github.com/kataras/iris/v12@latest
go mod tidy
```

---

## 3. 第一个 Iris 应用

```go
package main

import "github.com/kataras/iris/v12"

func main() {
	app := iris.New()

	app.Get("/", func(ctx iris.Context) {
		ctx.JSON(iris.Map{
			"message": "Hello from Iris!",
		})
	})

	app.Listen(":8080")
}
```

运行：

```bash
go run .
```

访问：

```text
http://localhost:8080
```

你会得到一个 JSON 响应。

这段代码里最重要的几个概念：

- `iris.New()`：创建应用实例
- `app.Get()`：注册 GET 路由
- `iris.Context`：请求上下文，处理请求和响应的核心对象
- `ctx.JSON()`：返回 JSON
- `app.Listen(":8080")`：启动 HTTP 服务

---

## 4. Iris 的核心结构理解

在 Iris 中，可以把一个请求的处理流程理解为：

`请求 -> 路由匹配 -> 中间件 -> 处理函数 -> 响应输出`

你需要重点理解 4 个对象：

- `Application`：整个应用
- `Party`：路由分组
- `Context`：一次请求的上下文
- `Handler`：处理函数或中间件

典型项目里常见的职责分层是：

`Router -> Controller -> Service -> Repository -> Database`

虽然 Iris 不强制项目结构，但这是很常见、也很适合新手入门的一种方式。

---

## 5. 路由系统

### 5.1 基本路由

```go
app.Get("/users", listUsers)
app.Post("/users", createUser)
app.Put("/users/{id:int}", updateUser)
app.Delete("/users/{id:int}", deleteUser)
```

Iris 支持常见 HTTP 方法：

- `Get`
- `Post`
- `Put`
- `Delete`
- `Patch`
- `Head`
- `Options`

### 5.2 路由参数

```go
app.Get("/users/{id:int}", func(ctx iris.Context) {
	id, err := ctx.Params().GetInt("id")
	if err != nil {
		ctx.StopWithStatus(iris.StatusBadRequest)
		return
	}

	ctx.JSON(iris.Map{"id": id})
})
```

常见参数写法：

- `{id:int}`：整数参数
- `{id:uint64}`：无符号整数
- `{name:string}`：字符串

路由参数的优势是：

- 路由定义更清晰
- 参数类型约束更直接
- 能提前减少非法请求

### 5.3 路由分组 Party

```go
api := app.Party("/api")
v1 := api.Party("/v1")

v1.Get("/users", listUsers)
v1.Post("/users", createUser)
```

`Party` 可以理解为带公共前缀、公共中间件、公共配置的一组路由。实际项目中非常常用。

---

## 6. Context：请求上下文

`iris.Context` 是你写 Iris 时最常接触的对象。它负责：

- 读取请求数据
- 获取路由参数
- 获取查询参数
- 获取表单数据
- 写入响应
- 设置状态码
- 设置响应头
- 在中间件链中控制流程

常见用法：

```go
app.Get("/hello", func(ctx iris.Context) {
	name := ctx.URLParamDefault("name", "world")
	ctx.StatusCode(iris.StatusOK)
	ctx.Header("X-App", "iris-demo")
	ctx.WriteString("hello " + name)
})
```

---

## 7. 读取请求数据

### 7.1 读取 URL Query 参数

```go
app.Get("/search", func(ctx iris.Context) {
	keyword := ctx.URLParam("keyword")
	page, _ := ctx.URLParamInt("page")

	ctx.JSON(iris.Map{
		"keyword": keyword,
		"page":    page,
	})
})
```

请求示例：

```text
/search?keyword=iris&page=2
```

### 7.2 读取 JSON 请求体

```go
type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

app.Post("/users", func(ctx iris.Context) {
	var req CreateUserRequest
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.StopWithError(iris.StatusBadRequest, err)
		return
	}

	ctx.JSON(iris.Map{
		"message": "created",
		"data":    req,
	})
})
```

### 7.3 读取 XML / YAML / Form / Query

Iris 不只支持 JSON，也支持多种请求数据绑定，例如：

- `ctx.ReadXML(&v)`
- `ctx.ReadYAML(&v)`
- 表单与 URL Query 的绑定
- 在依赖注入模式下，直接把结构体作为处理函数参数

对新手来说，最常用的是：

- API 场景用 `ReadJSON`
- 表单提交用 `ctx.FormValue`
- 查询筛选用 `ctx.URLParam`

### 7.4 请求校验

Iris 官方文档给出了配合 `go-playground/validator` 的校验写法：

```go
type User struct {
	Name     string `json:"name" validate:"required,min=2"`
	Email    string `json:"email" validate:"required,email"`
	Age      int    `json:"age" validate:"required,gte=0,lte=130"`
	Password string `json:"password" validate:"required,min=8"`
}
```

示例：

```go
import "github.com/go-playground/validator/v10"

validate := validator.New()

app.Post("/users", func(ctx iris.Context) {
	var user User
	if err := ctx.ReadJSON(&user); err != nil {
		ctx.StopWithError(iris.StatusBadRequest, err)
		return
	}

	if err := validate.Struct(user); err != nil {
		ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().
			Title("Validation Error").
			DetailErr(err))
		return
	}

	ctx.JSON(iris.Map{"message": "ok"})
})
```

建议：

- 请求结构体和数据库模型分开
- 校验失败统一返回 400
- 对外隐藏底层实现细节

---

## 8. 响应输出

Iris 支持非常丰富的响应类型。

### 8.1 返回文本

```go
app.Get("/text", func(ctx iris.Context) {
	ctx.Text("Hello World")
})
```

### 8.2 返回 JSON

```go
app.Get("/json", func(ctx iris.Context) {
	ctx.JSON(iris.Map{
		"code":    0,
		"message": "success",
	})
})
```

### 8.3 返回格式化 JSON

```go
app.Get("/pretty", func(ctx iris.Context) {
	ctx.JSON(iris.Map{
		"name": "iris",
		"type": "framework",
	}, iris.JSON{
		Indent: "  ",
	})
})
```

### 8.4 返回 XML

```go
type Person struct {
	ID   int    `xml:"id"`
	Name string `xml:"name"`
}

app.Get("/xml", func(ctx iris.Context) {
	ctx.XML(Person{ID: 1, Name: "Tom"})
})
```

### 8.5 设置状态码与响应头

```go
app.Get("/status", func(ctx iris.Context) {
	ctx.StatusCode(iris.StatusCreated)
	ctx.Header("X-Trace-ID", "abc123")
	ctx.JSON(iris.Map{"message": "created"})
})
```

新手建议统一响应格式，例如：

```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

这样前后端联调会更顺畅。

---

## 9. 中间件 Middleware

中间件是 Iris 的核心能力之一。一个中间件本质上就是：

```go
func(ctx iris.Context)
```

它常用于：

- 记录日志
- 鉴权
- 统一异常处理
- 跨域处理
- 请求追踪
- 性能统计

Iris 的中间件体系比很多 Go Web 框架更完整。它不只有“请求前执行”的一个入口，而是提供了多种注册层级和执行阶段，让你可以精确控制：

- 在路由匹配前还是匹配后执行
- 只对正常路由生效，还是对错误路由也生效
- 在主处理器之前执行，还是在主处理器之后收尾
- 对整个应用生效，还是只对某个路由组生效

### 9.1 基本中间件示例

```go
func authMiddleware(ctx iris.Context) {
	token := ctx.GetHeader("Authorization")
	if token == "" {
		ctx.StopWithStatus(iris.StatusUnauthorized)
		return
	}

	ctx.Next()
}

app.Use(authMiddleware)
```

### 9.2 分组中间件

```go
admin := app.Party("/admin")
admin.Use(authMiddleware)

admin.Get("/dashboard", func(ctx iris.Context) {
	ctx.WriteString("admin dashboard")
})
```

### 9.3 Iris 提供的多种中间件注册方式

根据 Iris 官方文档，常见的中间件注册方式主要有这些：

- `Application.WrapRouter`
- `Party.UseRouter`
- `Application.UseGlobal`
- `Party.Use`
- `Party.UseOnce`
- `Party.UseError`
- `Party.Done`
- `Application.DoneGlobal`

下面逐个解释。

#### 9.3.1 `Application.WrapRouter`

这是最底层的拦截方式，作用在整个应用最前面，甚至早于 Iris 自己的路由处理。

它不是普通的 `iris.Handler`，而是一个更底层的包装函数，适合做：

- 非常早期的安全检查
- 最底层 CORS 处理
- 请求重写
- 在进入 Iris 路由器前修改行为

它的特点是：

- 作用范围是整个应用
- 执行时机最早
- 执行顺序是“后注册先执行”
- 能直接接触 `http.ResponseWriter` 和 `*http.Request`

新手建议：

- 如果你只是做常规鉴权、日志、请求计时，不要优先用它
- 只有当你明确需要“在路由器之前”拦截请求时再使用

#### 9.3.2 `Party.UseRouter`

`UseRouter` 作用在某个 `Party` 路由组前缀下，并且它会在路由匹配阶段参与执行。

它的特点是：

- 对指定 `Party` 及其子路由组生效
- 即使路由最终没有匹配成功，它也可能执行
- 执行时机早于 `UseGlobal` 和 `Use`
- 很适合做“路径级过滤”

适用场景：

- 某个前缀下的统一访问限制
- 对某组接口做早期过滤
- 需要在错误页也参与处理的逻辑

示例：

```go
api := app.Party("/api")
api.UseRouter(func(ctx iris.Context) {
	println("before router matched under /api")
	ctx.Next()
})
```

#### 9.3.3 `Application.UseGlobal`

`UseGlobal` 是应用级中间件，对整个应用中的已有路由和未来路由都生效，也能作用于通过 `OnErrorCode` 注册的错误处理页面。

它的特点是：

- 作用范围是整个应用
- 所有 Party、子域名、后续新增路由都能继承
- 执行时机在 `UseRouter` 之后、`Use` 之前
- 非常适合放全局公共逻辑

适合做：

- 全局日志
- trace id 注入
- 请求耗时统计
- 全局恢复处理
- 公共 header 设置

示例：

```go
app.UseGlobal(func(ctx iris.Context) {
	ctx.Header("X-App", "iris-demo")
	ctx.Next()
})
```

#### 9.3.4 `Party.Use`

这是最常见、也是新手最应该优先掌握的中间件注册方式。

它的特点是：

- 对某个 `Party` 及其子路由生效
- 执行在主处理器之前
- 只对“正常匹配到的路由”执行
- 对纯错误路由不会生效

适合做：

- 登录校验
- 权限校验
- 参数预处理
- 当前用户注入

示例：

```go
admin := app.Party("/admin")
admin.Use(func(ctx iris.Context) {
	if ctx.GetHeader("Authorization") == "" {
		ctx.StopWithStatus(iris.StatusUnauthorized)
		return
	}

	ctx.Next()
})
```

#### 9.3.5 `Party.UseOnce`

`UseOnce` 和 `Use` 类似，但它的设计目标是避免重复注册相同中间件。

适合场景：

- 多个模块初始化时可能重复调用注册逻辑
- 想避免同一个中间件被重复挂载

对新手来说：

- 日常开发使用频率不如 `Use`
- 但在大型项目或插件式注册中会更稳妥

#### 9.3.6 `Party.UseError`

`UseError` 专门服务于错误处理链。

它的特点是：

- 只在 HTTP 错误场景执行
- 例如 `404 Not Found`、`401 Unauthorized` 等
- 对普通成功路由不会执行
- 对子路由组同样生效

适合做：

- 统一错误页
- 错误日志
- 错误响应格式包装

示例：

```go
app.UseError(func(ctx iris.Context) {
	println("error middleware:", ctx.GetStatusCode())
	ctx.Next()
})

app.OnErrorCode(iris.StatusNotFound, func(ctx iris.Context) {
	ctx.JSON(iris.Map{
		"code":    404,
		"message": "resource not found",
	})
})
```

注意：

- `Use` 不处理错误页
- 错误链应交给 `UseError`

#### 9.3.7 `Party.Done`

`Done` 用于在主路由处理器之后执行“收尾逻辑”。

它的特点是：

- 在主 handler 之后执行
- 对某个 `Party` 及其子路由生效
- 默认需要前面的 handler 调用了 `ctx.Next()`，它才能继续执行

适合做：

- 请求结束日志
- 收尾统计
- 响应后埋点
- 清理临时状态

示例：

```go
app.Done(func(ctx iris.Context) {
	println("request finished:", ctx.Path())
})
```

#### 9.3.8 `Application.DoneGlobal`

这是全局版的 `Done`。

它的特点是：

- 作用于整个应用
- 执行时机在整个链路靠后的位置
- 适合做全局收尾工作

适合做：

- 统一访问日志落盘
- 全局链路统计
- 统一资源清理

### 9.4 各种注册方式的对照理解

你可以把 Iris 中间件粗略理解成三层：

1. 最前层：`WrapRouter`
2. 前置层：`UseRouter`、`UseGlobal`、`Use`
3. 错误与收尾层：`UseError`、`Done`、`DoneGlobal`

如果只给新手一个实用建议：

- 90% 的日常开发先用 `Use` 和 `UseGlobal`
- 需要错误处理时用 `UseError`
- 需要收尾逻辑时用 `Done`
- 只有真的要抢在路由器前面时才用 `WrapRouter`

### 9.5 中间件执行顺序

根据官方文档，典型的正常请求链路可以理解为：

`WrapRouter -> UseRouter -> UseGlobal -> Use -> 主处理器 -> Done -> DoneGlobal`

如果是错误请求链路，则通常会变成：

`WrapRouter -> UseRouter -> UseGlobal -> UseError -> 错误处理器`

注意两点：

- `Use` 主要处理正常匹配路由
- `UseError` 主要处理错误链

### 9.6 中间件执行流程控制

这是 Iris 中间件最关键的部分。真正决定链路怎么走的，不是“注册了什么”，而是“你在 handler 里如何控制执行”。

最核心的方法是：

- `ctx.Next()`
- `ctx.StopExecution()`
- `ctx.StopWithStatus(code)`
- `ctx.StopWithText(code, text)`
- `ctx.StopWithError(code, err)`

#### 9.6.1 `ctx.Next()`

`ctx.Next()` 的意思是：继续执行下一个处理器。

如果你不调用它，请求链通常不会继续向后传递。

示例：

```go
func logger(ctx iris.Context) {
	println("before")
	ctx.Next()
	println("after")
}
```

这里的执行效果是：

1. 先打印 `before`
2. 进入下一个 handler
3. 下一个 handler 返回后，再打印 `after`

这说明 `ctx.Next()` 不只是“继续向后走”，它还让你可以写出“前置 + 后置”结构。

#### 9.6.2 不调用 `ctx.Next()`

如果一个中间件没有调用 `ctx.Next()`，那么后面的处理器通常不会执行。

这经常用于：

- 未登录直接拦截
- 参数不合法直接返回
- 命中黑名单直接拒绝

示例：

```go
func auth(ctx iris.Context) {
	if ctx.GetHeader("Authorization") == "" {
		ctx.StatusCode(iris.StatusUnauthorized)
		ctx.JSON(iris.Map{"message": "unauthorized"})
		return
	}

	ctx.Next()
}
```

#### 9.6.3 `ctx.StopExecution()`

`ctx.StopExecution()` 用于显式停止后续执行。

适合在已经写好了响应，但还想明确告诉框架“链路到此结束”时使用。

示例：

```go
func deny(ctx iris.Context) {
	ctx.StatusCode(iris.StatusForbidden)
	ctx.WriteString("forbidden")
	ctx.StopExecution()
}
```

#### 9.6.4 `StopWith...` 系列

Iris 提供了一组更方便的“停止并响应”方法：

- `ctx.StopWithStatus(code)`
- `ctx.StopWithText(code, text)`
- `ctx.StopWithError(code, err)`

它们适合直接用于失败分支，代码更简洁。

示例：

```go
func requireJSON(ctx iris.Context) {
	if ctx.GetHeader("Content-Type") != "application/json" {
		ctx.StopWithText(iris.StatusUnsupportedMediaType, "JSON required")
		return
	}

	ctx.Next()
}
```

### 9.7 前置、中置、后置写法

中间件不一定只能写在 `app.Use(...)` 里，也可以直接挂在某一条路由上，形成非常直观的链路。

```go
func before(ctx iris.Context) {
	ctx.Values().Set("request_id", "abc123")
	println("before")
	ctx.Next()
}

func mainHandler(ctx iris.Context) {
	println("main")
	ctx.WriteString("ok")
	ctx.Next()
}

func after(ctx iris.Context) {
	println("after")
}

app.Get("/demo", before, mainHandler, after)
```

上面这个例子里：

- `before` 是前置处理
- `mainHandler` 是主业务处理
- `after` 是后置处理

要注意：

- 如果 `mainHandler` 不调用 `ctx.Next()`，那么 `after` 不会执行
- 这也是 `Done` 风格处理常常需要关注的地方

### 9.8 如何在中间件之间传递数据

Iris 官方推荐通过 `ctx.Values()` 在同一次请求的处理链中共享数据。

示例：

```go
func currentUserMiddleware(ctx iris.Context) {
	ctx.Values().Set("current_user_id", 1001)
	ctx.Next()
}

func profileHandler(ctx iris.Context) {
	userID := ctx.Values().GetIntDefault("current_user_id", 0)
	ctx.JSON(iris.Map{"user_id": userID})
}
```

适合传递：

- 当前用户信息
- trace id
- 权限判断结果
- 预处理后的业务上下文

建议：

- 共享的数据命名统一
- 不要在 `Values` 里塞太多杂乱数据
- 对关键字段做好类型约定

### 9.9 Execution Rules：强制执行 Begin/Done

官方文档还提到 `ExecutionRules`，它可以修改默认执行规则，强制某些 Begin/Done 类中间件即使没有 `ctx.Next()` 也继续执行。

示例：

```go
app.SetExecutionRules(iris.ExecutionRules{
	Done: iris.ExecutionOptions{Force: true},
})
```

这意味着：

- 即使主 handler 没有调用 `ctx.Next()`
- `Done` 阶段的 handler 也可以被强制执行

适用场景：

- 你必须保证收尾逻辑总会执行
- 例如统一审计、监控、资源清理

但新手要注意：

- 一旦改动执行规则，理解链路会更复杂
- 团队里最好统一规范，不要有人默认规则、有人强制规则混用

### 9.10 一个完整的执行流程示例

```go
app.UseGlobal(func(ctx iris.Context) {
	println("global before")
	ctx.Next()
	println("global after")
})

app.Use(func(ctx iris.Context) {
	println("use before")
	ctx.Next()
	println("use after")
})

app.Done(func(ctx iris.Context) {
	println("done")
})

app.Get("/hello", func(ctx iris.Context) {
	println("handler")
	ctx.WriteString("hello")
	ctx.Next()
})
```

一个可能的执行顺序是：

1. `global before`
2. `use before`
3. `handler`
4. `done`
5. `use after`
6. `global after`

这能帮助你理解：

- `Next()` 会进入后续处理器
- 后续处理器执行完成后，会回到当前 handler 的 `Next()` 之后
- `Done` 处于主处理器之后的阶段

### 9.11 新手使用中间件的实战建议

- 全局日志、trace、公共 header 用 `UseGlobal`
- 登录校验、权限校验优先用 `Party.Use`
- 错误格式统一用 `UseError` + `OnErrorCode`
- 请求结束统计可以考虑 `Done` 或 `DoneGlobal`
- 共享数据用 `ctx.Values()`
- 只有确实要在 Iris 路由器之前拦截时才用 `WrapRouter`

### 9.12 中间件执行的关键点总结

- 调用 `ctx.Next()` 才会进入下一个处理器
- 不调用 `ctx.Next()` 可以中断后续执行
- `ctx.Next()` 之后的代码会在后续处理器返回后继续执行
- 错误链和正常链是分开的，`Use` 不等于 `UseError`
- `Done` 默认依赖前面 handler 正确调用 `ctx.Next()`
- 适合把公共逻辑从业务处理函数中抽离出来
- 实际开发优先掌握 `UseGlobal`、`Use`、`UseError`、`Done`

---

## 10. 模板渲染 View

如果你不是只写 API，而是要返回 HTML 页面，那么 Iris 的模板系统会很有用。

### 10.1 注册 HTML 模板

```go
tmpl := iris.HTML("./views", ".html")
tmpl.Reload(true)

app.RegisterView(tmpl)
```

### 10.2 渲染页面

```go
app.Get("/", func(ctx iris.Context) {
	ctx.View("index.html", iris.Map{
		"title":   "首页",
		"message": "欢迎来到 Iris",
	})
})
```

### 10.3 使用 Layout

```go
tmpl := iris.HTML("./views", ".html")
tmpl.Layout("layouts/main.html")
tmpl.Reload(true)
app.RegisterView(tmpl)
```

也可以在单个路由中覆盖布局：

```go
app.Get("/custom", func(ctx iris.Context) {
	ctx.ViewLayout("layouts/custom.html")
	ctx.View("page.html")
})
```

### 10.4 自定义模板函数

```go
tmpl := iris.HTML("./views", ".html")
tmpl.AddFunc("greet", func(name string) string {
	return "Hello " + name
})
app.RegisterView(tmpl)
```

模板中可直接使用：

```html
{{ greet "Iris" }}
```

### 10.5 支持的模板引擎

官方文档展示了多种模板引擎支持，包括：

- 标准 HTML 模板
- Django 模板
- Pug 模板

对新手最推荐的是标准 HTML 模板，因为它最稳定、最接近 Go 原生生态。

---

## 11. 静态文件与文件服务

前端静态资源通常包括：

- CSS
- JavaScript
- 图片
- 上传文件

Iris 可以通过 `HandleDir` 快速托管静态资源目录。

### 11.1 基本静态目录

```go
app.HandleDir("/static", iris.Dir("./assets"))
```

访问示例：

```text
http://localhost:8080/static/app.css
```

### 11.2 开启目录选项

```go
app.HandleDir("/files", iris.Dir("./uploads"), iris.DirOptions{
	Compress: true,
	ShowList: true,
})
```

### 11.3 文件下载

```go
app.Get("/download", func(ctx iris.Context) {
	ctx.SendFile("./files/report.pdf", "report.pdf")
})
```

适用场景：

- 管理后台附件下载
- 图片资源托管
- 上传结果预览

---

## 12. 压缩 Compression

Iris 可以通过一行代码开启压缩能力：

```go
app.Use(iris.Compression)
```

它可以对响应内容进行压缩，也支持读取压缩请求体。官方文档列出的内建压缩算法包括：

- `gzip`
- `deflate`
- `br`
- `snappy`

对新手的建议：

- API 服务默认可以开启压缩
- 大量文本响应场景收益明显
- 文件下载要结合实际情况测试

---

## 13. Session 与 Cookie

Session 适合保存用户登录态、访问次数、临时会话数据等。

### 13.1 创建 Session 管理器

```go
import "github.com/kataras/iris/v12/sessions"

sess := sessions.New(sessions.Config{
	Cookie: "mysession",
})

app.Use(sess.Handler())
```

### 13.2 使用 Session

```go
app.Get("/visit", func(ctx iris.Context) {
	session := sess.Start(ctx)
	count := session.Increment("visits", 1)

	ctx.JSON(iris.Map{
		"visits": count,
	})
})
```

### 13.3 安全配置建议

官方安全文档特别提到 Session/Cookie 的配置应关注：

- Cookie 名称
- 子域共享策略
- `SameSite`
- 域名范围
- 是否允许回收与持久化

新手至少要注意：

- 生产环境使用 HTTPS
- Cookie 配置结合域名策略
- 登录态不要只依赖前端判断

---

## 14. 依赖注入 Dependency Injection

依赖注入是 Iris 很有辨识度的一项能力。你可以把“请求需要什么输入”“业务逻辑依赖什么服务”“响应返回什么结果”直接写在函数签名里，让框架在运行时自动完成绑定和调用。

对服务端开发初学者来说，可以先把它理解成一句话：

- 传统写法：先拿 `ctx`，再手动从 `ctx` 里读参数、读 JSON、调服务、写响应
- Iris DI 写法：把这些需要的东西直接写到函数参数和返回值中，Iris 帮你做大部分样板代码

### 14.1 先理解：Iris 的 DI 在解决什么问题

先看一个不使用 DI 的常见写法：

```go
app.Put("/user/{id:uint64}", func(ctx iris.Context) {
	id, _ := ctx.Params().GetUint64("id")

	var req UpdateUserRequest
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.StopWithStatus(iris.StatusBadRequest)
		return
	}

	resp := UpdateUserResponse{
		ID:      id,
		Message: "User updated successfully",
	}

	ctx.JSON(resp)
})
```

这个写法并没有错，但你会发现有很多重复劳动：

- 手动读取路由参数
- 手动解析请求体
- 手动处理响应输出
- 每个处理函数都要写一遍类似流程

而 Iris 的 DI 想做的是：把这些“通用搬运工作”交给框架，让处理函数更接近业务本身。

### 14.2 最小示例：从函数签名读懂 DI

```go
type UpdateUserRequest struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

type UpdateUserResponse struct {
	ID      uint64 `json:"id"`
	Message string `json:"message"`
}

func updateUser(id uint64, input UpdateUserRequest) UpdateUserResponse {
	return UpdateUserResponse{
		ID:      id,
		Message: "User updated successfully",
	}
}

func main() {
	app := iris.New()

	app.Party("/user").ConfigureContainer(func(api *iris.APIContainer) {
		api.Put("/{id:uint64}", updateUser)
	})

	app.Listen(":8080")
}
```

这里有三件关键的事：

- `id uint64` 会自动绑定到路由参数 `{id:uint64}`
- `input UpdateUserRequest` 会自动绑定客户端发送的数据
- 返回值 `UpdateUserResponse` 会自动作为响应写回客户端，通常会序列化为 JSON

换句话说，Iris 会先分析你的函数签名，再决定“每个参数从哪里来”“返回值该怎么发出去”。

### 14.3 工作原理：Iris 是怎么把参数注入进去的

官方文档里把这套能力挂在 `ConfigureContainer` / `APIContainer` 上。它的大体机制可以理解为以下几步：

1. 在应用启动阶段，你通过 `ConfigureContainer` 注册路由和依赖
2. Iris 预先分析处理函数的参数列表和返回值列表
3. 请求真正到来时，Iris 按参数类型逐个准备所需值
4. 所有参数准备好后，调用你的处理函数
5. 再根据返回值类型，把结果写入 HTTP 响应

对初学者来说，这里最重要的认识不是“内部实现细节”，而是两条规则：

- Iris 的 DI 首先依赖“类型”和“路由定义”来匹配参数
- 这不是随便猜，它是按照固定规则去绑定输入和处理输出的

因此，函数签名在 Iris 中非常重要。它不只是 Go 函数签名，也是在描述一个 HTTP 接口的输入契约。

### 14.4 输入参数从哪里来

在 Iris 的 DI 模式下，一个处理函数的输入参数，通常来自下面几类来源。

#### 1. 路由参数

例如：

```go
func getUser(id uint64) iris.Map {
	return iris.Map{"id": id}
}
```

配合：

```go
api.Get("/{id:uint64}", getUser)
```

则 `id uint64` 会绑定到 URL 中的 `{id:uint64}`。

常见场景：

- `/users/{id:uint64}` 里的用户 ID
- `/articles/{slug:string}` 里的文章标识
- `/orders/{orderID:int}` 里的订单号

#### 2. 请求体或请求数据

当参数是一个结构体时，Iris 会尝试把客户端提交的数据绑定到这个结构体中。官方示例里提到，常见输入可以来自：

- JSON
- 表单
- URL Query
- 其他请求数据格式

示例：

```go
type CreateUserRequest struct {
	Firstname string `json:"firstname" form:"firstname" url:"firstname"`
	Lastname  string `json:"lastname" form:"lastname" url:"lastname"`
}

func createUser(req CreateUserRequest) iris.Map {
	return iris.Map{
		"firstname": req.Firstname,
		"lastname":  req.Lastname,
	}
}
```

如果客户端发送：

```json
{
  "firstname": "Tom",
  "lastname": "Lee"
}
```

那么 `req` 就会被自动填充。

#### 3. 上下文和 HTTP 底层对象

当你确实需要直接操作请求和响应，也可以把底层对象作为参数：

- `iris.Context`
- `*http.Request`
- `http.ResponseWriter`

例如：

```go
func hello(ctx iris.Context, r *http.Request) string {
	ctx.Header("X-Request-Path", r.URL.Path)
	return "ok"
}
```

这说明 Iris 的 DI 并不是“强制你不用 `ctx`”，而是“你只在需要时才拿 `ctx`”。

#### 4. 已注册的服务对象

这是服务端开发里最实用的一部分。你可以把数据库连接、Service、Repository、配置对象、日志对象等注册到容器里，然后让处理函数直接接收它们。

### 14.5 注册自定义依赖：把 Service 注入处理函数

先看一个最简单的静态依赖注册。

```go
type UserService struct{}

func (s *UserService) Create(name string) string {
	return "created: " + name
}

type CreateUserRequest struct {
	Firstname string `json:"firstname"`
}

func createUser(svc *UserService, req CreateUserRequest) iris.Map {
	return iris.Map{
		"message": svc.Create(req.Firstname),
	}
}

func main() {
	app := iris.New()

	app.ConfigureContainer(func(api *iris.APIContainer) {
		api.RegisterDependency(&UserService{})
		api.Post("/users", createUser)
	})

	app.Listen(":8080")
}
```

这里的执行逻辑可以理解为：

- `req CreateUserRequest` 来自客户端请求
- `svc *UserService` 来自容器中预先注册的依赖
- Iris 在调用 `createUser` 前，会把这两个参数都准备好

这很适合把项目分层成：

- Handler：收发 HTTP
- Service：处理业务逻辑
- Repository：访问数据库

这样处理函数不会塞满数据库查询、参数解析和响应拼装代码，可读性会好很多。

### 14.6 依赖不只能是值，也可以是函数

根据官方文档，Iris 注册的 dependency 不一定非得是一个固定对象，也可以是一个函数。这个函数本身也能依赖其他已注册对象，然后返回新的依赖值。

这意味着你可以做“依赖工厂”。

例如：

```go
type AppConfig struct {
	AppName string
}

type UserService struct {
	AppName string
}

func newUserService(cfg AppConfig) *UserService {
	return &UserService{AppName: cfg.AppName}
}

func createUser(svc *UserService) string {
	return "service from: " + svc.AppName
}

app.ConfigureContainer(func(api *iris.APIContainer) {
	api.RegisterDependency(AppConfig{AppName: "demo"})
	api.RegisterDependency(newUserService)
	api.Post("/users", createUser)
})
```

可以把它理解成：

- 先注册基础依赖 `AppConfig`
- 再注册一个“构造函数依赖” `newUserService`
- 当处理函数需要 `*UserService` 时，Iris 会知道如何把它构造出来

这种方式在真实项目里很常见，特别适合：

- Service 依赖配置对象
- Repository 依赖数据库连接
- 权限对象依赖当前请求上下文

### 14.7 请求期间动态注册依赖

Iris 不仅支持在应用启动时注册依赖，还支持在请求处理中注册“本次请求专属”的依赖。官方提供了：

- `ctx.RegisterDependency(...)`

这类用法适合中间件先从请求里解析出一些信息，再把它作为依赖传给后续处理函数。

例如，一个认证中间件先解析当前用户角色：

```go
type Role struct {
	Name string
}

func RoleMiddleware(ctx iris.Context) {
	role := Role{Name: "admin"}
	ctx.RegisterDependency(role)
	ctx.Next()
}

func dashboard(role Role) string {
	return "role: " + role.Name
}
```

这样后续 handler 就不需要自己再去 `ctx.Values()` 里手动取角色值。

这对初学者有两个好处：

- 中间件负责准备通用上下文数据
- 业务处理函数只关心它真正需要的对象

### 14.8 返回值不只是普通返回值，它也是响应定义

Iris 的 DI 不只体现在输入参数上，也体现在返回值上。很多情况下你不需要手写：

```go
ctx.JSON(...)
```

而是可以直接返回结果对象：

```go
type Response struct {
	ID      uint64 `json:"id,omitempty"`
	Message string `json:"message"`
}

func getUser(id uint64) Response {
	return Response{
		ID:      id,
		Message: "ok",
	}
}
```

这会让 handler 更像“纯函数”：输入明确，输出明确。

对初学者来说，这样的代码有一个很明显的优点：

- 看函数签名就能大致知道接口收什么、回什么

### 14.9 Preflight：在响应发出前做最后处理

官方 DI 文档里还有一个很实用的机制：如果你的返回结构体实现了 `Preflight(iris.Context) error` 方法，那么 Iris 会在发送响应前自动调用它。

这适合做：

- 动态设置状态码
- 补充响应字段
- 在响应发出前统一调整输出

示例：

```go
type DeleteUserResponse struct {
	Message   string `json:"message"`
	Code      int    `json:"code"`
	Timestamp int64  `json:"timestamp,omitempty"`
}

func (r *DeleteUserResponse) Preflight(ctx iris.Context) error {
	if r.Code > 0 {
		ctx.StatusCode(r.Code)
	}

	r.Timestamp = time.Now().Unix()
	return nil
}

func deleteUser(id uint64) *DeleteUserResponse {
	return &DeleteUserResponse{
		Message: "user has been marked for deletion",
		Code:    iris.StatusAccepted,
	}
}
```

这样做的好处是：状态码和响应体的定制逻辑可以跟响应类型绑定在一起，而不是散落在多个 handler 里。

### 14.10 错误处理与结果处理

在 `ConfigureContainer` 上，官方还提供了两个很值得知道的能力：

- `OnError`
- `UseResultHandler`

它们分别解决两个问题。

#### 1. `OnError`

用于统一处理依赖注入过程或处理函数执行过程中返回的错误。

适合：

- 统一记录日志
- 统一返回错误 JSON
- 避免每个路由都单独写一遍错误响应

#### 2. `UseResultHandler`

用于接管或扩展“返回值如何变成响应”的逻辑。

适合：

- 统一包装响应格式
- 把某些业务错误对象转换成特定视图或 JSON 结构
- 在 MVC 和 DI handler 上复用统一响应策略

初学者不用一开始就上这两个高级能力，但要知道 Iris 的 DI 不只是“参数自动传进去”，它连“结果怎么出去”也一起管了。

### 14.11 Iris DI 的特性总结

结合官方文档和示例，Iris DI 的几个核心特点可以总结为：

- 以函数签名为中心：参数和返回值本身就是接口定义的一部分
- 支持内建对象注入：如 `iris.Context`、`*http.Request` 等
- 支持请求数据绑定：路由参数、请求体、表单、查询参数等
- 支持自定义依赖：静态值和函数工厂都可以注册
- 支持请求期动态依赖：可以由中间件在单次请求中注入
- 支持结果处理：返回值可以直接成为 HTTP 响应
- 支持响应预处理：通过 `Preflight` 在响应前做最后控制
- 支持统一错误处理：通过 `OnError` 集中管理 DI 或 handler 的错误

### 14.12 Iris 的 DI 为什么很适合服务端初学者

它对初学者友好的点，不在于“高级”，而在于“把重复样板收起来了”。

你学习普通 Go Web 开发时，经常会反复写下面这些动作：

- 读路径参数
- 读 JSON
- 校验和转换
- 调业务服务
- 返回 JSON

Iris DI 让你可以更快地把注意力放到“接口的输入和输出”以及“业务逻辑本身”上。

它特别适合帮助新手建立这几个习惯：

- 把 HTTP 层和业务层分开
- 让 handler 保持短小
- 用结构体表达输入输出模型
- 用服务对象承载业务逻辑

### 14.13 但也不要把 DI 神化

DI 不是“高级项目才配用”，也不是“用了就一定更优雅”。它更像一个提高组织性的工具。

新手常见误区有三个：

#### 1. 把所有东西都做成依赖

如果只是一个很简单的小接口，直接用 `iris.Context` 写清楚也完全没问题。不要为了 DI 而 DI。

#### 2. 看不清参数来源

当函数参数越来越多时，要明确每个参数到底来自：

- 路由参数
- 请求体
- 容器依赖
- 中间件动态依赖

否则调试时会比较痛苦。

#### 3. 业务逻辑和传输模型混在一起

虽然 Iris 可以直接把请求体结构体注入 handler，但这不代表所有结构体都应该一路传到底层。实际项目中，通常仍然要区分：

- HTTP 请求/响应 DTO
- 业务层参数对象
- 数据库模型

### 14.14 推荐学习顺序

如果你是 Iris 或服务端开发初学者，建议按下面顺序理解这一节：

1. 先掌握最基础的 `iris.Context` 写法
2. 再理解“路由参数 + 请求体结构体”自动绑定
3. 然后学习 `RegisterDependency` 注入 Service
4. 最后再接触 `Preflight`、`OnError`、`UseResultHandler` 这些进阶能力

这样你不会一上来就被“自动注入”搞得失去方向，也更容易理解 Iris 为什么这样设计。

### 14.15 一个更贴近真实项目的示例

下面这个例子把“路由参数 + 请求体 + Service 注入 + 自动返回 JSON”串在一起。

```go
type UpdateUserRequest struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

type UpdateUserResponse struct {
	ID      uint64 `json:"id"`
	Message string `json:"message"`
}

type UserService struct{}

func (s *UserService) Update(id uint64, req UpdateUserRequest) UpdateUserResponse {
	return UpdateUserResponse{
		ID:      id,
		Message: req.Firstname + " " + req.Lastname + " updated",
	}
}

func updateUser(svc *UserService, id uint64, req UpdateUserRequest) UpdateUserResponse {
	return svc.Update(id, req)
}

func main() {
	app := iris.New()

	app.Party("/users").ConfigureContainer(func(api *iris.APIContainer) {
		api.RegisterDependency(&UserService{})
		api.Put("/{id:uint64}", updateUser)
	})

	app.Listen(":8080")
}
```

这个例子非常能体现 Iris DI 的风格：

- handler 只负责组合参数和调用服务
- service 只关心业务逻辑
- 路由参数、请求体和服务对象都由框架自动准备
- 返回结构体直接输出给客户端

如果你已经能读懂这个例子，说明你已经掌握了 Iris DI 的核心思路。

---

## 15. MVC 开发模式

Iris 对 MVC 提供了一等支持，这在 Go Web 框架里比较少见。它不是一个和普通路由完全割裂的新系统，而是建立在 Iris 的路由、依赖注入和结果分发能力之上的一层高级封装。

如果你前面已经理解了第 14 节的依赖注入，那么这一节就更容易掌握，因为 Iris MVC 本质上是在“Controller 结构体”上复用了 DI 的思路。

### 15.1 先理解：什么是 MVC

MVC 是一种经典的软件架构模式，全称：

- Model：模型
- View：视图
- Controller：控制器

它的核心目标是把不同职责拆开，避免所有代码都堆在一个地方。

在 Web 后端场景里，可以先这样理解：

- Model：业务数据和业务规则相关对象，比如用户、订单、文章
- View：把数据展示给客户端的方式，比如 HTML 页面、JSON 结果
- Controller：接收请求、组织流程、调用业务逻辑、返回结果

虽然很多现代 API 项目已经不严格强调“传统 MVC 三件套”的边界，但它背后的思想仍然很重要：

- 分层
- 解耦
- 让请求处理流程更清晰

对初学者来说，MVC 最值得学习的不是名词，而是“控制器不要把所有事情都做完”。Controller 应该主要负责：

- 接收和路由请求
- 调用 Service 或 Model 层
- 组织输出结果

而不是把数据库读写、业务规则、响应渲染全部揉成一个大函数。

### 15.2 MVC 在 Web 后端里通常长什么样

以“查询用户详情”为例，一个典型流程可以是：

1. 路由把 `GET /users/42` 分发给 Controller
2. Controller 从路径里拿到 `id`
3. Controller 调用 `UserService` 查询业务数据
4. Service 再调用 Repository 访问数据库
5. Controller 把结果作为 JSON 或 HTML 返回

如果是更传统的页面应用：

- Model 负责准备数据
- View 负责渲染模板
- Controller 负责把它们衔接起来

如果是今天常见的 REST API：

- View 往往不再是模板页面
- 而是 JSON 输出结构

所以在 Go Web 开发里，“MVC”很多时候可以理解成：

- 用 Controller 组织路由和请求处理
- 用 Model/Service/Repository 承担业务和数据职责
- 用 JSON/View 作为输出层

### 15.3 Iris 的 MVC 和传统 MVC 有什么关系

Iris 的 MVC 保留了 MVC 的核心思想，但做了更符合 Go 服务端开发习惯的实现。

它不是强行要求你必须有一个“View 文件夹 + Model 文件夹 + Controller 文件夹”的死板结构，而是提供一套 Controller 驱动的开发方式，让你可以：

- 用结构体表示 Controller
- 用方法名自动生成路由
- 用 DI 注入 `Context`、Service、Session 等依赖
- 用方法返回值自动生成响应
- 用生命周期方法在请求前后做统一处理

因此，Iris MVC 更准确的理解应该是：

它是“基于 Controller 的约定式路由 + 依赖注入 + 结果分发”的开发模式。

### 15.4 Iris MVC 的实现思路

从实现机制上看，Iris MVC 主要做了这几件事：

1. 启动时扫描 Controller 的导出方法
2. 根据方法命名规则推导 HTTP 方法和路由路径
3. 把这些方法注册到底层 Iris Router 上
4. 请求到来时创建或准备 Controller 实例
5. 通过 DI 给字段和方法参数注入依赖
6. 执行 Controller 方法
7. 把返回值自动转换为 HTTP 响应

所以 Iris MVC 不是“绕过路由系统”，而是“帮你更自动地生成路由并调用处理逻辑”。

你可以把它理解成：

- 路由系统负责 URL 分发
- DI 系统负责参数和依赖准备
- MVC 系统负责把 Controller 方法接到前两者上

### 15.5 基本用法

```go
import "github.com/kataras/iris/v12/mvc"

type UserController struct{}

func (c *UserController) Get() string {
	return "user list"
}

func (c *UserController) GetBy(id int) string {
	return "user id: " + strconv.Itoa(id)
}

func main() {
	app := iris.New()

	m := mvc.New(app.Party("/users"))
	m.Handle(new(UserController))

	app.Listen(":8080")
}
```

这段代码背后的意思是：

- `mvc.New(app.Party("/users"))` 创建一个挂载在 `/users` 路由组上的 MVC 应用
- `m.Handle(new(UserController))` 把 `UserController` 注册为一个控制器
- Iris 会扫描这个控制器的方法，并自动注册路由

说明：

- `Get()` 对应 `GET /users`
- `GetBy(id int)` 对应 `GET /users/{param:int}`

### 15.6 方法命名规则：Iris 如何把方法变成路由

Iris MVC 很依赖约定。官方文档给出的规则核心是：

- 方法名前缀是 HTTP Method 时，会被识别成路由处理器
- 方法名中的大写分段，会转换成路径片段
- `By` 关键字表示后面开始是动态路径参数

例如下面这个 Controller：

```go
type UserController struct{}

func (c *UserController) Get() string                  { return "list" }
func (c *UserController) GetLogin() string             { return "login page" }
func (c *UserController) PostLogin() string            { return "submit login" }
func (c *UserController) GetProfileFollowers() string  { return "followers" }
func (c *UserController) GetBy(id int64) string        { return "detail" }
func (c *UserController) GetUserBy(username string) string { return "user" }
```

如果它挂在 `/users` 下，大致会映射成：

- `Get()` -> `GET /users`
- `GetLogin()` -> `GET /users/login`
- `PostLogin()` -> `POST /users/login`
- `GetProfileFollowers()` -> `GET /users/profile/followers`
- `GetBy(id int64)` -> `GET /users/{param:int64}`
- `GetUserBy(username string)` -> `GET /users/user/{param:string}`

这也是 Iris MVC 最核心的“自动化”来源之一：你不是手工为每个方法写一条 `app.Get(...)`，而是让框架按照命名约定为你生成路由。

### 15.7 路径参数是怎么绑定的

当 Controller 方法带参数时，Iris 会结合 `By` 关键字和参数类型，自动把 URL 路径中的动态部分绑定到方法参数。

例如：

```go
func (c *UserController) GetBy(id uint64) iris.Map {
	return iris.Map{"id": id}
}
```

如果访问：

```text
GET /users/42
```

那么 `id` 就会被自动解析成 `42`。

支持的常见类型包括官方文档提到的：

- `int`
- `int64`
- `bool`
- `string`

此外也有像 `Wildcard` 这样的约定写法，用于匹配路径型参数。

这和第 14 节介绍的 DI 是同一思路：参数从哪里来，不需要你手动 `ctx.Params().Get(...)`，而是让框架按规则自动准备。

### 15.8 Controller 中可以注入什么

Iris MVC 并不只是“方法转路由”，它还能配合依赖注入使用。依赖可以注入到：

- Controller 的字段
- Controller 方法的参数

最常见的包括：

- `iris.Context`
- `*sessions.Session`
- 已注册的 Service
- 路由参数
- 其他可被 DI 解析的对象

例如：

```go
type UserService struct{}

func (s *UserService) List() []string {
	return []string{"tom", "jerry"}
}

type UserController struct {
	Ctx iris.Context
	Svc *UserService
}

func (c *UserController) Get() iris.Map {
	return iris.Map{
		"path":  c.Ctx.Path(),
		"users": c.Svc.List(),
	}
}
```

这表示：

- `Ctx` 由 Iris 注入当前请求上下文
- `Svc` 由容器注入已注册的业务服务

因此，Iris MVC 不是一套独立于 DI 的机制，而是基于 DI 的 Controller 模式。

### 15.9 返回值也是 MVC 的一部分

在 Iris MVC 中，Controller 方法的返回值不只是 Go 语言层面的返回值，它还是“响应怎么输出”的描述。

官方文档列出了很多可支持的返回形式，例如：

- `string`
- `int`
- `error`
- `(value, error)`
- 自定义结构体
- `mvc.Result`
- `(mvc.Result, error)`

例如：

```go
func (c *UserController) Get() string {
	return "hello"
}

func (c *UserController) GetProfile() iris.Map {
	return iris.Map{"name": "Tom"}
}

func (c *UserController) GetWelcome() mvc.Result {
	return mvc.Response{
		ContentType: "text/plain",
		Text:        "welcome",
	}
}
```

这背后说明 Iris MVC 会在方法执行后，再根据返回值类型做结果分发。

这和第 14 节 DI 里的“返回值自动写出响应”是一脉相承的。

### 15.10 生命周期方法

官方文档提到 Controller 可以定义：

- `BeforeActivation(b mvc.BeforeActivation)`
- `BeginRequest(ctx)`
- `EndRequest(ctx)`

它们分别在不同阶段工作。

#### 1. `BeforeActivation`

这是启动阶段的钩子，在 Controller 注册为路由之前执行。

它最重要的用途是：

- 给当前 Controller 添加自定义路由
- 给当前 Controller 的路由加中间件
- 调整该 Controller 作用域内的依赖

例如：

```go
func (c *UserController) BeforeActivation(b mvc.BeforeActivation) {
	b.Handle("GET", "/custom/{id:int64}", "Custom")
}

func (c *UserController) Custom(id int64) string {
	return "custom id: " + strconv.FormatInt(id, 10)
}
```

这意味着你并不一定只能依赖命名规则，也可以在 MVC 内部手动注册更自由的路由。

#### 2. `BeginRequest`

在每次请求真正进入目标 Controller 方法之前执行。

适合做：

- 初始化公共字段
- 预先准备通用数据
- 控制器级别的通用前置逻辑

#### 3. `EndRequest`

在请求结束后执行。

适合做：

- 清理逻辑
- 记录控制器级日志
- 收尾处理

### 15.11 错误处理

Iris MVC 对错误处理也有专门支持。官方文档提到：

- Controller 可以定义 `HandleHTTPError`
- 方法可以返回 `error`
- 方法也可以返回自定义 `mvc.Result`

这几种机制适合处理不同层面的错误。

#### 1. HTTP 错误

例如 404、500 这类 HTTP 错误，可以通过 `HandleHTTPError` 在 Controller 级别统一处理。

#### 2. 方法执行错误

如果 Controller 方法返回 `error` 或 `(value, error)`，框架会进入相应错误处理流程。

#### 3. 自定义结果对象

如果你希望把错误渲染成统一 JSON 或 HTML，也可以用自定义 `mvc.Result` 来控制输出。

这让 MVC 项目在“统一错误响应格式”方面更容易组织。

### 15.12 Controller 能做什么

结合官方文档，Iris MVC 支持：

- 控制器结构体
- 命名约定生成路由
- 依赖注入
- 生命周期钩子
- 视图渲染
- Session 注入
- 自定义结果输出
- 控制器级错误处理
- WebSocket 控制器

这也是它为什么经常被认为是 Iris 比较“全家桶”的能力之一。

### 15.13 Iris MVC 的优点

对中大型项目来说，Iris MVC 的优点主要在于：

- 路由和处理逻辑围绕 Controller 聚合，结构更集中
- 配合 DI 后，Controller 会更像业务入口层
- 同一组接口更容易统一挂中间件、依赖和错误处理
- 对模板站点、后台系统、多模块应用更友好

尤其当一个模块下有很多相关接口时，把它们聚合到一个 Controller 中，通常比散落成很多独立 handler 更容易维护。

### 15.14 Iris MVC 的代价和边界

但 MVC 也不是“越早用越好”。它的代价主要在于：

- 命名约定需要学习成本
- 自动路由会带来一点理解门槛
- 对非常简单的 API 来说，可能比普通 handler 更重

所以它更适合：

- 路由较多、模块边界清晰的项目
- 需要视图渲染的后台或站点
- 希望统一组织控制器、依赖和生命周期逻辑的团队

对于只有几个接口的小项目，普通路由 + handler 往往已经足够。

### 15.15 新手是否应该直接上 MVC

建议：

- 先学普通路由 + handler
- 再学第 14 节的依赖注入
- 最后再上 MVC

原因很简单：如果你一开始连 `iris.Context`、路由参数、请求体绑定、返回响应这些基础动作都还不熟，就直接学 MVC，容易只记住“方法名这样写能跑”，但不理解它背后的工作机制。

更好的顺序是：

1. 先会用普通 handler 写接口
2. 再理解 DI 如何减少样板代码
3. 最后理解 MVC 是如何把路由、DI、返回值和生命周期组合到一起的

这样你对 Iris MVC 的理解会更扎实。

### 15.16 一个更完整的 MVC 示例

```go
import (
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

type UserService struct{}

func (s *UserService) List() []string {
	return []string{"tom", "jerry"}
}

func (s *UserService) GetByID(id int64) iris.Map {
	return iris.Map{"id": id, "name": "tom"}
}

type UserController struct {
	Ctx iris.Context
	Svc *UserService
}

func (c *UserController) BeginRequest(ctx iris.Context) {
	ctx.Header("X-Controller", "UserController")
}

func (c *UserController) Get() iris.Map {
	return iris.Map{"users": c.Svc.List()}
}

func (c *UserController) GetBy(id int64) iris.Map {
	return c.Svc.GetByID(id)
}

func (c *UserController) BeforeActivation(b mvc.BeforeActivation) {
	b.Handle("GET", "/custom/{id:int64}", "Custom")
}

func (c *UserController) Custom(id int64) string {
	return "custom user: " + strconv.FormatInt(id, 10)
}

func main() {
	app := iris.New()

	m := mvc.New(app.Party("/users"))
	m.Register(&UserService{})
	m.Handle(new(UserController))

	app.Listen(":8080")
}
```

这个例子串起了 Iris MVC 的几个关键点：

- `mvc.New(app.Party("/users"))` 让 Controller 挂在 `/users`
- `m.Register(&UserService{})` 注册业务依赖
- `Get()` 和 `GetBy(id int64)` 按命名约定自动映射路由
- `BeforeActivation` 添加自定义路由
- `BeginRequest` 做控制器级前置处理
- 返回值直接成为 HTTP 响应

如果你能理解这个例子，说明你已经掌握了 Iris MVC 最核心的设计思路。

---

## 16. WebSocket

Iris 提供 WebSocket 支持，并且能和 MVC 结合使用。

官方文档中的思路是：把一个 Go 结构体注册成 WebSocket Controller，结构体中的导出方法就是事件处理器。

典型适用场景：

- 在线聊天室
- 实时通知
- 协同编辑
- 后台监控面板实时刷新

### 16.1 认识方式

WebSocket 的思维和普通 HTTP 不一样：

- HTTP 是请求一次、响应一次
- WebSocket 是连接建立后持续通信

### 16.2 新手建议

如果你刚接触 Web 开发：

- 先掌握普通路由和 JSON API
- 再理解 WebSocket 连接生命周期
- 最后再接入 Iris 的 WebSocket Controller

WebSocket 在 Iris 中是高级能力，适合第二阶段学习。

---

## 17. API 版本控制

如果你要长期维护 API，版本控制会非常重要。Iris 提供了专门的 `versioning` 子包。

### 17.1 基本示例

```go
import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/versioning"
)

func main() {
	app := iris.New()

	api := app.Party("/api")
	api.UseRouter(versioning.Aliases(versioning.AliasMap{
		versioning.Empty: "1.0.0",
		"latest":         "2.0.0",
	}))

	v1 := versioning.NewGroup(api, ">=1.0.0 <2.0.0")
	v1.Get("/users", func(ctx iris.Context) {
		ctx.JSON(iris.Map{"version": "v1"})
	})

	v2 := versioning.NewGroup(api, ">=2.0.0 <3.0.0")
	v2.Get("/users", func(ctx iris.Context) {
		ctx.JSON(iris.Map{"version": "v2"})
	})

	app.Listen(":8080")
}
```

### 17.2 版本来源

官方文档说明版本通常从以下请求头获取：

- `Accept`
- `Accept-Version`

也支持你自己从 Query 中读取后再设置。

### 17.3 为什么有用

适合这些情况：

- 老客户端不能立刻升级
- 新旧接口要同时维护
- 需要接口弃用通知

---

## 18. 常见安全能力

Iris 官方文档中单独列出了多类安全主题，例如：

- Basic Authentication
- CORS
- Sessions & Cookies
- CSRF
- JWT
- CAPTCHA

对于新手，最先要掌握的是这 4 件事：

### 18.1 基本认证思路

- 登录态与鉴权逻辑放到中间件
- 不要在每个 handler 里重复判断
- API 返回统一错误格式

### 18.2 CORS

前后端分离项目经常需要跨域。跨域配置一般应集中处理，不要散落在单个接口中。

### 18.3 Session/Cookie 安全

- 生产环境启用 HTTPS
- 配置合适的 `SameSite`
- 控制 Cookie 域与作用范围

### 18.4 JWT

如果你使用 JWT，建议把：

- 解析
- 验签
- 续期
- 黑名单或失效策略

都设计清楚，不要只做最基础的“生成一个 token”。

---

## 19. 测试

Iris 具备较好的测试支持，官方示例里常使用 `httptest` 子包。

### 19.1 基本思路

```go
import "github.com/kataras/iris/v12/httptest"
```

测试时先创建应用，再发起请求断言结果。

### 19.2 典型写法

```go
func TestIndex(t *testing.T) {
	app := iris.New()
	app.Get("/", func(ctx iris.Context) {
		ctx.JSON(iris.Map{"message": "ok"})
	})

	e := httptest.New(t, app)
	e.GET("/").Expect().Status(httptest.StatusOK).
		JSON().Object().ValueEqual("message", "ok")
}
```

### 19.3 新手测试建议

至少覆盖：

- 路由是否注册成功
- 参数绑定是否正确
- 校验失败是否返回 400
- 鉴权失败是否返回 401/403
- 成功响应格式是否稳定

---

## 20. 一个适合新手的 Iris 项目结构

下面是一个比较容易维护的目录建议：

```text
myapp/
├── main.go
├── go.mod
├── configs/
│   └── config.go
├── routes/
│   └── router.go
├── controllers/
│   └── user_controller.go
├── services/
│   └── user_service.go
├── repositories/
│   └── user_repository.go
├── models/
│   └── user.go
├── middlewares/
│   └── auth.go
├── views/
│   ├── layouts/
│   └── index.html
├── assets/
│   ├── css/
│   └── js/
└── tests/
    └── user_test.go
```

建议职责：

- `main.go`：程序入口
- `routes/`：统一注册路由
- `controllers/`：接收请求、组织响应
- `services/`：业务逻辑
- `repositories/`：数据库访问
- `models/`：数据结构
- `middlewares/`：认证、日志、跨域等
- `views/`：模板文件
- `assets/`：静态资源

---

## 21. 从零写一个最小 REST API 示例

```go
package main

import (
	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12"
)

type CreateUserRequest struct {
	Name  string `json:"name" validate:"required,min=2"`
	Email string `json:"email" validate:"required,email"`
}

func main() {
	app := iris.New()
	validate := validator.New()

	app.Use(iris.Compression)

	api := app.Party("/api")

	api.Get("/ping", func(ctx iris.Context) {
		ctx.JSON(iris.Map{"message": "pong"})
	})

	api.Post("/users", func(ctx iris.Context) {
		var req CreateUserRequest
		if err := ctx.ReadJSON(&req); err != nil {
			ctx.StopWithError(iris.StatusBadRequest, err)
			return
		}

		if err := validate.Struct(req); err != nil {
			ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().
				Title("Validation Error").
				DetailErr(err))
			return
		}

		ctx.StatusCode(iris.StatusCreated)
		ctx.JSON(iris.Map{
			"code":    0,
			"message": "user created",
			"data": iris.Map{
				"name":  req.Name,
				"email": req.Email,
			},
		})
	})

	app.Listen(":8080")
}
```

这个例子已经具备了最常见的开发要素：

- 路由
- JSON 读取
- 参数校验
- JSON 响应
- 状态码
- 压缩中间件

---

## 22. 新手最容易踩的坑

### 22.1 把业务逻辑全写进 handler

问题：

- 难测试
- 难复用
- 难维护

建议：handler 只负责“接收请求 + 调业务 + 返回响应”。

### 22.2 路由层和数据层耦合过深

建议中间加 `service` 层，不要让路由直接操作数据库。

### 22.3 没有统一错误处理

建议统一：

- 参数错误
- 权限错误
- 系统错误

返回格式和日志方式。

### 22.4 响应格式不统一

同一个项目里，尽量不要有的接口返回字符串，有的返回裸数组，有的返回复杂对象。

### 22.5 一开始就上过重架构

建议学习顺序：

1. 基本路由
2. Context 读写
3. 中间件
4. JSON 请求与响应
5. 模板和静态资源
6. Session
7. DI
8. MVC
9. WebSocket
10. 版本控制与高级安全能力

---

## 23. 学习路线建议

### 第一阶段：能跑起来

先学：

- 安装
- `iris.New`
- 基本路由
- `Context`
- JSON 请求与响应

目标：写出一个最小 API 服务。

### 第二阶段：能组织项目

再学：

- Party 路由分组
- 中间件
- 模板渲染
- 静态文件
- Session

目标：写出一个小型完整 Web 项目。

### 第三阶段：能做中大型项目

继续学：

- 依赖注入
- MVC
- API 版本控制
- WebSocket
- 测试
- 安全能力

目标：搭建可维护、可扩展的业务系统。

---

## 24. 快速对照表

### 24.1 常用应用方法

| 方法 | 作用 |
| --- | --- |
| `iris.New()` | 创建应用 |
| `app.Listen(":8080")` | 启动服务 |
| `app.Get/Post/Put/Delete` | 注册路由 |
| `app.Party("/api")` | 创建路由分组 |
| `app.Use(middleware)` | 注册中间件 |
| `app.RegisterView(...)` | 注册模板引擎 |
| `app.HandleDir(...)` | 托管静态目录 |

### 24.2 常用 Context 方法

| 方法 | 作用 |
| --- | --- |
| `ctx.Params().Get("id")` | 获取路由参数 |
| `ctx.URLParam("name")` | 获取查询参数 |
| `ctx.ReadJSON(&v)` | 读取 JSON |
| `ctx.JSON(v)` | 返回 JSON |
| `ctx.XML(v)` | 返回 XML |
| `ctx.Text("...")` | 返回文本 |
| `ctx.View("a.html", data)` | 渲染模板 |
| `ctx.StatusCode(code)` | 设置状态码 |
| `ctx.Header(key, value)` | 设置响应头 |
| `ctx.Next()` | 执行下一个处理器 |
| `ctx.StopWithStatus(code)` | 中断并返回状态码 |

---

## 25. 总结

Iris 的特点可以概括为：

- 不只是一个轻量路由器，而是一个完整 Web 框架
- API、模板站点、Session、MVC、WebSocket 都能覆盖
- 依赖注入和 MVC 是它比较有辨识度的优势
- 对新手来说，上手简单，但深入后功能也足够丰富

如果你刚开始学 Iris，最推荐的顺序是：

先把基础路由、请求读取、JSON 响应、中间件掌握好，再逐步学习模板、Session、DI、MVC 和 WebSocket。这样学习成本最低，也最符合真实项目演进过程。

---

## 26. 官方参考资料

以下资料是整理本手册时重点参考的 Iris 官方文档与官方仓库：

- 官方文档首页：https://iris-go.com/docs/
- 安装与快速开始：https://iris-go.com/docs/getting-started/installation.html
- 路由中间件：https://iris-go.com/docs/routing/middleware.html
- 模板视图：https://iris-go.com/docs/view/view.html
- 静态文件服务：https://iris-go.com/docs/file-server/introduction.html
- 依赖注入：https://iris-go.com/docs/dependency-injection/dependency-injection.html
- DI 输入类型：https://iris-go.com/docs/dependency-injection/inputs.html
- MVC：https://iris-go.com/docs/mvc/mvc.html
- MVC WebSocket：https://iris-go.com/docs/mvc/mvc-websockets.html
- MVC Session：https://iris-go.com/docs/mvc/mvc-sessions.html
- API 版本控制：https://iris-go.com/docs/routing/api-versioning.html
- JSON 请求处理：https://iris-go.com/docs/requests/json.html
- 请求校验：https://iris-go.com/docs/requests/validation.html
- JSON 响应：https://iris-go.com/docs/responses/json.html
- 压缩支持：https://iris-go.com/docs/compression/compression.html
- Session 与 Cookie 安全：https://iris-go.com/docs/security/security-sessions-cookies.html
- 官方 GitHub 仓库：https://github.com/kataras/iris

---

## 27. 使用说明

这份文档是面向入门的工具书，重点放在：

- 解释概念
- 给出最常用写法
- 帮你建立 Iris 的整体认知

如果你后续要继续深入，建议下一步结合官方 `_examples` 和真实项目一起练习。
