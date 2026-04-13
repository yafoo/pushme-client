# PushMeClient

<div align="center">

💬 Windows 桌面消息通知客户端 | 实时接收手机推送消息

[![GitHub release](https://img.shields.io/github/v/release/yafoo/pushme-client)](https://github.com/yafoo/pushme-client/releases)
[![Go Version](https://img.shields.io/badge/Go-1.21.8+-blue)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)
[![Platform](https://img.shields.io/badge/platform-Windows%2010%2B-lightgrey)](https://www.microsoft.com/windows)

</div>

---

## 📖 项目简介

PushMeClient 是一款专为 Windows (Win10+) 设计的桌面消息通知客户端，能够实时接收并展示来自 [PushMe](https://push.i-i.me/) APP 转发的各类消息。让您在工作期间无需频繁查看手机，即可在电脑上快速浏览短信、应用通知等重要信息。

### ✨ 核心特性

- 📨 **多消息类型支持**：文本、Markdown、数据消息、小屏展示、图表消息
- 📊 **动态图表渲染**：支持实时更新的图表消息展示（基于 ECharts）
- 🔌 **灵活的接口协议**：支持 GET、POST/FORM、POST/JSON 多种请求方式
- 🧩 **插件系统**：可扩展的插件管理机制
- 🏠 **自建服务支持**：可连接自建的 PushMeServer 服务
- 🚀 **系统集成**：开机自启动、Windows 原生通知推送
- ⚙️ **可视化配置**：友好的图形化设置界面
- 💻 **轻量高效**：基于 WebView 技术，内存占用低，界面清晰流畅

### 🔄 适用场景

除了接收 PushMe APP 的消息转发，还可以对接其他应用的消息推送，例如：
- [SmsForwarder](https://github.com/pppscn/SmsForwarder) - 短信转发器
- 自定义应用通知推送
- 服务器监控告警消息

---

## 📦 快速开始

### 下载安装

前往 [Releases](https://github.com/yafoo/pushme-client/releases) 页面下载最新版本，解压后即可直接运行，无需安装。

### 基础配置

1. **启动客户端**：双击运行 PushMeClient.exe
2. **获取局域网 IP**：软件默认监听 `3200` 端口
3. **配置 PushMe APP**：
   - 打开 PushMe 安卓客户端
   - 进入消息转发设置
   - 填写转发 URL：`http://您的局域网IP:3200`

### 公网访问（可选）

如果您的电脑拥有公网 IP，可以进行以下配置实现远程接收：
1. 在路由器中配置端口映射（将外部端口映射到电脑的 3200 端口）
2. 在 PushMe APP 中配置公网 IP 地址
3. 确保防火墙允许 3200 端口的入站连接

> ⚠️ **注意**：PushMeClient 目前不能直接接收 PushMe 官网接口的消息推送，必须通过 PushMe APP 转发或自建 PushMeServer 服务。如需独立使用，请参考 [PushMeServer](https://github.com/yafoo/pushme-server)。

---

## 🎯 特色功能

### 📈 图表消息

PushMeClient 独有功能，支持动态图表消息展示：
- 实时数据更新
- 多种图表类型（折线图、柱状图、饼图等）
- 支持全屏查看和交互操作

### 🧩 插件系统

通过插件扩展更多功能：
- 自定义消息处理逻辑
- 集成第三方服务
- 灵活的消息过滤和转换

---

## 🛠️ 技术架构

### 技术栈

| 层级 | 技术选型 |
|------|----------|
| 后端语言 | Go 1.21.8+ |
| UI 框架 | webview/webview_go |
| 前端框架 | Vue.js 3 (Global Build) |
| 图表库 | ECharts 5 |
| Markdown 解析 | markdown-it |
| 数据库 | SQLite + GORM |
| 系统交互 | go-ole (Windows COM)、go-toast |
| 安全加固 | DOMPurify (XSS 防护) |

### 版本演进

- **v1.0**：纯 Go + FyneUI 方案
  - ❌ 字体显示模糊
  - ❌ 内存占用较高
  - ❌ 开发体验不佳

- **v2.0**：Go + WebView 方案（当前版本）
  - ✅ 界面清晰流畅
  - ✅ 内存占用优化
  - ✅ 开发效率提升
  - ✅ 参考了 [http-win-notice](https://github.com/shanghaobo/http-win-notice) 的优秀实现

---

## 📸 界面预览

<div align="center">

| 主页 | 消息列表 |
|------|----------|
| ![主页](https://github.com/user-attachments/assets/1c6a9a8b-ac09-4d3e-9988-f00c2afe59fa) | ![消息](https://github.com/user-attachments/assets/cc928c04-3821-483e-a44c-93e7a81e96c3) |

| 小屏展示 | Markdown 渲染 |
|----------|---------------|
| ![小屏](https://github.com/user-attachments/assets/8e80d762-1e19-4f12-a58f-b832a9a252b8) | ![markdown](https://github.com/user-attachments/assets/7a23640d-3077-4c7e-bf3d-ef22e1816579) |

| 图表消息 | 系统通知 |
|----------|----------|
| ![图表](https://github.com/user-attachments/assets/7e5d619e-b21c-422a-9052-e325567c13f9) | ![通知](https://github.com/user-attachments/assets/6e353862-a75d-42e9-99cd-8d2542498776) |

| 插件管理 | 设置界面 |
|----------|----------|
| ![插件](https://github.com/user-attachments/assets/655dbb7b-1a3a-47c8-b448-19caacfac2cf) | ![设置](https://github.com/user-attachments/assets/9ca54ee2-ad15-4d24-b967-f48f834e172b) |

| 图表全屏 |
|----------|
| ![图表全屏](https://github.com/user-attachments/assets/8e0d6c00-d3a9-49f6-87d3-650b6cb74a83) |

</div>

---

## 🗂️ 项目结构

```
PushMeClient/
├── constant/          # 全局常量定义
├── db/                # 数据库模块
│   ├── msg/           # 消息数据操作
│   ├── plugin/        # 插件数据操作
│   └── db.go          # 数据库初始化
├── server/            # Web 服务器
│   ├── html/          # 前端资源
│   │   ├── static/    # 静态文件 (JS/CSS)
│   │   ├── view/      # HTML 页面模板
│   │   └── html.go    # 资源加载辅助
│   └── server.go      # HTTP 服务器实现
├── utils/             # 工具类
│   ├── request/       # HTTP 请求封装
│   ├── setting/       # 配置管理
│   ├── sys/           # 系统工具
│   └── init.go        # 初始化逻辑
├── win/               # Windows 特定功能
├── makelnk/           # 快捷方式创建（开机自启）
├── main.go            # 程序入口
└── build.bat          # 构建脚本
```

---

## 🔧 开发指南

### 环境要求

- **操作系统**：Windows 10 及以上
- **Go 版本**：>= 1.21.8
- **编译器**：GCC（用于编译 CGO 依赖，如 sqlite3 和 webview）

### 本地开发

```bash
# 1. 克隆仓库
git clone https://github.com/yafoo/pushme-client.git
cd pushme-client

# 2. 下载依赖
go mod tidy

# 3. 运行程序
go run main.go
```

### 编译构建

```bash
# 使用提供的构建脚本
build.bat

# 或手动构建
go build -o PushMeClient.exe .
```

> 💡 确保 `CGO_ENABLED=1` 以支持 sqlite3 和 webview 的编译。

---

## 🤝 贡献指南

欢迎提交 Issue 和 Pull Request！

- 🐛 发现 Bug？请提交 [Issue](https://github.com/yafoo/pushme-client/issues)
- 💡 有新想法？欢迎讨论
- 🔧 改进代码？请提交 PR

---

## 📝 特别说明

> ⚠️ **开发者说明**  
> 本项目是作者学习 Go 语言的实践作品，代码质量仍有提升空间。如果您发现任何问题或有改进建议，非常欢迎提出宝贵意见！🙏

---

## 📄 开源协议

本项目采用 MIT 协议开源，详见 [LICENSE](LICENSE) 文件。

---

## 🌟 支持项目

如果这个项目对您有帮助，欢迎给个 Star ⭐ 鼓励一下！

**仓库地址：**
- GitHub：https://github.com/yafoo/pushme-client
- Gitee：https://gitee.com/yafu/pushme-client
