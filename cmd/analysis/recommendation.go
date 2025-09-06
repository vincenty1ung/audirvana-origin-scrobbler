package analysis

import (
	"context"
	"fmt"
	"math"
	"github.com/vincenty1ung/lastfm-scrobbler/model"
	"github.com/vincenty1ung/lastfm-scrobbler/log"
	"go.uber.org/zap"
)

// MusicRecommendation 音乐推荐结构
type MusicRecommendation struct {
	Artist string
	Album  string
	Track  string
	Score  float64
}

// GenerateMusicRecommendations 基于历史播放记录生成音乐推荐
func GenerateMusicRecommendations(ctx context.Context, limit int) ([]MusicRecommendation, error) {
	
	
	// 获取所有播放统计记录
	allTracks, err := getAllTrackPlayCounts(ctx)
	if err != nil {
		log.Error(ctx, "Failed to get track play counts", zap.Error(err))
		return nil, err
	}
	
	// 获取最近播放的记录
	recentRecords, err := getRecentPlayRecords(ctx, 50)
	if err != nil {
		log.Error(ctx, "Failed to get recent play records", zap.Error(err))
		return nil, err
	}
	
	// 基于播放统计和最近播放记录生成推荐
	recommendations := calculateRecommendations(allTracks, recentRecords, limit)
	
	return recommendations, nil
}

// getAllTrackPlayCounts 获取所有播放统计记录
func getAllTrackPlayCounts(ctx context.Context) ([]*model.TrackPlayCount, error) {
	return model.GetAllTrackPlayCounts(ctx)
}

// getRecentPlayRecords 获取最近播放的记录
func getRecentPlayRecords(ctx context.Context, limit int) ([]*model.TrackPlayRecord, error) {
	return model.GetRecentPlayRecords(ctx, limit)
}

// calculateRecommendations 计算推荐分数并生成推荐列表
func calculateRecommendations(allTracks []*model.TrackPlayCount, recentRecords []*model.TrackPlayRecord, limit int) []MusicRecommendation {
	// 创建艺术家和专辑的播放频率映射
	artistFrequency := make(map[string]int)
	albumFrequency := make(map[string]int)
	
	// 统计最近播放记录中的艺术家和专辑频率
	for _, record := range recentRecords {
		artistFrequency[record.Artist]++
		albumFrequency[record.Album]++
	}
	
	// 计算推荐分数
	var recommendations []MusicRecommendation
	for _, track := range allTracks {
		score := calculateScore(track, artistFrequency, albumFrequency)
		recommendations = append(recommendations, MusicRecommendation{
			Artist: track.Artist,
			Album:  track.Album,
			Track:  track.Track,
			Score:  score,
		})
	}
	
	// 按分数排序并返回前limit个推荐
	sortRecommendations(recommendations)
	
	if len(recommendations) > limit {
		return recommendations[:limit]
	}
	
	return recommendations
}

// calculateScore 计算单个曲目的推荐分数
func calculateScore(track *model.TrackPlayCount, artistFrequency, albumFrequency map[string]int) float64 {
	// 基础分数基于播放次数
	baseScore := float64(track.PlayCount)
	
	// 艺术家频率加权
	artistWeight := 1.0
	if freq, ok := artistFrequency[track.Artist]; ok {
		artistWeight = 1.0 + math.Log(float64(freq)+1)
	}
	
	// 专辑频率加权
	albumWeight := 1.0
	if freq, ok := albumFrequency[track.Album]; ok {
		albumWeight = 1.0 + math.Log(float64(freq)+1)
	}
	
	// 计算最终分数
	finalScore := baseScore * artistWeight * albumWeight
	
	return finalScore
}

// sortRecommendations 按分数降序排序推荐列表
func sortRecommendations(recommendations []MusicRecommendation) {
	for i := 0; i < len(recommendations)-1; i++ {
		for j := i + 1; j < len(recommendations); j++ {
			if recommendations[i].Score < recommendations[j].Score {
				recommendations[i], recommendations[j] = recommendations[j], recommendations[i]
			}
		}
	}
}

// PrintRecommendations 打印音乐推荐
func PrintRecommendations(recommendations []MusicRecommendation) {
	fmt.Println("=== 音乐推荐 ===")
	for i, rec := range recommendations {
		fmt.Printf("%d. %s - %s - %s (推荐分数: %.2f)\n", i+1, rec.Artist, rec.Album, rec.Track, rec.Score)
	}
}

// GetArtistRecommendations 获取特定艺术家的推荐曲目
func GetArtistRecommendations(ctx context.Context, artist string, limit int) ([]MusicRecommendation, error) {
	
	
	// 获取该艺术家的所有曲目
	tracks, err := getTracksByArtist(ctx, artist)
	if err != nil {
		log.Error(ctx, "Failed to get tracks by artist", zap.String("artist", artist), zap.Error(err))
		return nil, err
	}
	
	// 计算推荐分数
	var recommendations []MusicRecommendation
	for _, track := range tracks {
		recommendations = append(recommendations, MusicRecommendation{
			Artist: track.Artist,
			Album:  track.Album,
			Track:  track.Track,
			Score:  float64(track.PlayCount),
		})
	}
	
	// 按分数排序
	sortRecommendations(recommendations)
	
	if len(recommendations) > limit {
		return recommendations[:limit], nil
	}
	
	return recommendations, nil
}

// getTracksByArtist 获取特定艺术家的所有曲目
func getTracksByArtist(ctx context.Context, artist string) ([]*model.TrackPlayCount, error) {
	return model.GetTracksByArtist(ctx, artist)
}