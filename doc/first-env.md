# Go开发环境搭建

1. 设备：MacBookPro M5+16+1T
2. VPN：https://vgvg.vg/#/stage/dashboard （https://185.45.192.166:4433/public-Q/#/login）
3. HomeBrew：Mac上软件包管理工具 https://github.com/Homebrew/brew/releases/latest 下载pkg就行
4. Goland：Jetbrains出品的GoIDE
5. DataGrip：Jetbrains出品的数据库管理工具 个人免费使用 完美平替Navicat
6. Git：brew install git 需要配置ssh方式连接
7. ApiFox：Api接口测试工具
8. Proxyman：抓包工具完美平替Charles
9. SwitchHost：本地域名IP映射
10. Snipaste：桌面截图软件
11. TinyRDM: RedisGUI强推比DataGrip好用

# 常见问题

### Mac如何彻底卸载软件

1. AppClean：https://freemacsoft.net/appcleaner/#google_vignette

### Goland无法Debug

1. 排查一下版本是不是2024.*版本(社区issue https://youtrack.jetbrains.com/projects/GO/issues/GO-18407)
2. 建议升级到最新版本(闲鱼有教程)

### Goland无法配置GoSDK

1. 尝试官网下载 https://golang.google.cn/dl/ 需要配置GOPROXY代理
