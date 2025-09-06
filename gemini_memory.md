# Project Memory: lastfm-scrobbler

## 1. Project Overview

`lastfm-scrobbler` is a Go-based background service for macOS. Its primary function is to monitor music playback from **Audirvana** and **Roon** players and "scrobble" the track information to the user's [Last.fm](https://www.last.fm/) profile. This automates the process of logging listening history for users of these high-fidelity audio players.

## 2. Core Logic & Technologies

- **Language**: Go
- **Concurrency**: The application runs two concurrent goroutines to monitor Audirvana and Roon independently.
- **Player Interaction**: It uses **AppleScript** to communicate with Audirvana and Roon to get the player state (playing, paused) and current track details (artist, album, title, duration, position). The `github.com/andybrewer/mack` library is used to execute AppleScript from Go.
- **Last.fm API**: It uses the `github.com/shkh/lastfm-go` library to send "Now Playing" updates and scrobble tracks to the Last.fm API.
- **Configuration**: The application is configured via a central `config/config.yaml` file, managed by the `github.com/spf13/viper` library.
- **CLI**: `github.com/spf13/cobra` is used to build the command-line interface.
- **Logging**: Structured logging is implemented using `go.uber.org/zap`.
- **Daemonization**: The project includes shell scripts (`shell/`) and `.plist` files (`shell/launch/`) to run the application as a persistent background service on macOS using `launchd`.

## 3. How to Use

### Configuration (Crucial Step)

1.  Edit the `config/config.yaml` file.
2.  Fill in the `lastfm` section with your own API key, shared secret, username, and password. API keys can be obtained from the [Last.fm API website](https://www.last.fm/api/account/create).

### Running the Application

There are two main ways to run the project:

1.  **Manual/Debug Mode**:

    - Build the binary: `go builastfm-scrobbler
    - Run the executable: `./audirvana-lastfm-scrobbler
    - Monitor logs: `tail -f .logs/go_lastfm-scrobbler.log`

2.  **Background Service (Recommended)**:
    - Use the provided shell scripts to manage the `lastfm-scrobbler
    - **Build & Install**: `sh shell/script/build_lastfm-scrobblers_launchctl.sh`
    - **Start Service**: `sh shell/script/start_lastfm-scrobblers.sh`
    - **Stop Service**: `sh shell/script/stop_lastfm-scrobblers.sh`

## 4. Key Project Structure

```
/
├── main.go               # Main application entry point, handles startup and shutdown.
├── go.mod                # Go module dependencies.
├── config/
│   ├── config.go         # Go struct definitions for the config.
│   └── config.yaml       # User configuration file (API keys, etc.).
├── scrobbler/            # Core scrobbling logic.
│   ├── lastfm.go         # Handles all communication with the Last.fm API.
│   └── track_check_playing.go # Contains the main monitoring loops for Audirvana and Roon.
├── applesciprt/
│   └── applesciprt.go    # Wrappers for executing AppleScript commands.
├── log/
│   └── log.go            # Logging setup and initialization.
├── shell/                # Scripts for building and managing the service.
│   ├── script/           # Build, start, and stop scripts.
│   └── launch/           # `launchd` .plist templates.
└── study.md              # Detailed project documentation and learning guide.
```

## 5. Development Guidelines

- **Localization (Important)**: **The primary development environment is Chinese.** To maintain consistency, logs, comments, and documentation should prioritize using **Chinese**.
- **Go Style Guide**: Code must adhere to Go's idiomatic style (formatting, naming, comments). Follow the guidelines from the [Uber Go Style Guide (Chinese Translation)](https://github.com/xxjwxc/uber_go_guide_cn/blob/master/README.md).
- **API Design**: All APIs must follow RESTful standards, using appropriate HTTP verbs and resource-based paths.
- **Database Usage**: For future database integration (MySQL, Redis), prioritize performance by using space-for-time trade-offs (e.g., caching, indexing).
- **Testing**: New business logic must be covered by unit tests.
- **Logging**: Logs must be accurate and detailed. Use log levels (`info`, `error`, `warn`, `debug`) appropriately to describe the event.

## 6. Feature Development and Memory Protocol

- **Feature Manifest**: For any new module or feature, a `feature_manifest.md` file must be created within the module's directory, detailing its scope, functionality, and implementation notes.
- **Memory Indexing**: Upon creating a new feature, an entry must be added to a central `memory_index.md` file. This entry must include:
  - **Date**: The date the feature was added.
  - **Feature Summary**: A concise, one-sentence summary of the new feature.
  - **Link**: A link to the feature's `feature_manifest.md`.
- **Memory Scalability**: If the primary `gemini_memory.md` becomes too large, create supplementary markdown files for specific domains (e.g., `database_memory.md`, `api_design_memory.md`) and link to them from the main memory file to maintain clarity.
