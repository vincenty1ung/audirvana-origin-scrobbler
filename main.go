package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/audirvana-origin-scrobbler/config"
	"github.com/audirvana-origin-scrobbler/log"
	"github.com/audirvana-origin-scrobbler/scrobbler"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	c := make(chan struct{})
	_ = log.LogInit(config.ConfigObj.Log.Path, config.ConfigObj.Log.Level, c)
	scrobbler.InitLastfmApi(
		config.ConfigObj.Lastfm.ApiKey,
		config.ConfigObj.Lastfm.SharedSecret,
		config.ConfigObj.Lastfm.UserLoginToken, false, "", "",
	)

	go scrobbler.CheckPlayingTrack(c)

	select {
	case <-ctx.Done():
		fmt.Println("system exiting")
		close(c)
		stop()
	}
}
