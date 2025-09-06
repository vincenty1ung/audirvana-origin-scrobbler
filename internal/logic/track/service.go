package track

import (
	"context"

	model2 "github.com/vincenty1ung/lastfm-scrobbler/internal/model"
)

// TrackService 定义曲目相关服务接口
type TrackService interface {
	GetTrackPlayCounts(ctx context.Context, limit, offset int) ([]*model2.TrackPlayCount, error)
	GetTrackPlayCount(ctx context.Context, artist, album, track string) (*model2.TrackPlayCount, error)
	InsertTrackPlayRecord(ctx context.Context, record *model2.TrackPlayRecord) error
	IncrementTrackPlayCount(ctx context.Context, artist, album, track string) error
}

// TrackServiceImpl 实现TrackService接口
type TrackServiceImpl struct{}

// NewTrackService 创建TrackService实例
func NewTrackService() TrackService {
	return &TrackServiceImpl{}
}

// GetTrackPlayCounts 获取曲目播放统计列表
func (s *TrackServiceImpl) GetTrackPlayCounts(ctx context.Context, limit, offset int) (
	[]*model2.TrackPlayCount, error,
) {
	return model2.GetTrackPlayCounts(ctx, limit, offset)
}

// GetTrackPlayCount 获取特定曲目的播放统计
func (s *TrackServiceImpl) GetTrackPlayCount(ctx context.Context, artist, album, track string) (
	*model2.TrackPlayCount, error,
) {
	return model2.GetTrackPlayCount(ctx, artist, album, track)
}

func (s *TrackServiceImpl) InsertTrackPlayRecord(ctx context.Context, record *model2.TrackPlayRecord) error {
	return model2.InsertTrackPlayRecord(ctx, record)
}

func (s *TrackServiceImpl) IncrementTrackPlayCount(ctx context.Context, artist, album, track string) error {
	return model2.IncrementTrackPlayCount(ctx, artist, album, track)
}
