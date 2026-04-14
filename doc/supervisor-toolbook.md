# Supervisor 工具使用指南

> 基于 Supervisor 官方文档和实践经验整理，面向需要进程管理的开发者和运维人员。
> 
> 目标：读完后能独立安装、配置并使用 Supervisor 管理进程，解决常见问题。

---

## 1. Supervisor 是什么

Supervisor 是用 Python 开发的一套通用的进程管理程序，能将一个普通的命令行进程变为后台 daemon，并监控进程状态，异常退出时能自动重启。

Supervisor 的核心能力包括：

- 进程后台运行
- 自动重启异常退出的进程
- 统一管理多个进程
- 提供命令行和 Web 界面操作
- 详细的日志记录

适合使用 Supervisor 的场景：

- 需要常驻运行的脚本或服务
- 后台任务处理
- 微服务架构中的服务管理
- 任何需要保证持续运行的进程

官网：[supervisord.org](http://supervisord.org/)

---

## 2. 安装与环境准备

### 2.1 推荐：使用 apt-get 安装（全局生效）

在 Debian/Ubuntu 系统上，可以直接通过 apt 安装：

```bash
apt-get install supervisor
```

安装完成后，Supervisor 会自动配置为系统服务，并且全局可用。

### 2.2 使用 pip 安装

使用 pip 安装：

```bash
pip install supervisor
```

### 2.3 pip 安装到虚拟环境的问题

如果您使用 pip 安装到虚拟环境中，可能会遇到全局不生效的问题。这是因为虚拟环境中的命令默认只在该环境内可用。

**解决方案：**

1. **方法一：激活虚拟环境后使用**
   ```bash
   source /path/to/venv/bin/activate
   supervisord
   ```

2. **方法二：使用绝对路径**
   ```bash
   /path/to/venv/bin/supervisord
   ```

3. **方法三：将虚拟环境的 bin 目录添加到系统 PATH**
   ```bash
   export PATH=/path/to/venv/bin:$PATH
   ```

---

## 3. 配置文件

### 3.1 主配置文件

- **apt-get 安装**：主配置文件位于 `/etc/supervisor/supervisord.conf`
- **pip 安装**：需要手动生成配置文件
  ```bash
  mkdir -p /etc/supervisor
  echo_supervisord_conf > /etc/supervisor/supervisord.conf
  ```

### 3.2 子进程配置文件

子进程配置文件通常放在 `/etc/supervisor/conf.d/` 目录下，一个进程对应一个配置文件，后缀为 `.conf`。

---

## 4. 启动 Supervisor

### 4.1 apt-get 安装的启动方式

```bash
/etc/init.d/supervisor start
# 或
systemctl start supervisor
```

### 4.2 pip 安装的启动方式

```bash
# 默认查找配置文件路径：/usr/etc/supervisord.conf, /usr/supervisord.conf, supervisord.conf, etc/supervisord.conf, /etc/supervisord.conf, /etc/supervisor/supervisord.conf
supervisord

# 或指定配置文件
supervisord -c /etc/supervisor/supervisord.conf
```

---

## 5. 常用命令

### 5.1 查看状态

```bash
supervisorctl status
```

### 5.2 管理进程

```bash
# 重新读取配置文件
supervisorctl reread

# 更新进程组
supervisorctl update

# 启动所有进程
supervisorctl start all

# 启动指定进程
supervisorctl start <process_name>

# 停止所有进程
supervisorctl stop all

# 停止指定进程
supervisorctl stop <process_name>

# 重启所有进程
supervisorctl restart all

# 重启指定进程
supervisorctl restart <process_name>

# 重启 supervisord 服务
supervisorctl reload
```

---

## 6. 配置文件详解

### 6.1 主配置文件示例

```ini
[unix_http_server]
file=/tmp/supervisor.sock   ; socket 文件路径

[supervisord]
logfile=/tmp/supervisord.log ; 日志文件
logfile_maxbytes=50MB        ; 日志文件最大大小
logfile_backups=10           ; 日志文件备份数
loglevel=info                ; 日志级别
pidfile=/tmp/supervisord.pid ; pid 文件
nodaemon=false               ; 是否在前台运行
minfds=1024                  ; 最小文件描述符数
minprocs=200                 ; 最小进程数

[rpcinterface:supervisor]
supervisor.rpcinterface_factory = supervisor.rpcinterface:make_main_rpcinterface

[supervisorctl]
serverurl=unix:///tmp/supervisor.sock ; 连接 supervisor 的 URL

[include]
files = /etc/supervisor/conf.d/*.conf ; 包含的子配置文件
```

### 6.2 子进程配置文件示例

```ini
[program:echo_time]
command=sh /tmp/echo_time.sh    ; 执行命令
priority=999                    ; 启动优先级（默认 999）
autostart=true                  ; supervisord 启动时自动启动
autorestart=true                ; 进程意外退出时自动重启
startsecs=10                    ; 进程必须运行的秒数才算正常启动
startretries=3                  ; 启动失败的最大重试次数
exitcodes=0,2                   ; 预期的退出码
stopsignal=QUIT                 ; 停止信号
stopwaitsecs=10                 ; 发送停止信号后等待的最大秒数
user=root                       ; 运行进程的用户
log_stdout=true                 ; 记录标准输出
log_stderr=true                 ; 记录标准错误
logfile=/tmp/echo_time.log      ; 日志文件
logfile_maxbytes=1MB            ; 日志文件最大大小
logfile_backups=10              ; 日志文件备份数
stdout_logfile_maxbytes=20MB    ; 标准输出日志文件最大大小
stdout_logfile_backups=20       ; 标准输出日志文件备份数
stdout_logfile=/tmp/echo_time.stdout.log ; 标准输出日志文件路径
```

---

## 7. 实战示例

### 7.1 创建一个简单的测试脚本

创建 `/tmp/echo_time.sh` 文件：

```bash
#!/bin/bash

while true; do
    echo `date +%Y-%m-%d,%H:%M:%S`
    sleep 2
done
```

赋予执行权限：

```bash
chmod +x /tmp/echo_time.sh
```

### 7.2 创建子进程配置文件

创建 `/etc/supervisor/conf.d/echo_time.conf` 文件：

```ini
[program:echo_time]
command=sh /tmp/echo_time.sh
autostart=true
autorestart=true
startsecs=10
startretries=3
exitcodes=0,2
stopsignal=QUIT
stopwaitsecs=10
user=root
log_stdout=true
log_stderr=true
logfile=/tmp/echo_time.log
logfile_maxbytes=1MB
logfile_backups=10
stdout_logfile_maxbytes=20MB
stdout_logfile_backups=20
stdout_logfile=/tmp/echo_time.stdout.log
```

### 7.3 启动进程

```bash
supervisorctl reread
supervisorctl update
```

### 7.4 查看运行状态

```bash
supervisorctl status
# 输出：echo_time                        RUNNING   pid 12345, uptime 0:00:10

# 查看日志
tail -f /tmp/echo_time.stdout.log
```

---

## 8. Web 界面操作

### 8.1 开启 Web 界面

修改主配置文件 `/etc/supervisor/supervisord.conf`，取消注释以下部分：

```ini
[inet_http_server]         ; inet (TCP) server disabled by default
port=*:9001                ; ip_address:port specifier, *:port for all iface
username=user              ; default is no username (open server)
password=123               ; default is no password (open server)
```

### 8.2 重启 Supervisor

```bash
supervisorctl reload
```

### 8.3 访问 Web 界面

在浏览器中访问 `http://your-server-ip:9001`，输入用户名和密码即可。

---

## 9. 常见问题

### 9.1 pip 安装后全局不生效

**问题**：使用 pip 安装到虚拟环境后，全局无法使用 supervisor 命令。

**解决方案**：
- 激活虚拟环境后使用
- 使用绝对路径执行
- 将虚拟环境的 bin 目录添加到系统 PATH

### 9.2 缺少 meld3 依赖

**问题**：安装后执行 `supervisor -V` 出现 `pkg_resources.DistributionNotFound: meld3>=0.6.5` 错误。

**解决方案**：手动安装 meld3

```bash
wget https://pypi.python.org/packages/source/m/meld3/meld3-1.0.2.tar.gz
tar -zxf meld3-1.0.2.tar.gz
cd meld3-1.0.2
python setup.py install
```

### 9.3 pip 不可用

**问题**：执行 pip 命令出现 `ImportError: No module named sysconfig` 错误。

**解决方案**：重新安装 python-setuptools

```bash
rm -rf /usr/lib/python2.6/site-packages/pkg_resources*
yum reinstall python-setuptools
```

---

## 10. 总结

Supervisor 是一个功能强大的进程管理工具，通过它可以轻松实现进程的后台运行、自动重启等功能。推荐使用 apt-get 方式安装，这样可以全局生效且配置更简单。如果使用 pip 安装到虚拟环境，需要注意全局使用的问题。

通过本文的介绍，相信您已经掌握了 Supervisor 的基本使用方法，可以在实际项目中应用它来管理需要常驻运行的进程。