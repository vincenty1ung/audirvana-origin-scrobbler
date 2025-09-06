package analysis

import (
	"context"
	"fmt"
	"github.com/vincenty1ung/lastfm-scrobbler/model"
	"github.com/vincenty1ung/lastfm-scrobbler/log"
	"go.uber.org/zap"
	"time"
)

// GenerateMusicPreferenceReport 生成音乐偏好分析报告
func GenerateMusicPreferenceReport(ctx context.Context) error {
	logger := log.Logger
	
	// 获取播放统计总数
	totalTracks, err := model.GetTrackCounts(ctx)
	if err != nil {
		logger.Error("Failed to get track counts", zap.Error(err))
		return err
	}
	
	// 获取播放次数最多的曲目
	topTracks, err := model.GetTrackPlayCounts(ctx, 10, 0)
	if err != nil {
		logger.Error("Failed to get top tracks", zap.Error(err))
		return err
	}
	
	// 获取最近播放的曲目
	recentRecords, err := getRecentPlayRecords(ctx, 20)
	if err != nil {
		logger.Error("Failed to get recent play records", zap.Error(err))
		return err
	}
	
	// 打印报告
	fmt.Println("=== 音乐偏好分析报告 ===")
	fmt.Printf("总曲目数: %d\n", totalTracks)
	fmt.Println("\n播放次数最多的曲目:")
	for i, track := range topTracks {
		fmt.Printf("%d. %s - %s - %s (播放次数: %d)\n", i+1, track.Artist, track.Album, track.Track, track.PlayCount)
	}
	
	fmt.Println("\n最近播放的曲目:")
	for i, record := range recentRecords {
		fmt.Printf("%d. %s - %s - %s (%s)\n", i+1, record.Artist, record.Album, record.Track, record.PlayTime.Format("2006-01-02 15:04:05"))
	}
	
	return nil
}

// getRecentPlayRecords 获取最近播放的记录


// ScheduleReport 定时生成报告
func ScheduleReport(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			if err := GenerateMusicPreferenceReport(ctx); err != nil {
				log.Logger.Error("Failed to generate music preference report", zap.Error(err))
			}
		case <-ctx.Done():
			return
		}
	}
}