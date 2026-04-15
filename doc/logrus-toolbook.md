# logrus 工具箱

> 基于 logrus 官方文档和实践经验整理，面向需要日志管理的 Go 开发者。
> 
> 目标：读完后能独立使用 logrus 进行日志记录，掌握常见配置和最佳实践。

---

## 1. logrus 是什么

logrus 是 Go 语言中最流行的日志库之一，提供了结构化日志记录能力，支持多种输出格式和字段类型。它是标准库 `log` 包的增强版，提供了更丰富的功能和更灵活的配置选项。

logrus 的核心能力包括：

- 结构化日志输出（JSON 格式）
- 多种日志级别（Debug、Info、Warn、Error、Fatal、Panic）
- 字段（Fields）支持，便于日志分析和过滤
- 钩子（Hooks）机制，可扩展日志处理行为
- 自定义格式化器
- 线程安全

适合使用 logrus 的场景：

- 生产环境的应用日志
- 需要结构化分析的日志
- 微服务架构中的日志管理
- 任何需要灵活日志配置的 Go 项目

---

## 2. 安装与环境准备

### 2.1 安装 logrus

使用 go get 命令安装 logrus：

```bash
go get github.com/sirupsen/logrus
```

### 2.2 导入 logrus

在 Go 代码中导入 logrus：

```go
import "github.com/sirupsen/logrus"
```

---

## 3. 基本使用

### 3.1 简单日志输出

```go
package main

import "github.com/sirupsen/logrus"

func main() {
    logrus.Info("This is an info message")
    logrus.Warn("This is a warning message")
    logrus.Error("This is an error message")
}
```

### 3.2 日志级别

logrus 支持以下日志级别（按严重程度递增）：

- Debug
- Info
- Warn
- Error
- Fatal
- Panic

默认日志级别是 Info，可以通过以下方式设置：

```go
logrus.SetLevel(logrus.DebugLevel)
```

### 3.3 字段（Fields）

使用字段可以为日志添加结构化信息：

```go
logrus.WithFields(logrus.Fields{
    "user_id": 123,
    "action": "login",
}).Info("User logged in")
```

输出：
```
time="2023-10-01T10:00:00Z" level=info msg="User logged in" user_id=123 action=login
```

---

## 4. 高级特性

### 4.1 JSON 格式输出

```go
logrus.SetFormatter(&logrus.JSONFormatter{})
```

输出：
```json
{"level":"info","msg":"User logged in","time":"2023-10-01T10:00:00Z","user_id":123,"action":"login"}
```

### 4.2 自定义日志格式

```go
logrus.SetFormatter(&logrus.TextFormatter{
    FullTimestamp:   true,
    TimestampFormat: "2006-01-02 15:04:05",
    ForceColors:     true,
})
```

### 4.3 输出到文件

```go
file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
if err != nil {
    logrus.Fatal("Failed to open log file", err)
}
logrus.SetOutput(file)
```

### 4.4 同时输出到控制台和文件

```go
file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
if err != nil {
    logrus.Fatal("Failed to open log file", err)
}

// 同时输出到控制台和文件
logrus.SetOutput(io.MultiWriter(os.Stdout, file))
```

#### 配置执行时机

**重要**：输出目标的配置只需要执行一次，而不是每次输出日志时都执行。建议在应用启动时完成配置：

1. **在 `init()` 函数中配置**：
   ```go
   func init() {
       // 设置日志级别
       logrus.SetLevel(logrus.InfoLevel)
       
       // 设置输出格式
       logrus.SetFormatter(&logrus.JSONFormatter{})
       
       // 同时输出到控制台和文件
       file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
       if err != nil {
           logrus.Fatal("Failed to open log file", err)
       }
       logrus.SetOutput(io.MultiWriter(os.Stdout, file))
   }
   ```

2. **使用日志包装器**：
   创建一个专门的日志包，在初始化时配置一次，然后在整个应用中使用（详见 6.2 节）。

这样配置后，所有的日志输出都会自动同时发送到控制台和文件，无需每次输出日志时重复配置。

---

## 5. 钩子（Hooks）

钩子允许在日志记录时执行额外的操作，例如发送邮件、记录到数据库等。

### 5.1 飞书群通知钩子

下面是一个完整的飞书群通知 webhook 钩子实现示例：

