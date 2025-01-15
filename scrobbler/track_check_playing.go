package scrobbler

import (
	"fmt"
	"log"
	"time"

	"go.uber.org/zap"

	"github.com/audirvana-origin-scrobbler/audirvana"
	"github.com/audirvana-origin-scrobbler/common"
	"github.com/audirvana-origin-scrobbler/exec"
	alog "github.com/audirvana-origin-scrobbler/log"
)

const (
	percentScrobble = 0.55
	defaultSleep    = 3
	longSleep       = 60 // 休眠间隔六十秒
	checkCount      = 100
)

var maped = make(map[string]bool)

func CheckPlayingTrack(stop <-chan struct{}) {
	timer := time.NewTicker(time.Second * defaultSleep)
	var (
		now           time.Time
		currentTrack  = ""
		previousTrack = ""
		tmpCount      = 0
	)
	for {
		select {
		case <-timer.C:
			alog.Logger.Debug("Checking playing track..." + time.Now().String())
			tmpCount++
			if tmpCount > checkCount { // 检查100次依旧没有播放检查轮训放大到60秒
				timer.Reset(time.Second * longSleep)
			}
			running := audirvana.IsRunning()
			alog.Logger.Debug("程序运行是否运行", zap.Bool("running", running))
			var audirvanaTrackInfo *audirvana.TrackInfo
			if running {
				audirvanaTrackInfo = nil
				state, _ := audirvana.GetState()
				alog.Logger.Debug("audirvana 播放状态", zap.Any("state", state))
				if state == common.PlayerStatePlaying {
					if tmpCount > checkCount {
						tmpCount = 0
						timer.Reset(time.Second * defaultSleep)
					}
					audirvanaTrackInfo = audirvana.GetNowPlayingTrackInfo()
				}
			}
			if audirvanaTrackInfo != nil {
				currentTrack = audirvanaTrackInfo.Url
				position := audirvanaTrackInfo.Position
				duration := audirvanaTrackInfo.Duration
				if position/float64(duration) > percentScrobble && !maped[currentTrack] {
					// 说明在听歌存在有效数据的
					exiftoolInfo, err := exec.ExiftoolHandle(audirvanaTrackInfo.Url)
					if err != nil {
						alog.Logger.Warn("exec ExiftoolHandle", zap.Error(err))
						continue
					}
					// 标记听歌完成
					pushTrackScrobbleReq := &PushTrackScrobbleReq{
						Artist:      audirvanaTrackInfo.Artist,
						AlbumArtist: audirvanaTrackInfo.Artist,
						Track:       audirvanaTrackInfo.Title,
						Album:       audirvanaTrackInfo.Album,
						Duration:    audirvanaTrackInfo.Duration,
						Timestamp:   now.UTC().Unix(),
					}
					if exiftoolInfo != nil {
						pushTrackScrobbleReq.TrackNumber = exiftoolInfo.GetTrackNumber()
						pushTrackScrobbleReq.MusicBrainzTrackID = exiftoolInfo.GetMusicBrainzTrackId()
					}
					_, err = PushTrackScrobble(pushTrackScrobbleReq)
					if err != nil {
						log.Fatal(err)
					}
					maped[currentTrack] = true
					alog.Logger.Info("标记听歌完成")
				}
				// 上传听歌ing
				if currentTrack != previousTrack {
					// 产生新歌曲
					delete(maped, previousTrack)
					now = time.Now()
					// 说明在听歌存在有效数据的
					exiftoolInfo, err := exec.ExiftoolHandle(audirvanaTrackInfo.Url)
					if err != nil {
						alog.Logger.Warn("exec ExiftoolHandle", zap.Error(err))
						continue
					}
					playingReq := TrackUpdateNowPlayingReq{
						Artist:      audirvanaTrackInfo.Artist,
						AlbumArtist: audirvanaTrackInfo.Artist,
						Track:       audirvanaTrackInfo.Title,
						Album:       audirvanaTrackInfo.Album,
						Duration:    audirvanaTrackInfo.Duration,
					}
					if exiftoolInfo != nil {
						playingReq.TrackNumber = exiftoolInfo.GetTrackNumber()
						playingReq.MusicBrainzTrackID = exiftoolInfo.GetMusicBrainzTrackId()
					}

					alog.Logger.Info("NowPlayingTrackInfo", zap.Any("audirvanaTrackInfo", audirvanaTrackInfo))
					err = TrackUpdateNowPlaying(&playingReq)
					if err != nil {
						alog.Logger.Warn("TrackUpdateNowPlaying", zap.Error(err))
						continue
					}
				}
				previousTrack = audirvanaTrackInfo.Url
			}
			// todo插入听歌流水
		case <-stop:
			fmt.Println("check playing track exit")
			return
		}
	}

}

func CompareTrackData() {

}
