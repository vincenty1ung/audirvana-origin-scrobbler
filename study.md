# Audirvana Origin Scrobbler 项目学习指南

## 1. 项目概述

`lastfm-scrobbler` 是一个用 Go 语言编写的后台服务程序,专为 macOS 用户设计。其核心功能是监控 **Audirvana** 和 **Roon** 两款高品质音乐播放器的播放状态,并将正在播放的曲目信息实时同步(Scrobble)到 [Last.fm](https://www.last.fm/) 网站。

对于同时使用 Audirvana/Roon 和 Last.fm 的音乐爱好者来说,这个项目解决了播放器原生不支持或支持不佳的 Scrobble 功能,能够自动化、无缝地记录个人的听歌历史。

## 2. 工作原理

项目的整体逻辑非常清晰,可以分解为以下几个步骤:

1.  **启动与初始化**:

    - 程序启动后,首先会读取位于 `config/config.yaml` 的配置文件。这里包含了 Last.fm 的 API 密钥、用户信息以及日志级别等关键设置。
    - 接着,它会根据配置初始化 Last.fm API 的客户端和 `zap` 日志记录器。

2.  **并发监控**:

    - 程序利用 Go 语言的并发特性,启动两个独立的 Goroutine (可以理解为轻量级线程):
      - 一个用于监控 Audirvana (`scrobbler.AudirvanaCheckPlayingTrack`)。
      - 另一个用于监控 Roon (`scrobbler.RoonCheckPlayingTrack`)。
    - 这两个监控任务并行运行,互不干扰,实现了对两个播放器的同时支持。

3.  **信息获取 (核心)**:

    - 每个监控任务都会进入一个无限循环,定期检查对应播放器的状态。
    - 检查的核心手段是调用 AppleScript。AppleScript 是 macOS 上一种强大的脚本语言,可以与图形界面应用程序进行交互。项目通过 `github.com/andybrewer/mack` 这个库在 Go 代码中执行预设好的 AppleScript 命令,从而向 Audirvana 或 Roon 查询当前播放歌曲的艺术家、专辑、标题和播放进度等信息。

4.  **Scrobble 逻辑**:

    - 获取到曲目信息后,程序会进行一系列判断,例如:
      - 当前是否有歌曲在播放?
      - 播放的歌曲是否和上一首相同?
      - 歌曲的播放时长或进度是否达到了 Last.fm 规定的 Scrobble 阈值(通常是播放一半或超过 4 分钟)?
    - 当满足所有 Scrobble 条件后,程序会调用 `github.com/shkh/lastfm-go` 库,使用用户的凭据,将这首歌曲的信息发送到 Last.fm 的服务器。

5.  **后台运行**:
    - 整个程序被设计为一个守护进程(Daemon)。它会安静地在后台运行,直到接收到系统中断信号 (如 `Ctrl+C` 或关机命令) 时才会优雅地退出。
    - 项目内的 `shell/` 目录和 `.plist` 文件提供了完整的解决方案,用于通过 macOS 的 `launchd` 服务来管理程序的启停,实现开机自启和稳定运行。

## 3. 如何使用

### 步骤一: 配置

这是使用该项目的关键第一步。

1.  找到并打开 `config/config.yaml` 文件。
2.  **重点修改 `lastfm` 部分**:
    - `apiKey` 和 `sharedSecret`: 你需要去 [Last.fm API 官网](https://www.last.fm/api/account/create) 申请自己的 API key。
    - `userUsername`: 你的 Last.fm 用户名。
    - `userPassword`: 你的 Last.fm 密码。

```yaml
lastfm:
  applicationName: lastfm-scrobbler
  apiKey: 9c7d3bxxxxx6bab # <-- 替换成你自己的
  sharedSecret: 80c9e7cxxxxxe0ec3b5 # <-- 替换成你自己的
  registeredTo: vincxxxch1n
  userLoginToken:
  userUsername: vincentch1n # <-- 替换成你自己的
  userPassword: your_xxxxword # <-- 替换成你自己的

log:
  path: ./.logs
  level: info

musixmatch:
  apiKey: 4xxxx5xxxxx81b6654790
```

### 步骤二: 编译与运行

你有两种方式来运行此项目:

**方式一: 手动运行 (用于调试)**

1.  确保你已经安装了 Go 语言环境。
2.  在项目根目录下打开终端,执行编 lastfm-scrobbler
    ```shell
    golastfm-scrobbler
    ```
3.  编译成功后,会生成一个名为 `lastfm-scrobbler` 的可执行文件。运行它:
    ```shell
    ./audirvana-origlastfm-scrobbler
    ```
4.  此时程序已在前台开始运行。你可以通过查看日志来了解其工作状态:
    ```shell
    tail -f .logs/go_lastfm-scrobbler.log
    ```

**方式二: 作为后台服务运行 (推荐的日常使用方式)**

项目提供了非常方便的 shell 脚本来一键完成服务的部署。
lastfm-scrobbler

1.  在项目根目录下打开终端。
2.  执行构建和部署脚本:
    ```shell
    sh shell/script/build_lastfm-scrobblers_launchctl.sh
    ```
    这个脚本会自动编译项目,并生成符合 `launchd` 规范的 `.plist` 配置文件。
3.  启动服务:

    ```shelllastfm-scrobbler
    sh shell/script/start_audirlastfm-scrobbler
    ```

    现在,程序已经在后台运行,并且会随系统开机自动启动。

4.  **其他管理命令**:
    - 停止服务: `sh shell/script/stop_lastfm-scrobblers.sh`
    - 查看日志: `tail -f .logs/go_lastfm-scrobbler.log`

## 4. 项目结构解析

了解项目的文件结构有助于更好地理解代码和进行二次开发。

```
/
├── main.go               # 程序主入口,负责初始化和启动服务
├── go.mod                # Go 模块依赖管理文件
├── config/               # 配置文件目录
│   ├── config.go         # 定义了与 YAML 文件对应的 Go 结构体
│   └── config.yaml       # 用户配置文件,非常重要
├── scrobbler/            # 核心逻辑目录
│   ├── lastfm.go         # 封装了与 Last.fm API 交互的逻辑
│   └── track_check_playing.go # 包含了监控 Audirvana 和 Roon 的核心循环与逻辑
├── applesciprt/          # AppleScript 脚本目录 (实际脚本可能硬编码在 .go 文件中)
│   └── applesciprt.go    # 封装了执行 AppleScript 的函数
├── log/                  # 日志模块
│   └── log.go            # 初始化 zap 日志记录器
├── shell/                # 自动化脚本目录
│   ├── script/           # 存放管理服务的启停、构建脚本
│   └── launch/           # 存放 launchd 使用的 .plist 模板文件
└── README.md             # 项目自述文件
```