```go
import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "time"
    "github.com/sirupsen/logrus"
)

// FeishuHook 飞书群通知钩子
type FeishuHook struct {
    webhookURL string
    title      string
}

// FeishuMessage 飞书消息格式
type FeishuMessage struct {
    MsgType string `json:"msg_type"`
    Content Content `json:"content"`
}

// Content 消息内容
type Content struct {
    Post Post `json:"post"`
}

// Post 飞书富文本消息
type Post struct {
    ZhCN ZhCN `json:"zh_cn"`
}

// ZhCN 中文内容
type ZhCN struct {
    Title   string     `json:"title"`
    Content [][]string `json:"content"`
}

// NewFeishuHook 创建飞书钩子
func NewFeishuHook(webhookURL, title string) *FeishuHook {
    return &FeishuHook{
        webhookURL: webhookURL,
        title:      title,
    }
}

// Fire 触发钩子
func (h *FeishuHook) Fire(entry *logrus.Entry) error {
    // 构建飞书消息
    content := [][]string{
        {"**日志级别**", entry.Level.String()},
        {"**消息内容**", entry.Message},
        {"**时间**", time.Now().Format("2006-01-02 15:04:05")},
    }

    // 添加额外字段
    for key, value := range entry.Data {
        content = append(content, []string{"**" + key + "**", fmt.Sprintf("%v", value)})
    }

    message := FeishuMessage{
        MsgType: "post",
        Content: Content{
            Post: Post{
                ZhCN: ZhCN{
                    Title:   h.title,
                    Content: content,
                },
            },
        },
    }

    // 序列化消息
    payload, err := json.Marshal(message)
    if err != nil {
        return err
    }

    // 发送 HTTP 请求
    client := &http.Client{Timeout: 5 * time.Second}
    _, err = client.Post(h.webhookURL, "application/json", bytes.NewBuffer(payload))
    return err
}

// Levels 返回需要触发的日志级别
func (h *FeishuHook) Levels() []logrus.Level {
    return []logrus.Level{logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel}
}

// 添加钩子
feishuHook := NewFeishuHook(
    "https://open.feishu.cn/open-apis/bot/v2/hook/your-webhook-url",
    "应用日志通知",
)
logrus.AddHook(feishuHook)
```

### 5.2 飞书群通知配置最佳实践

1. **错误处理**：确保 webhook 发送失败不会影响应用运行
2. **超时设置**：为 HTTP 请求设置合理的超时时间（建议 5 秒）
3. **异步处理**：使用 goroutine 异步发送通知，避免阻塞主流程
4. **安全考虑**：webhook URL 应存储在环境变量或配置文件中，不要硬编码
5. **消息格式**：保持消息格式清晰，包含关键信息
6. **重试机制**：可添加简单的重试逻辑，提高可靠性

#### 异步飞书通知示例

```go
func (h *FeishuHook) Fire(entry *logrus.Entry) error {
    // 异步发送
    go func() {
        // 构建和发送消息的逻辑
        // ...
    }()
    return nil
}
```

### 5.3 飞书群机器人设置

1. **创建飞书群机器人**：
   - 打开飞书群聊，点击群设置
   - 选择「群机器人」→「添加机器人」→「自定义」
   - 填写机器人名称，获取 webhook URL

2. **权限设置**：
   - 确保机器人有发送消息的权限
   - 如需发送卡片消息，需开启相应权限

3. **测试通知**：
   - 部署钩子后，可通过触发错误日志测试通知是否正常

---

## 6. 实际应用示例

### 6.1 基本应用结构

```go
package main

import (
    "github.com/sirupsen/logrus"
    "os"
    "io"
)

func init() {
    // 设置日志级别
    logrus.SetLevel(logrus.InfoLevel)
    
    // 设置输出格式
    logrus.SetFormatter(&logrus.JSONFormatter{})
    
    // 同时输出到控制台和文件
    file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        logrus.Fatal("Failed to open log file", err)
    }
    logrus.SetOutput(io.MultiWriter(os.Stdout, file))
}

func main() {
    logrus.Info("Application started")
    
    // 带字段的日志
    logrus.WithFields(logrus.Fields{
        "module": "user",
        "operation": "create",
    }).Info("User creation initiated")
    
    // 错误日志
    if err := someFunction(); err != nil {
        logrus.WithError(err).Error("Failed to process request")
    }
    
    logrus.Info("Application finished")
}

func someFunction() error {
    return fmt.Errorf("simulated error")
}
```

