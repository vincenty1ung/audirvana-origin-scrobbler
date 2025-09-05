# 项目概述

这是一个用 Go 语言编写的音乐播放记录同步工具，主要用于将 Audirvana 和 Roon 播放的音乐曲目同步（Scrobble）到 Last.fm。项目通过定时检查播放状态，获取当前播放曲目的信息，并在满足一定条件（如播放进度达到 55%）时将曲目信息上报到 Last.fm。

## 新增功能概览

1. **本地数据库存储**: 使用 SQLite 通过 GORM 实现本地数据持久化，存储播放记录和播放统计。
2. **播放记录追踪**: 记录每次播放的详细信息，包括艺术家、专辑、曲目、播放时间等。
3. **播放统计**: 统计每首曲目的播放次数，使用乐观锁机制保证并发安全。
4. **数据同步**: 将未同步到 Last.fm 的播放记录进行同步，并标记同步状态。

## 主要技术栈

- **语言**: Go 1.24
- **依赖管理**: Go Modules
- **数据库**: SQLite (通过 GORM)
- **链路跟踪**: OpenTelemetry
- **主要依赖库**:
  - `github.com/spf13/cobra` - 命令行接口
  - `github.com/spf13/viper` - 配置管理
  - `go.uber.org/zap` - 日志记录
  - `github.com/shkh/lastfm-go` - Last.fm API 客户端
  - `github.com/milindmadhukar/go-musixmatch` - Musixmatch API 客户端 (当前被注释)
  - `github.com/andybrewer/mack` - AppleScript 执行
  - `gorm.io/gorm` - ORM 框架
  - `gorm.io/driver/sqlite` - SQLite 驱动
  - `go.opentelemetry.io/otel` - OpenTelemetry SDK

## 项目架构

- `main.go`: 程序入口，使用 Cobra 设置命令行参数并启动服务。
- `config/`: 配置管理模块，使用 Viper 解析 `config.yaml`。
- `scrobbler/`: 核心逻辑模块，负责检查播放状态、获取曲目信息、与 Last.fm 交互。
- `audirvana/`: 与 Audirvana 应用交互的模块，通过 AppleScript 获取播放信息。
- `roon/`: 与 Roon 应用交互的模块 (具体实现未在本次分析中详细查看)。
- `log/`: 日志模块，基于 Zap 实现。
- `model/`: 数据模型和数据库操作模块，使用 GORM 实现。
- `shell/`: 包含用于构建、启动和停止服务的 shell 脚本。
- `storage/`: 本地数据存储目录。
- `telemetry/`: 链路跟踪模块，集成 OpenTelemetry 实现分布式追踪。

## 数据模型

### TrackPlayRecord (播放记录)

存储每次播放的详细信息：

- ID: 主键
- Artist: 艺术家
- AlbumArtist: 专辑艺术家
- Track: 曲目名
- Album: 专辑名
- Duration: 持续时间
- PlayTime: 播放时间
- Scrobbled: 是否已同步到 Last.fm
- MusicBrainzID: MusicBrainz ID
- TrackNumber: 音轨号
- Source: 数据来源（Audirvana 或 Roon）
- CreatedAt: 创建时间
- UpdatedAt: 更新时间

### TrackPlayCount (播放统计)

统计每首曲目的播放次数：

- ID: 主键
- Artist: 艺术家
- Album: 专辑名
- Track: 曲目名
- PlayCount: 播放次数
- Version: 乐观锁版本号
- CreatedAt: 创建时间
- UpdatedAt: 更新时间

## 核心功能实现

### 数据库初始化

在 `model/db.go` 中实现数据库连接和初始化，使用自定义日志记录器集成 zap 和 OpenTelemetry。

### 链路跟踪

在 `telemetry/telemetry.go` 中实现 OpenTelemetry 的初始化和配置，包括 tracer provider 和 exporter 的设置。在各个关键模块中创建 span 来跟踪请求和操作的执行过程。

### 播放记录存储

在 `model/track_play_record.go` 中实现播放记录的插入、更新和查询功能。

### 播放统计

在 `model/track_play_count.go` 中实现播放次数的增加和查询功能，使用乐观锁机制处理并发更新。

# 构建和运行

## 配置

在运行程序前，需要配置 `config/config.yaml` 文件，填入 Last.fm 和 Musixmatch 的 API 密钥等信息。

## 构建

```bash
go build
```

## 运行

````bash
./lastfm-scrobbler
```lastfm-scrobbler

## 使用脚本运行
lastfm-scrobbler
项目提供了 shell 脚本来简化构建和运行过程：

```bashlastfm-scrobbler
# 构建 launchctl 服务
sh shell/script/build_lastfm-scrobblers_launchctl.sh

# 启动服务
sh shell/script/start_lastfm-scrobblers.sh
lastfm-scrobbler
# 停止服务
sh shell/script/stop_lastfm-scrobblers.sh
````

## 查看日志

```bash
tail -f logs/go_lastfm-scrobbler.log
```

# 开发约定

- **代码风格**: 遵循 Go 语言惯用风格。
- **日志**: 使用 `go.uber.org/zap` 进行日志记录，区分不同日志级别。
- **配置**: 使用 `github.com/spf13/viper` 管理配置，配置文件为 `config/config.yaml`。
- **命令行接口**: 使用 `github.com/spf13/cobra` 构建命令行接口。
- **数据库**: 使用 `gorm.io/gorm` 作为 ORM 框架，使用 SQLite 作为本地存储。

## 特性开发与记忆协议

- **特性清单**: 为任何新模块或特性，必须在 `memory/{date}` 目录内创建一个特性清单文件，详细说明其范围、功能和实现要点。文件名应具有描述性，如 `feature_name_feature_manifest.md`。
- **记忆索引**: 创建新特性后，必须在中央 `memory_index.md` 文件中添加一个条目，包含:
  - **日期**: 添加特性的日期。
  - **特性摘要**: 一句话总结新特性。
  - **链接**: 指向该特性的特性清单文件的链接。
- **记忆扩展**: 如果主 `CURSOR.md` 文件变得过于庞大，应为特定领域创建补充的 markdown 文件，并在主文件中链接到它们，以保持清晰。
- **日期分类管理**: 特性清单文件应按创建日期归档到 `memory/{date}` 目录中，以便更好地组织和管理。
