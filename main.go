package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/lastfm-scrobbler/api"
	"github.com/lastfm-scrobbler/cmd"
	"github.com/lastfm-scrobbler/conf
	"github.com/lastfm-scrobbler/log"
	"github.com/lastfm-scrobbler/model"
l"
	"github.com/lastfm-scrobbler/scrobbl
	"github.com/lastfm-scrobbler/telemetry"
)

var (
	configFile = new(string)
	isMobile   = new(bool)
)

func main() {
	rootCmd := NewCommand("lastfm-scrobbler", "", "")
	// command.SetHelpTemplate("使用-c 设置配置文件路径\n使用-m 设置true/false")
	rootCmd.Version = "1.0.0"
	rootCmd.Args = cobra.NoArgs
	rootCmd.RunE = func(cmd *cobra.Command, args []string) error { return initServer() }

	flags := rootCmd.Flags()
	flags.SortFlags = false
	flags.StringVarP(configFile, "config", "c", "config/config.yaml", "config file")
	flags.BoolVarP(isMobile, "mobile", "m", false, "it a mobile")

	// Add sync-records subcommand
	rootCmd.AddCommand(newSyncRecordsCommand())

	// Add memory-tool subcommand
	rootCmd.AddCommand(newMemoryToolCommand())

	cobra.CheckErr(rootCmd.Execute())
}

func newSyncRecordsCommand() *cobra.Command {
	return cmd.NewSyncRecordsCommand()
}

func newMemoryToolCommand() *cobra.Command {
	return cmd.NewMemoryToolCommand()
}

func initServer() error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	c := make(chan struct{})
	config.InitConfig(*configFile)
	_ = log.LogInit(config.ConfigObj.Log.Path, config.ConfigObj.Log.Level, c)

	// Initialize telemetry
	if err := telemetry.Init(config.ConfigObj.Telemetry); err != nil {
		return fmt.Errorf("failed to initialize telemetry: %w", err)
	}
	defer telemetry.Shutdown(context.Background())

	// Initialize database
	if err := model.InitDB(config.ConfigObj.Database.Path); err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}

	// Start HTTP server in a separate goroutine
	go api.StartHTTPServer(ctx, config.ConfigObj.Telemetry.Name)

	_ = run(c)
	<-ctx.Done()
	fmt.Println("system exiting")
	close(c)
	return nil
}

func run(c <-chan struct{}) error {
	ctx := context.Background()
	scrobbler.InitLastfmApi(
		ctx,
		config.ConfigObj.Lastfm.ApiKey,
		config.ConfigObj.Lastfm.SharedSecret,
		config.ConfigObj.Lastfm.UserLoginToken,
		*isMobile,
		config.ConfigObj.Lastfm.UserUsername,
		config.ConfigObj.Lastfm.UserPassword,
	)

	// musixmatch.InitMxmClient(config.ConfigObj.Musixmatch.ApiKey)
	// 音乐检查
	go scrobbler.AudirvanaCheckPlayingTrack(ctx, c)
	go scrobbler.RoonCheckPlayingTrack(ctx, c)
	return nil
}