### 6.2 日志包装器

为了更好地组织日志代码，可以创建一个日志包装器：

```go
package logger

import (
    "github.com/sirupsen/logrus"
    "os"
    "io"
)

var Log *logrus.Logger

func Init() {
    Log = logrus.New()
    Log.SetLevel(logrus.InfoLevel)
    
    // 设置输出格式
    Log.SetFormatter(&logrus.JSONFormatter{})
    
    // 同时输出到控制台和文件
    file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        Log.Fatal("Failed to open log file", err)
    }
    Log.SetOutput(io.MultiWriter(os.Stdout, file))
}

func Info(message string, fields ...logrus.Fields) {
    if len(fields) > 0 {
        Log.WithFields(fields[0]).Info(message)
    } else {
        Log.Info(message)
    }
}

func Error(message string, err error, fields ...logrus.Fields) {
    entry := Log.WithError(err)
    if len(fields) > 0 {
        entry = entry.WithFields(fields[0])
    }
    entry.Error(message)
}
```

使用方式：

```go
package main

import "your/project/logger"

func main() {
    logger.Init()
    logger.Info("Application started")
    
    if err := someFunction(); err != nil {
        logger.Error("Failed to process request", err, logrus.Fields{
            "context": "main",
        })
    }
}
```

---

## 7. 最佳实践

### 7.1 日志级别使用建议

- **Debug**：开发和调试时使用，包含详细的内部状态信息
- **Info**：正常运行时的信息，如服务启动、请求处理开始等
- **Warn**：需要注意但不影响正常运行的情况，如配置项缺失
- **Error**：错误情况，但应用可以继续运行，如请求处理失败
- **Fatal**：严重错误，应用无法继续运行，会调用 os.Exit(1)
- **Panic**：极其严重的错误，会触发 panic

### 7.2 字段使用建议

- 使用一致的字段命名规范
- 为关键操作添加上下文信息
- 避免在字段中包含敏感信息（如密码、token）
- 合理使用字段分组，便于日志分析

### 7.3 性能考虑

- 对于高频日志，考虑使用 `logrus.WithFields` 预创建日志条目
- 避免在日志消息中进行复杂的字符串拼接
- 对于生产环境，合理设置日志级别，避免过多的 Debug 日志

### 7.4 日志文件管理

- 定期轮转日志文件，避免单个文件过大
- 设置合理的日志保留策略
- 考虑使用日志聚合服务（如 ELK Stack、Graylog 等）

---

## 8. 常见问题

### 8.1 日志输出格式问题

**问题**：日志输出格式不符合预期。

**解决方案**：检查格式化器配置，确保使用了正确的格式化器。

```go
// JSON 格式
logrus.SetFormatter(&logrus.JSONFormatter{})

// 文本格式
logrus.SetFormatter(&logrus.TextFormatter{
    FullTimestamp: true,
})
```

### 8.2 日志级别不生效

**问题**：设置了日志级别但所有级别的日志都输出了。

**解决方案**：确保在输出日志之前设置日志级别。

```go
// 正确顺序
logrus.SetLevel(logrus.WarnLevel)
logrus.Debug("This should not be logged")

// 错误顺序
logrus.Debug("This will be logged")
logrus.SetLevel(logrus.WarnLevel)
```

### 8.3 钩子不工作

**问题**：添加了钩子但没有执行。

**解决方案**：检查钩子的 `Levels()` 方法是否包含了要触发的日志级别。

```go
func (h *MyHook) Levels() []logrus.Level {
    return []logrus.Level{logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel}
}
```

---

## 9. 总结

logrus 是一个功能强大、灵活易用的 Go 日志库，通过本文的介绍，相信您已经掌握了它的基本使用方法和高级特性。在实际项目中，合理使用 logrus 可以帮助您更好地监控应用运行状态、排查问题，提高系统的可维护性。

记住以下几点：

- 根据实际需求选择合适的日志级别
- 合理使用字段添加上下文信息
- 配置适当的输出格式和目标
- 利用钩子扩展日志功能
- 遵循最佳实践，保持日志的一致性和可读性

通过这些实践，您可以充分发挥 logrus 的优势，为您的应用提供可靠的日志记录能力。