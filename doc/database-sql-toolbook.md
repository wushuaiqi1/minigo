# Go database/sql 工具箱

> 基于 Go 官方 database/sql 包文档和实践经验整理，面向需要在 Go 中操作数据库的开发者。
> 
> 目标：读完后能独立使用 database/sql 包进行数据库操作，掌握常见配置和最佳实践。

---

## 1. database/sql 是什么

database/sql 是 Go 语言标准库中提供的一个通用 SQL 数据库接口包。它提供了一个抽象层，使得开发者可以使用统一的 API 操作不同的数据库（如 MySQL、PostgreSQL、SQLite 等），而无需关心底层数据库驱动的具体实现。

database/sql 的核心能力包括：

- 统一的数据库操作接口
- 连接池管理
- 事务支持
- 预编译语句
- 上下文支持（Context）
- 空值处理
- 数据类型转换

适合使用 database/sql 的场景：

- 需要操作关系型数据库的 Go 应用
- 需要支持多种数据库的应用
- 需要高性能数据库连接管理的应用
- 任何需要事务支持的数据库操作

---

## 2. 安装与环境准备

### 2.1 安装数据库驱动

database/sql 包需要与具体的数据库驱动配合使用。常用的数据库驱动：

```bash
# MySQL 驱动
go get github.com/go-sql-driver/mysql

# PostgreSQL 驱动
go get github.com/lib/pq

# SQLite 驱动
go get github.com/mattn/go-sqlite3

# SQL Server 驱动
go get github.com/denisenkom/go-mssqldb
```

### 2.2 导入包

```go
import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql" // MySQL 驱动
)
```

注意：使用空白标识符 `_` 导入驱动，只执行其 `init()` 函数注册驱动，不直接使用驱动包。

---

## 3. 基本使用

### 3.1 连接数据库

```go
import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
)

func main() {
    // 连接 MySQL
    db, err := sql.Open("mysql", "username:password@tcp(localhost:3306)/dbname")
    if err != nil {
        panic(err)
    }
    defer db.Close()
    
    // 测试连接
    err = db.Ping()
    if err != nil {
        panic(err)
    }
    
    fmt.Println("数据库连接成功")
}
```

### 3.2 数据源格式

不同数据库的数据源格式：

```go
// MySQL
"username:password@tcp(host:port)/dbname?parseTime=true"

// PostgreSQL
"host=localhost port=5432 user=username password=password dbname=dbname sslmode=disable"

// SQLite
"file:test.db?cache=shared&mode=memory"

// SQL Server
"server=localhost;user id=username;password=password;database=dbname"
```

### 3.3 执行查询

```go
// 查询单行
var name string
var age int
err := db.QueryRow("SELECT name, age FROM users WHERE id = ?", 1).Scan(&name, &age)
if err != nil {
    if err == sql.ErrNoRows {
        fmt.Println("没有找到记录")
    } else {
        panic(err)
    }
}
fmt.Printf("姓名: %s, 年龄: %d\n", name, age)

// 查询多行
rows, err := db.Query("SELECT id, name, age FROM users WHERE age > ?", 18)
if err != nil {
    panic(err)
}
defer rows.Close()

for rows.Next() {
    var id int
    var name string
    var age int
    err := rows.Scan(&id, &name, &age)
    if err != nil {
        panic(err)
    }
    fmt.Printf("ID: %d, 姓名: %s, 年龄: %d\n", id, name, age)
}

// 检查遍历过程中的错误
if err = rows.Err(); err != nil {
    panic(err)
}
```

### 3.4 执行插入、更新、删除

