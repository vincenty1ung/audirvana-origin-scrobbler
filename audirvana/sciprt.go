package audirvana

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"go.uber.org/zap"

	"github.com/audirvana-origin-scrobbler/applesciprt"
	"github.com/audirvana-origin-scrobbler/common"
	alog "github.com/audirvana-origin-scrobbler/log"
)

type (
	TrackBase struct {
		TrackID    string
		Title      string
		Album      string
		Artist     string
		Duration   int64
		Position   float64
		Url        string
		AirfoiLogo string
	}
	TrackInfo struct {
		TrackBase
	}
)

func IsRunning() bool {
	tell, err := applesciprt.Tell(
		common.AppSystemEvents, fmt.Sprintf(
			`set listApplicationProcessNames to name of every application process
			if listApplicationProcessNames contains "%s" then
				set AUDIRVANA_RUNNING_STATE to "true"
			else
				set AUDIRVANA_RUNNING_STATE to "false"
			end if`, common.AppAudirvanaOrigin,
		),
	)
	if err != nil {
		return false
	}
	parseBool, err := strconv.ParseBool(tell)
	if err != nil {
		alog.Logger.Warn("err:", zap.Error(err))
		return false
	}
	return parseBool
}

func GetState() (playerState common.PlayerState, err error) {
	result, err := applesciprt.Tell(common.AppAudirvanaOrigin, `set audirvanaState to get player state`)
	if err != nil {
		alog.Logger.Warn("err:", zap.Error(err))
		return "", err
	}
	return common.PlayerState(result), nil
}

func GetNowPlayingTrackInfo() *TrackInfo {
	tell, err := applesciprt.Tell(
		common.AppAudirvanaOrigin,
		`set playingTrack to playing track title
	set playingAlbum to playing track album
	set playingArtist to playing track artist
	set playingDuration to playing track duration
	set playingDuration to playingDuration as string
	set playingPosition to player position
	set playingPosition to playingPosition as string
	set playingUrl to playing track url
	set result to playingTrack & "|" & playingAlbum & "|" & playingArtist & "|" & playingDuration & "|" & playingPosition & "|" & playingUrl`,
	)
	if err != nil {
		alog.Logger.Warn("err:", zap.Error(err))
		return nil
	}
	split := strings.Split(tell, "|")
	info := &TrackInfo{
		TrackBase: TrackBase{},
	}
	for i, s := range split {
		switch i {
		case 0:
			info.Title = strings.TrimSpace(s)
		case 1:
			info.Album = strings.TrimSpace(s)
		case 2:
			info.Artist = strings.TrimSpace(s)
		case 3:
			parseInt, err := strconv.ParseInt(strings.TrimSpace(s), 10, 64)
			if err != nil {
				alog.Logger.Warn("err:", zap.Error(err))
				return nil
			}
			info.Duration = parseInt
		case 4:
			parseInt, err := strconv.ParseFloat(strings.TrimSpace(s), 32)
			if err != nil {
				alog.Logger.Warn("err:", zap.Error(err))
				return nil
			}
			info.Position = parseInt
		case 5:
			unescape, err := url.PathUnescape(s)
			if err != nil {
				alog.Logger.Warn("err:", zap.Error(err))
				return nil
			}
			info.Url = strings.TrimSpace(unescape)
		case 6:
			// info.AirfoiLogo = s
		}
	}
	return info
}
