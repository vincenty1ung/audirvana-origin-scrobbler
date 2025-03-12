package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/audirvana-origin-scrobbler/config"
	"github.com/audirvana-origin-scrobbler/log"
	"github.com/audirvana-origin-scrobbler/scrobbler"
)

func main() {
	command := NewCommand("audirvana-origin-scrobbler", "", "")
	// command.SetHelpTemplate("使用-c 设置配置文件路径\n使用-m 设置true/false")
	command.Version = "1.0.0"
	command.Args = cobra.NoArgs
	command.RunE = func(cmd *cobra.Command, args []string) error { return initServer() }

	flags := command.Flags()
	flags.SortFlags = false
	flags.StringVarP(configFile, "config", "c", "config/config.yaml", "config file")
	flags.BoolVarP(isMobile, "mobile", "m", false, "it a mobile")
	cobra.CheckErr(command.Execute())
}

func initServer() error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	c := make(chan struct{})
	config.InitConfig(*configFile)
	_ = log.LogInit(config.ConfigObj.Log.Path, config.ConfigObj.Log.Level, c)
	_ = run(c)
	select {
	case <-ctx.Done():
		fmt.Println("system exiting")
		close(c)
	}
	return nil
}

func run(c <-chan struct{}) error {
	scrobbler.InitLastfmApi(
		config.ConfigObj.Lastfm.ApiKey,
		config.ConfigObj.Lastfm.SharedSecret,
		config.ConfigObj.Lastfm.UserLoginToken,
		*isMobile,
		config.ConfigObj.Lastfm.UserUsername,
		config.ConfigObj.Lastfm.UserPassword,
	)
	// musixmatch.InitMxmClient(config.ConfigObj.Musixmatch.ApiKey)
	// 音乐检查
	go scrobbler.AudirvanaCheckPlayingTrack(c)
	go scrobbler.RoonCheckPlayingTrack(c)
	return nil
}