```go
// 插入数据
result, err := db.Exec("INSERT INTO users(name, age) VALUES(?, ?)", "张三", 25)
if err != nil {
    panic(err)
}

// 获取插入的 ID
lastID, err := result.LastInsertId()
if err != nil {
    panic(err)
}
fmt.Printf("插入的 ID: %d\n", lastID)

// 获取影响的行数
rowsAffected, err := result.RowsAffected()
if err != nil {
    panic(err)
}
fmt.Printf("影响的行数: %d\n", rowsAffected)

// 更新数据
result, err = db.Exec("UPDATE users SET age = ? WHERE id = ?", 26, 1)
if err != nil {
    panic(err)
}

// 删除数据
result, err = db.Exec("DELETE FROM users WHERE id = ?", 1)
if err != nil {
    panic(err)
}
```

---

## 4. 高级特性

### 4.1 事务处理

```go
// 开始事务
tx, err := db.Begin()
if err != nil {
    panic(err)
}

// 确保事务被回滚或提交
defer func() {
    if p := recover(); p != nil {
        tx.Rollback()
        panic(p) // 重新抛出 panic
    }
}()

// 执行多个操作
_, err = tx.Exec("UPDATE accounts SET balance = balance - ? WHERE id = ?", 100, 1)
if err != nil {
    tx.Rollback()
    panic(err)
}

_, err = tx.Exec("UPDATE accounts SET balance = balance + ? WHERE id = ?", 100, 2)
if err != nil {
    tx.Rollback()
    panic(err)
}

// 提交事务
err = tx.Commit()
if err != nil {
    panic(err)
}
fmt.Println("事务提交成功")
```

### 4.2 预编译语句

```go
// 准备预编译语句
stmt, err := db.Prepare("SELECT name, age FROM users WHERE id = ?")
if err != nil {
    panic(err)
}
defer stmt.Close()

// 多次执行预编译语句
for _, id := range []int{1, 2, 3} {
    var name string
    var age int
    err := stmt.QueryRow(id).Scan(&name, &age)
    if err != nil {
        if err == sql.ErrNoRows {
            continue
        }
        panic(err)
    }
    fmt.Printf("ID: %d, 姓名: %s, 年龄: %d\n", id, name, age)
}
```

### 4.3 上下文支持

```go
import (
    "context"
    "time"
)

// 带超时的查询
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

var name string
err := db.QueryRowContext(ctx, "SELECT name FROM users WHERE id = ?", 1).Scan(&name)
if err != nil {
    if err == context.DeadlineExceeded {
        fmt.Println("查询超时")
    } else {
        panic(err)
    }
}

// 带上下文的事务
ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

tx, err := db.BeginTx(ctx, nil)
if err != nil {
    panic(err)
}
defer tx.Rollback()

// 执行事务操作
_, err = tx.ExecContext(ctx, "UPDATE users SET age = age + 1 WHERE id = ?", 1)
if err != nil {
    panic(err)
}

err = tx.Commit()
if err != nil {
    panic(err)
}
```

### 4.4 连接池配置

```go
// 设置最大打开连接数
db.SetMaxOpenConns(25)

// 设置最大空闲连接数
db.SetMaxIdleConns(25)

// 设置连接最大存活时间
db.SetConnMaxLifetime(5 * time.Minute)

// 设置连接最大空闲时间
db.SetConnMaxIdleTime(2 * time.Minute)

// 获取连接池统计信息
stats := db.Stats()
fmt.Printf("打开的连接数: %d\n", stats.OpenConnections)
fmt.Printf("空闲连接数: %d\n", stats.Idle)
fmt.Printf("使用中的连接数: %d\n", stats.InUse)
```

---

## 5. 空值处理

### 5.1 使用 Null 类型

database/sql 包提供了多种 Null 类型来处理数据库中的 NULL 值：

