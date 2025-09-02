# Audirvana Origin Scrobbler 项目记忆

## 1. 项目概述 (Project Overview)

`audirvana-origin-scrobbler` 是一个用 Go 语言编写的后台服务程序,专为 macOS 用户设计。其核心功能是监控 **Audirvana** 和 **Roon** 两款高品质音乐播放器的播放状态,并将正在播放的曲目信息实时同步(Scrobble)到用户的 [Last.fm](https://www.last.fm/) 个人主页。

本项目旨在为使用这些播放器的音乐爱好者提供一个自动化、无缝记录个人听歌历史的解决方案。

## 2. 核心技术与工作原理 (Core Logic & Technologies)

*   **语言 (Language)**: Go
*   **并发模型 (Concurrency)**: 程序利用 Go 的 Goroutine 并发运行两个独立的监控任务,分别对应 Audirvana 和 Roon,互不干扰。
*   **播放器交互 (Player Interaction)**:
    *   通过 `github.com/andybrewer/mack` 库在 Go 中执行 **AppleScript**, 以此与 Audirvana 和 Roon 通信,获取播放器状态和曲目详情。
    *   依赖一个名为 `nowplaying-cli-mac` 的命令行工具来辅助获取 Roon 的播放信息。
*   **Last.fm API**: 使用 `github.com/shkh/lastfm-go` 库向 Last.fm 发送 "正在播放" (Now Playing) 更新和 Scrobble 记录。
*   **配置管理 (Configuration)**: 通过 `config/config.yaml` 文件进行集中配置,并使用 `github.com/spf13/viper` 库进行解析。
*   **命令行接口 (CLI)**: 基于 `github.com/spf13/cobra` 构建。
*   **日志系统 (Logging)**: 采用 `go.uber.org/zap` 实现结构化日志记录。
*   **后台服务 (Daemonization)**: 项目包含 `shell/` 目录下的脚本和 `shell/launch/` 中的 `.plist` 文件,用于将程序注册为 macOS 的 `launchd` 后台服务,实现开机自启和稳定运行。

## 3. 如何使用 (How to Use)

### 步骤一: 配置 (Crucial Step)

1.  找到并用编辑器打开 `config/config.yaml` 文件。
2.  **务必填写 `lastfm` 部分**:
    *   `apiKey` 和 `sharedSecret`: 前往 [Last.fm API 官网](https://www.last.fm/api/account/create) 申请你自己的 API 凭证。
    *   `userUsername`: 你的 Last.fm 用户名。
    *   `userPassword`: 你的 Last.fm 密码。

### 步骤二: 运行程序

提供两种运行方式:

1.  **手动/调试模式 (Manual/Debug Mode)**:
    *   编译: `go build`
    *   运行: `./audirvana-origin-scrobbler`
    *   实时查看日志: `tail -f logs/go_audirvana-origin-scrobbler.log`

2.  **后台服务模式 (Recommended)**:
    *   **构建与安装**: `sh shell/script/build_audirvana-origin-scrobblers_launchctl.sh`
    *   **启动服务**: `sh shell/script/start_audirvana-origin-scrobblers.sh`
    *   **停止服务**: `sh shell/script/stop_audirvana-origin-scrobblers.sh`

## 4. 关键项目结构 (Key Project Structure)

```
/
├── main.go               # 程序主入口,负责初始化、命令行解析和启动服务
├── go.mod                # Go 模块依赖管理文件
├── config/
│   ├── config.go         # 定义了与 YAML 文件对应的 Go 结构体
│   └── config.yaml       # 用户配置文件 (API 密钥等)
├── scrobbler/            # 核心 Scrobble 逻辑
│   ├── lastfm.go         # 封装了与 Last.fm API 的所有交互
│   └── track_check_playing.go # 包含监控 Audirvana 和 Roon 的核心循环
├── applesciprt/
│   └── applesciprt.go    # 封装了执行 AppleScript 的函数
├── log/
│   └── log.go            # 日志系统初始化
├── shell/                # 自动化脚本
│   ├── script/           # 服务的构建、启停脚本
│   └── launch/           # `launchd` 服务所需的 .plist 模板
├── study.md              # 详细的项目学习指南
└── GEMINI.md             # Gemini 使用的项目上下文记忆文件
```

## 5. 开发指南 (Development Guidelines)

*   **本地化 (Localization)**: **项目的主要开发环境为中文**。为保持一致性,日志、注释和文档应优先使用**中文**。
*   **代码风格 (Go Style Guide)**: 代码必须遵循 Go 的惯用风格。推荐参考 [Uber Go 风格指南 (中文版)](https://github.com/xxjwxc/uber_go_guide_cn/blob/master/README.md)。
*   **API 设计 (API Design)**: (若涉及) 所有 API 都应遵循 RESTful 标准。
*   **数据库使用 (Database Usage)**: (若涉及) 对于未来的数据库集成(如 MySQL, Redis),应优先考虑性能,采用空间换时间策略(如缓存、索引)。
*   **测试 (Testing)**: 新增的业务逻辑必须有单元测试覆盖。
*   **日志 (Logging)**: 日志必须准确、详细。合理使用 `info`, `error`, `warn`, `debug` 等日志级别。

## 6. 特性开发与记忆协议 (Feature Development and Memory Protocol)

*   **特性清单 (Feature Manifest)**: 为任何新模块或特性,必须在其目录内创建一个 `feature_manifest.md` 文件,详细说明其范围、功能和实现要点。
*   **记忆索引 (Memory Indexing)**: 创建新特性后,必须在中央 `memory_index.md` 文件中添加一个条目,包含:
    *   **日期**: 添加特性的日期。
    *   **特性摘要**: 一句话总结新特性。
    *   **链接**: 指向该特性的 `feature_manifest.md` 的链接。
*   **记忆扩展 (Memory Scalability)**: 如果主 `GEMINI.md` 文件变得过于庞大,应为特定领域创建补充的 markdown 文件 (例如, `database_memory.md`, `api_design_memory.md`),并在主文件中链接到它们,以保持清晰。