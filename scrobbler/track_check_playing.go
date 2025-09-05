package scrobbler

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"go.uber.org/zap"

	"github.com/audirvana-origin-scrobbler/audirvana"
	"github.com/audirvana-origin-scrobbler/common"
	"github.com/audirvana-origin-scrobbler/exec"
	alog "github.com/audirvana-origin-scrobbler/log"
	"github.com/audirvana-origin-scrobbler/model"
	"github.com/audirvana-origin-scrobbler/telemetry"
)

const (
	percentScrobble = 0.55
	defaultSleep    = 3
	longSleep       = 60 // 休眠间隔六十秒
	checkCount      = 100
)

var maped = make(map[string]bool)
var maped2 = make(map[string]bool)
var pushCount = atomic.Uint32{} // 多渠道上报
var isLong bool
var isLong2 bool

func AudirvanaCheckPlayingTrack(ctx context.Context, stop <-chan struct{}) {
	timer := time.NewTicker(time.Second * defaultSleep)
	var (
		now           time.Time
		currentTrack  = ""
		previousTrack = ""
		tmpCount      = 0
	)
	counts, err := model.GetTrackCounts(ctx)
	if err != nil {
		panic(err)
	}
	pushCount.Store(uint32(counts))
	for {
		select {
		case <-timer.C:
			h := func(ctx context.Context) {
				// Start a new span for this check cycle
				checkCtx, span := telemetry.StartSpanForTracerName(
					ctx, _TracerName, "audirvanaCheckPlayingTrack",
				)
				defer span.End()

				alog.Debug(checkCtx, "AudirvanaCheckPlayingTrack Checking playing track..."+time.Now().String())

				// End the span at the end of this cycle
				// span.End()
				tmpCount++
				if tmpCount > checkCount && !isLong { // 检查100次依旧没有播放检查轮训放大到60秒
					timer.Reset(time.Second * longSleep)
					isLong = true
					alog.Info(
						checkCtx, "检查100次依旧没有播放检查轮训放大到60秒",
						zap.Uint32("共计上传歌曲标记", pushCount.Load()),
					)
				}
				if isLong {
					alog.Info(checkCtx, "60秒检查", zap.Uint32("共计上传歌曲标记", pushCount.Load()))
				}
				running := audirvana.IsRunning(checkCtx)
				alog.Debug(checkCtx, "程序运行是否运行", zap.Bool("running", running))
				var audirvanaTrackInfo *audirvana.TrackInfo
				if running {
					audirvanaTrackInfo = nil
					state, _ := audirvana.GetState(checkCtx)
					alog.Debug(checkCtx, "audirvana 播放状态", zap.Any("state", state))
					if state == common.PlayerStatePlaying {
						if tmpCount > checkCount {
							isLong = false
							timer.Reset(time.Second * defaultSleep)
						}
						tmpCount = 0
						audirvanaTrackInfo = audirvana.GetNowPlayingTrackInfo(checkCtx)
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
						if mataDataHandleCache := exec.FindMataDataHandleCache(
							checkCtx, audirvanaTrackInfo.Url,
						); mataDataHandleCache != nil {
							pushTrackScrobbleReq.TrackNumber = mataDataHandleCache.GetTrackNumber()
							pushTrackScrobbleReq.MusicBrainzTrackID = mataDataHandleCache.GetMusicBrainzTrackId()
							if artist := mataDataHandleCache.GetArtist(); len(artist) != 0 {
								pushTrackScrobbleReq.Artist = artist
							}
							if albumartist := mataDataHandleCache.GetAlbumartist(); len(albumartist) != 0 {
								pushTrackScrobbleReq.AlbumArtist = albumartist
							}
						}
						_, err := PushTrackScrobble(checkCtx, pushTrackScrobbleReq)
						if err != nil {
							alog.Warn(checkCtx, "TrackUpdateNowPlaying", zap.Error(err))
							return
						}
						// Save to database
						record := &model.TrackPlayRecord{
							Artist:        pushTrackScrobbleReq.Artist,
							AlbumArtist:   pushTrackScrobbleReq.AlbumArtist,
							Track:         pushTrackScrobbleReq.Track,
							Album:         pushTrackScrobbleReq.Album,
							Duration:      pushTrackScrobbleReq.Duration,
							PlayTime:      time.Unix(pushTrackScrobbleReq.Timestamp, 0),
							Scrobbled:     true,
							MusicBrainzID: pushTrackScrobbleReq.MusicBrainzTrackID,
							TrackNumber:   pushTrackScrobbleReq.TrackNumber,
							Source:        "Audirvana",
						}
						if err := model.InsertTrackPlayRecord(checkCtx, record); err != nil {
							alog.Warn(checkCtx, "Failed to insert track play record", zap.Error(err))
						}

						// Update track play count
						if err := model.IncrementTrackPlayCount(
							checkCtx, record.Artist, record.Album, record.Track,
						); err != nil {
							alog.Warn(checkCtx, "Failed to increment track play count", zap.Error(err))
						}

						maped[currentTrack] = true
						pushCount.Add(1)
						alog.Info(checkCtx, "标记听歌完成", zap.String("track", pushTrackScrobbleReq.Track))
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
						if mataDataHandleCache := exec.FindMataDataHandleCache(
							ctx, audirvanaTrackInfo.Url,
						); mataDataHandleCache != nil {
							playingReq.TrackNumber = mataDataHandleCache.GetTrackNumber()
							playingReq.MusicBrainzTrackID = mataDataHandleCache.GetMusicBrainzTrackId()
							if artist := mataDataHandleCache.GetArtist(); len(artist) != 0 {
								playingReq.Artist = artist
							}
							if albumartist := mataDataHandleCache.GetAlbumartist(); len(albumartist) != 0 {
								playingReq.AlbumArtist = albumartist
							}
						}
						alog.Info(checkCtx, "NowPlayingTrackInfo", zap.Any("audirvanaTrackInfo", audirvanaTrackInfo))
						err := TrackUpdateNowPlaying(checkCtx, &playingReq)
						if err != nil {
							alog.Warn(checkCtx, "TrackUpdateNowPlaying", zap.Error(err))
							return
						}
					}
					previousTrack = tmpTrack // 防止cue文件出现问题
				}
			}
			h(ctx)
		case <-stop:
			fmt.Println("check playing track exit")
			return
		}
	}
}
func RoonCheckPlayingTrack(ctx context.Context, stop <-chan struct{}) {
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
			h := func(ctx context.Context) {
				// Start a new span for this check cycle
				checkCtx, span := telemetry.StartSpanForTracerName(
					ctx, _TracerName, "roonCheckPlayingTrack",
				)
				defer span.End()

				alog.Debug(checkCtx, "RoonCheckPlayingTrack Checking playing track..."+time.Now().String())
				tmpCount++
				if tmpCount > checkCount && !isLong2 { // 检查100次依旧没有播放检查轮训放大到60秒
					timer.Reset(time.Second * longSleep)
					isLong2 = true
					alog.Info(
						checkCtx, "检查100次依旧没有播放检查轮训放大到60秒",
						zap.Uint32("共计上传歌曲标记", pushCount.Load()),
					)
				}
				if isLong2 {
					alog.Info(checkCtx, "60秒检查", zap.Uint32("共计上传歌曲标记", pushCount.Load()))
				}
				playing, err := exec.GetMRMediaNowPlaying()
				if err != nil {
					alog.Warn(checkCtx, "TrackUpdateNowPlaying", zap.Error(err))
					return
				}
				var roonTrackInfo *exec.MRMediaNowPlaying
				if playing.BundleIdentifier == exec.MRMediaNowPlayingAppRoon {
					roonTrackInfo = nil
					if playing.IsPlaying {
						if tmpCount > checkCount {
							isLong2 = false
							timer.Reset(time.Second * defaultSleep)
						}
						tmpCount = 0
						roonTrackInfo = playing
					}
				}
				if roonTrackInfo != nil {
					tmpTrack := roonTrackInfo.Title //
					currentTrack = tmpTrack
					position := roonTrackInfo.ElapsedTime
					duration := roonTrackInfo.Duration
					if position/float64(duration) > percentScrobble && !maped2[currentTrack] {
						// 标记听歌完成
						pushTrackScrobbleReq := &PushTrackScrobbleReq{
							Artist:      roonTrackInfo.Artist,
							AlbumArtist: roonTrackInfo.Artist,
							Track:       roonTrackInfo.Title,
							Album:       roonTrackInfo.Album,
							Duration:    int64(roonTrackInfo.Duration),
							Timestamp:   now.UTC().Unix(),
						}
						// 说明在听歌存在有效数据的
						_, err := PushTrackScrobble(checkCtx, pushTrackScrobbleReq)
						if err != nil {
							alog.Warn(checkCtx, "RoonCheckPlayingTrack TrackUpdateNowPlaying", zap.Error(err))
							return
						}
						// Save to database
						record := &model.TrackPlayRecord{
							Artist:        pushTrackScrobbleReq.Artist,
							AlbumArtist:   pushTrackScrobbleReq.AlbumArtist,
							Track:         pushTrackScrobbleReq.Track,
							Album:         pushTrackScrobbleReq.Album,
							Duration:      pushTrackScrobbleReq.Duration,
							PlayTime:      time.Unix(pushTrackScrobbleReq.Timestamp, 0),
							Scrobbled:     true,
							MusicBrainzID: pushTrackScrobbleReq.MusicBrainzTrackID,
							TrackNumber:   pushTrackScrobbleReq.TrackNumber,
							Source:        "Roon",
						}
						if err := model.InsertTrackPlayRecord(checkCtx, record); err != nil {
							alog.Warn(checkCtx, "Failed to insert track play record", zap.Error(err))
						}
						// Update track play count
						if err := model.IncrementTrackPlayCount(
							checkCtx, record.Artist, record.Album, record.Track,
						); err != nil {
							alog.Warn(checkCtx, "Failed to increment track play count", zap.Error(err))
						}

						maped2[currentTrack] = true
						pushCount.Add(1)
						alog.Info(
							checkCtx, "RoonCheckPlayingTrack 标记听歌完成",
							zap.String("track", pushTrackScrobbleReq.Track),
						)
					}
					// 上传听歌ing
					if currentTrack != previousTrack {
						// 产生新歌曲
						delete(maped2, previousTrack)
						now = time.Now()
						playingReq := TrackUpdateNowPlayingReq{
							Artist:      roonTrackInfo.Artist,
							AlbumArtist: roonTrackInfo.Artist,
							Track:       roonTrackInfo.Title,
							Album:       roonTrackInfo.Album,
							Duration:    int64(roonTrackInfo.Duration),
						}
						alog.Info(
							checkCtx, "RoonCheckPlayingTrack NowPlayingTrackInfo",
							zap.Any("roonTrackInfo", roonTrackInfo),
						)
						err := TrackUpdateNowPlaying(checkCtx, &playingReq)
						if err != nil {
							alog.Warn(ctx, "RoonCheckPlayingTrack TrackUpdateNowPlaying", zap.Error(err))
							return
						}
					}
					previousTrack = tmpTrack // 防止cue文件出现问题
				}
			}
			h(ctx)
		case <-stop:
			fmt.Println("RoonCheckPlayingTrack check playing track exit")
			return
		}
	}
}