```go
import (
    "database/sql"
    "fmt"
)

// NullString
var name sql.NullString
err := db.QueryRow("SELECT name FROM users WHERE id = ?", 1).Scan(&name)
if err != nil {
    panic(err)
}

if name.Valid {
    fmt.Printf("姓名: %s\n", name.String)
} else {
    fmt.Println("姓名为 NULL")
}

// NullInt64
var age sql.NullInt64
err = db.QueryRow("SELECT age FROM users WHERE id = ?", 1).Scan(&age)
if err != nil {
    panic(err)
}

if age.Valid {
    fmt.Printf("年龄: %d\n", age.Int64)
} else {
    fmt.Println("年龄为 NULL")
}

// NullFloat64
var salary sql.NullFloat64
err = db.QueryRow("SELECT salary FROM employees WHERE id = ?", 1).Scan(&salary)
if err != nil {
    panic(err)
}

if salary.Valid {
    fmt.Printf("薪资: %.2f\n", salary.Float64)
} else {
    fmt.Println("薪资为 NULL")
}

// NullBool
var active sql.NullBool
err = db.QueryRow("SELECT active FROM users WHERE id = ?", 1).Scan(&active)
if err != nil {
    panic(err)
}

if active.Valid {
    fmt.Printf("是否激活: %t\n", active.Bool)
} else {
    fmt.Println("激活状态为 NULL")
}

// NullTime
var createdAt sql.NullTime
err = db.QueryRow("SELECT created_at FROM users WHERE id = ?", 1).Scan(&createdAt)
if err != nil {
    panic(err)
}

if createdAt.Valid {
    fmt.Printf("创建时间: %s\n", createdAt.Time)
} else {
    fmt.Println("创建时间为 NULL")
}
```

### 5.2 插入 NULL 值

```go
// 插入 NULL 值
var name sql.NullString
name.String = "张三"
name.Valid = true // 非空

var age sql.NullInt64
age.Int64 = 0
age.Valid = false // NULL

_, err := db.Exec("INSERT INTO users(name, age) VALUES(?, ?)", name, age)
if err != nil {
    panic(err)
}

// 或者直接使用 nil
_, err = db.Exec("INSERT INTO users(name, age) VALUES(?, NULL)", "李四")
if err != nil {
    panic(err)
}
```

### 5.3 泛型 Null 类型（Go 1.18+）

```go
import "database/sql"

// 使用泛型 Null 类型
var name sql.Null[string]
err := db.QueryRow("SELECT name FROM users WHERE id = ?", 1).Scan(&name)
if err != nil {
    panic(err)
}

if name.Valid {
    fmt.Printf("姓名: %s\n", name.V)
} else {
    fmt.Println("姓名为 NULL")
}

// 插入 NULL 值
var age sql.Null[int64]
age.V = 25
age.Valid = true

_, err = db.Exec("INSERT INTO users(name, age) VALUES(?, ?)", "王五", age)
if err != nil {
    panic(err)
}
```

---

## 6. 实际应用示例

### 6.1 用户管理系统

