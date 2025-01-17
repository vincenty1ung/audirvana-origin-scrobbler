package scrobbler

import (
	"fmt"
	"log"
	"sync/atomic"
	"time"

	"github.com/vincenty1ung/yeung-go-study/lru"
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
var pushCount = atomic.Uint32{} // 多渠道上报
var isLong bool

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
			if tmpCount > checkCount && !isLong { // 检查100次依旧没有播放检查轮训放大到60秒
				timer.Reset(time.Second * longSleep)
				isLong = true
				alog.Logger.Info(
					"检查100次依旧没有播放检查轮训放大到60秒", zap.Uint32("共计上传歌曲标记", pushCount.Load()),
				)
			}
			if isLong {
				alog.Logger.Info(
					"60秒检查", zap.Uint32("共计上传歌曲标记", pushCount.Load()),
				)
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
						isLong = false
						timer.Reset(time.Second * defaultSleep)
					}
					tmpCount = 0
					audirvanaTrackInfo = audirvana.GetNowPlayingTrackInfo()
				}
			}
			if audirvanaTrackInfo != nil {
				tmpTrack := audirvanaTrackInfo.Url + audirvanaTrackInfo.Title // 防止cue文件出现问题
				currentTrack = tmpTrack
				position := audirvanaTrackInfo.Position
				duration := audirvanaTrackInfo.Duration
				if position/float64(duration) > percentScrobble && !maped[currentTrack] {
					// 标记听歌完成
					pushTrackScrobbleReq := &PushTrackScrobbleReq{
						Artist:      audirvanaTrackInfo.Artist,
						AlbumArtist: audirvanaTrackInfo.Artist,
						Track:       audirvanaTrackInfo.Title,
						Album:       audirvanaTrackInfo.Album,
						Duration:    audirvanaTrackInfo.Duration,
						Timestamp:   now.UTC().Unix(),
					}
					// 说明在听歌存在有效数据的
					if exiftoolInfo := FindExiftoolInfoCache(audirvanaTrackInfo.Url); exiftoolInfo != nil {
						pushTrackScrobbleReq.TrackNumber = exiftoolInfo.GetTrackNumber()
						pushTrackScrobbleReq.MusicBrainzTrackID = exiftoolInfo.GetMusicBrainzTrackId()
						if artist := exiftoolInfo.GetArtist(); len(artist) != 0 {
							pushTrackScrobbleReq.Artist = artist
						}
						if albumartist := exiftoolInfo.GetAlbumartist(); len(albumartist) != 0 {
							pushTrackScrobbleReq.AlbumArtist = albumartist
						}
					}
					_, err := PushTrackScrobble(pushTrackScrobbleReq)
					if err != nil {
						log.Fatal(err)
					}
					maped[currentTrack] = true
					pushCount.Add(1)
					alog.Logger.Info("标记听歌完成", zap.String("track", pushTrackScrobbleReq.Track))
				}
				// 上传听歌ing
				if currentTrack != previousTrack {
					// 产生新歌曲
					delete(maped, previousTrack)
					now = time.Now()
					playingReq := TrackUpdateNowPlayingReq{
						Artist:      audirvanaTrackInfo.Artist,
						AlbumArtist: audirvanaTrackInfo.Artist,
						Track:       audirvanaTrackInfo.Title,
						Album:       audirvanaTrackInfo.Album,
						Duration:    audirvanaTrackInfo.Duration,
					}
					// 说明在听歌存在有效数据的
					if exiftoolInfo := FindExiftoolInfoCache(audirvanaTrackInfo.Url); exiftoolInfo != nil {
						playingReq.TrackNumber = exiftoolInfo.GetTrackNumber()
						playingReq.MusicBrainzTrackID = exiftoolInfo.GetMusicBrainzTrackId()
						if artist := exiftoolInfo.GetArtist(); len(artist) != 0 {
							playingReq.Artist = artist
						}
						if albumartist := exiftoolInfo.GetAlbumartist(); len(albumartist) != 0 {
							playingReq.AlbumArtist = albumartist
						}
					}
					alog.Logger.Info("NowPlayingTrackInfo", zap.Any("audirvanaTrackInfo", audirvanaTrackInfo))
					err := TrackUpdateNowPlaying(&playingReq)
					if err != nil {
						alog.Logger.Warn("TrackUpdateNowPlaying", zap.Error(err))
						continue
					}
				}
				previousTrack = tmpTrack // 防止cue文件出现问题
				// todo插入听歌流水
			}
		case <-stop:
			fmt.Println("check playing track exit")
			return
		}
	}

}

var lruCache = lru.Constructor[string](100)

func FindExiftoolInfoCache(key string) *exec.ExiftoolInfo {
	var (
		exiftoolInfo *exec.ExiftoolInfo
		err          error
	)

	if cacheExiftoolInfo := lruCache.Get(key); cacheExiftoolInfo != nil {
		exiftoolInfo = cacheExiftoolInfo.(*exec.ExiftoolInfo)
	} else {
		exiftoolInfo, err = exec.ExiftoolHandle(key)
		if err != nil {
			alog.Logger.Warn("exec ExiftoolHandle", zap.Error(err))
			return exiftoolInfo
		}
		if exiftoolInfo != nil {
			lruCache.Put(key, exiftoolInfo)
		}
	}
	return exiftoolInfo
}
