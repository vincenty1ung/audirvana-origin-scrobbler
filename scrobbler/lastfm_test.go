package scrobbler

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/audirvana-origin-scrobbler/config"
	"github.com/audirvana-origin-scrobbler/log"
)

func init() {
	c := make(chan struct{})
	config.InitConfig("../config/config.yaml")
	_ = log.LogInit(config.ConfigObj.Log.Path, config.ConfigObj.Log.Level, c)
	InitLastfmApi(context.Background(),
		config.ConfigObj.Lastfm.ApiKey, config.ConfigObj.Lastfm.SharedSecret, "", true,
		config.ConfigObj.Lastfm.UserUsername, config.ConfigObj.Lastfm.UserPassword,
	)
}

func TestPushTrackScrobble(t *testing.T) {
	parse, err2 := time.Parse(time.DateTime, "2025-01-15 15:04:05")
	if err2 != nil {
		fmt.Println(err2)
	}
	unix := parse.Unix()
	unix2 := parse.UTC().Unix()
	if unix != unix2 {
		fmt.Println(unix, unix2)
	}
	res, err := PushTrackScrobble(context.Background(),
		&PushTrackScrobbleReq{
			base:   base{},
			Artist: "声音玩具",
			// AlbumArtist:        "声音玩具",
			Track:              "抚琴小夜曲",
			Album:              "爱是昂贵的",
			TrackNumber:        6,
			Timestamp:          unix,
			MusicBrainzTrackID: "1fa14539-2851-4982-bfda-4a78ad390a36",
			Context:            "",
			StreamId:           0,
			Duration:           479,
			ChosenByUser:       0,
			Sk:                 "",
		},
	)
	fmt.Println(res)
	if err != nil {
		t.Error(err)
	}
}

func TestPushTrackScrobbleReq(t *testing.T) {
	req := PushTrackScrobbleReq{
		base:   base{},
		Artist: "声音玩具",
		// AlbumArtist:        "声音玩具",
		Track:              "抚琴小夜曲",
		Album:              "爱是昂贵的",
		TrackNumber:        6,
		Timestamp:          time.Now().Unix(),
		MusicBrainzTrackID: "1fa14539-2851-4982-bfda-4a78ad390a36",
		Context:            "",
		StreamId:           0,
		Duration:           479,
		ChosenByUser:       0,
		Sk:                 "",
	}
	res, err := req.ToMap()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(res)
}