```go
package main

import (
    "database/sql"
    "fmt"
    "log"
    _ "github.com/go-sql-driver/mysql"
)

type User struct {
    ID        int
    Name      string
    Age       int
    Email     string
    CreatedAt time.Time
}

type UserRepository struct {
    db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
    return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *User) error {
    result, err := r.db.Exec(
        "INSERT INTO users(name, age, email) VALUES(?, ?, ?)",
        user.Name, user.Age, user.Email,
    )
    if err != nil {
        return err
    }

    id, err := result.LastInsertId()
    if err != nil {
        return err
    }

    user.ID = int(id)
    return nil
}

func (r *UserRepository) GetByID(id int) (*User, error) {
    var user User
    err := r.db.QueryRow(
        "SELECT id, name, age, email, created_at FROM users WHERE id = ?",
        id,
    ).Scan(&user.ID, &user.Name, &user.Age, &user.Email, &user.CreatedAt)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil
        }
        return nil, err
    }
    return &user, nil
}

func (r *UserRepository) GetAll() ([]User, error) {
    rows, err := r.db.Query("SELECT id, name, age, email, created_at FROM users")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var users []User
    for rows.Next() {
        var user User
        err := rows.Scan(&user.ID, &user.Name, &user.Age, &user.Email, &user.CreatedAt)
        if err != nil {
            return nil, err
        }
        users = append(users, user)
    }

    if err = rows.Err(); err != nil {
        return nil, err
    }

    return users, nil
}

func (r *UserRepository) Update(user *User) error {
    _, err := r.db.Exec(
        "UPDATE users SET name = ?, age = ?, email = ? WHERE id = ?",
        user.Name, user.Age, user.Email, user.ID,
    )
    return err
}

func (r *UserRepository) Delete(id int) error {
    _, err := r.db.Exec("DELETE FROM users WHERE id = ?", id)
    return err
}

func main() {
    db, err := sql.Open("mysql", "username:password@tcp(localhost:3306)/testdb")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    if err = db.Ping(); err != nil {
        log.Fatal(err)
    }

    repo := NewUserRepository(db)

    // 创建用户
    user := &User{
        Name:  "张三",
        Age:   25,
        Email: "zhangsan@example.com",
    }
    if err = repo.Create(user); err != nil {
        log.Fatal(err)
    }
    fmt.Printf("创建用户成功，ID: %d\n", user.ID)

    // 查询用户
    user, err = repo.GetByID(user.ID)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("查询用户: %+v\n", user)

    // 更新用户
    user.Age = 26
    if err = repo.Update(user); err != nil {
        log.Fatal(err)
    }
    fmt.Println("更新用户成功")

    // 删除用户
    if err = repo.Delete(user.ID); err != nil {
        log.Fatal(err)
    }
    fmt.Println("删除用户成功")
}
```

### 6.2 事务处理示例

```go
package main

import (
    "database/sql"
    "fmt"
    "log"
    _ "github.com/go-sql-driver/mysql"
)

type Account struct {
    ID      int
    Name    string
    Balance float64
}

type TransactionService struct {
    db *sql.DB
}

func NewTransactionService(db *sql.DB) *TransactionService {
    return &TransactionService{db: db}
}

func (s *TransactionService) Transfer(fromID, toID int, amount float64) error {
    tx, err := s.db.Begin()
    if err != nil {
        return err
    }
    defer func() {
        if p := recover(); p != nil {
            tx.Rollback()
            panic(p)
        }
    }()

    // 检查转出账户余额
    var fromBalance float64
    err = tx.QueryRow("SELECT balance FROM accounts WHERE id = ?", fromID).Scan(&fromBalance)
    if err != nil {
        tx.Rollback()
        return err
    }

    if fromBalance < amount {
        tx.Rollback()
        return fmt.Errorf("余额不足")
    }

    // 扣除转出账户余额
    _, err = tx.Exec("UPDATE accounts SET balance = balance - ? WHERE id = ?", amount, fromID)
    if err != nil {
        tx.Rollback()
        return err
    }

    // 增加转入账户余额
    _, err = tx.Exec("UPDATE accounts SET balance = balance + ? WHERE id = ?", amount, toID)
    if err != nil {
        tx.Rollback()
        return err
    }

    // 记录交易日志
    _, err = tx.Exec(
        "INSERT INTO transactions(from_id, to_id, amount) VALUES(?, ?, ?)",
        fromID, toID, amount,
    )
    if err != nil {
        tx.Rollback()
        return err
    }

    return tx.Commit()
}

func main() {
    db, err := sql.Open("mysql", "username:password@tcp(localhost:3306)/testdb")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    if err = db.Ping(); err != nil {
        log.Fatal(err)
    }

    service := NewTransactionService(db)

    // 执行转账
    err = service.Transfer(1, 2, 100.0)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("转账成功")
}
```

---

## 7. 最佳实践

### 7.1 连接管理

- **使用连接池**：合理配置连接池参数，避免连接泄漏
- **及时关闭连接**：使用 defer 确保连接被正确关闭
- **设置超时**：为数据库操作设置合理的超时时间
- **监控连接池**：定期检查连接池状态，优化连接配置

### 7.2 错误处理

- **检查所有错误**：不要忽略任何数据库操作的错误
- **处理 sql.ErrNoRows**：区分"没有记录"和真正的错误
- **使用事务**：对于需要原子性的操作，使用事务确保数据一致性
- **记录错误日志**：记录详细的错误信息，便于排查问题

### 7.3 性能优化

- **使用预编译语句**：对于重复执行的 SQL 语句，使用预编译语句提高性能
- **合理使用索引**：为查询字段创建合适的索引
- **避免 SELECT ***：只查询需要的字段，减少数据传输量
- **批量操作**：对于大量数据的插入或更新，使用批量操作

### 7.4 安全考虑

- **使用参数化查询**：避免 SQL 注入攻击
- **最小权限原则**：为数据库用户分配最小必要的权限
- **敏感信息保护**：不要在代码中硬编码数据库密码
- **使用 SSL/TLS**：在生产环境中使用加密连接

---

## 8. 常见问题

### 8.1 连接泄漏

**问题**：数据库连接数持续增长，最终导致无法获取连接。

**解决方案**：
- 确保 Rows 和 Stmt 被正确关闭
- 使用 defer 确保资源释放
- 监控连接池状态
- 设置合理的连接超时时间

```go
rows, err := db.Query("SELECT * FROM users")
if err != nil {
    panic(err)
}
defer rows.Close() // 确保关闭 Rows

for rows.Next() {
    // 处理数据
}
```

### 8.2 事务回滚

**问题**：事务中的某个操作失败，但事务没有正确回滚。

**解决方案**：
- 使用 defer 确保事务被回滚或提交
- 在 defer 中检查错误状态
- 使用 recover 处理 panic 情况

```go
tx, err := db.Begin()
if err != nil {
    panic(err)
}

defer func() {
    if p := recover(); p != nil {
        tx.Rollback()
        panic(p)
    }
}()

// 执行操作
_, err = tx.Exec("UPDATE users SET age = age + 1 WHERE id = ?", 1)
if err != nil {
    tx.Rollback()
    panic(err)
}

err = tx.Commit()
if err != nil {
    panic(err)
}
```

### 8.3 空值处理

**问题**：数据库中的 NULL 值导致 Scan 操作失败。

**解决方案**：
- 使用 sql.Null 类型处理 NULL 值
- 检查 Valid 字段判断是否为 NULL
- 使用指针类型处理 NULL 值

```go
var name sql.NullString
err := db.QueryRow("SELECT name FROM users WHERE id = ?", 1).Scan(&name)
if err != nil {
    panic(err)
}

if name.Valid {
    fmt.Printf("姓名: %s\n", name.String)
} else {
    fmt.Println("姓名为 NULL")
}
```

### 8.4 上下文超时

**问题**：数据库操作长时间阻塞，影响应用性能。

**解决方案**：
- 为数据库操作设置合理的超时时间
- 使用 context.WithTimeout 创建带超时的上下文
- 处理 context.DeadlineExceeded 错误

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

var name string
err := db.QueryRowContext(ctx, "SELECT name FROM users WHERE id = ?", 1).Scan(&name)
if err != nil {
    if err == context.DeadlineExceeded {
        fmt.Println("查询超时")
    } else {
        panic(err)
    }
}
```

---

## 9. 总结

database/sql 是 Go 语言中操作关系型数据库的标准包，提供了统一的接口和丰富的功能。通过本文的介绍，相信您已经掌握了它的基本使用方法和高级特性。

记住以下几点：

- 选择合适的数据库驱动
- 正确管理数据库连接
- 使用事务确保数据一致性
- 合理处理 NULL 值
- 遵循最佳实践，确保代码质量和性能
- 注意安全考虑，防止 SQL 注入

通过这些实践，您可以充分发挥 database/sql 的优势，为您的应用提供可靠的数据持久化能力。