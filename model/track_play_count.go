package model

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

type TrackPlayCount struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Artist    string    `gorm:"index;uniqueIndex:idx_track_album_artist" json:"artist"`
	Album     string    `gorm:"index;uniqueIndex:idx_track_album_artist" json:"album"`
	Track     string    `gorm:"index;uniqueIndex:idx_track_album_artist" json:"track"`
	PlayCount int       `json:"play_count"`
	Version   int       `gorm:"default:1" json:"version"` // 乐观锁版本号
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (TrackPlayCount) TableName() string {
	return "track_play_counts"
}

func IncrementTrackPlayCount(ctx context.Context, artist, album, track string) error {
	// 使用乐观锁机制更新播放次数
	for {
		var record TrackPlayCount
		err := GetDB().WithContext(ctx).Where(
			"artist = ? AND album = ? AND track = ?", artist, album, track,
		).First(&record).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// Create new record
				record = TrackPlayCount{
					Artist:    artist,
					Album:     album,
					Track:     track,
					PlayCount: 1,
				}
				err = GetDB().WithContext(ctx).Create(&record).Error
				if err != nil && !errors.Is(err, gorm.ErrDuplicatedKey) {
					return err
				}
				// 如果出现重复键错误，说明其他goroutine已经创建了记录，继续循环处理
				if errors.Is(err, gorm.ErrDuplicatedKey) {
					continue
				}
				return nil
			}
			return err
		}

		// Update existing record with optimistic locking
		updatedRecord := TrackPlayCount{
			PlayCount: record.PlayCount + 1,
			Version:   record.Version + 1,
		}

		result := GetDB().WithContext(ctx).Where(
			"artist = ? AND album = ? AND track = ? AND version = ?",
			artist, album, track, record.Version,
		).Updates(&updatedRecord)

		if result.Error != nil {
			return result.Error
		}

		// 如果更新成功，跳出循环
		if result.RowsAffected > 0 {
			break
		}
		// 如果更新失败（版本号不匹配），继续循环重试
	}

	return nil
}

func GetTrackPlayCounts(ctx context.Context, limit, offset int) ([]*TrackPlayCount, error) {
	var records []*TrackPlayCount
	err := GetDB().WithContext(ctx).Order("play_count DESC").Limit(limit).Offset(offset).Find(&records).Error
	if err != nil {
		return nil, err
	}
	return records, nil
}
func GetTrackCounts(ctx context.Context) (int64, error) {
	var count int64
	err := GetDB().WithContext(ctx).Model(&TrackPlayCount{}).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func GetTrackPlayCount(ctx context.Context, artist, album, track string) (*TrackPlayCount, error) {
	var record TrackPlayCount
	err := GetDB().WithContext(ctx).Where(
		"artist = ? AND album = ? AND track = ?", artist, album, track,
	).First(&record).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}

// GetAllTrackPlayCounts 获取所有播放统计记录
func GetAllTrackPlayCounts(ctx context.Context) ([]*TrackPlayCount, error) {
	var allTracks []*TrackPlayCount
	pageSize := 100
	offset := 0
	
	for {
		var tracks []*TrackPlayCount
		err := GetDB().WithContext(ctx).Order("play_count DESC").Limit(pageSize).Offset(offset).Find(&tracks).Error
		if err != nil {
			return nil, err
		}
		
		allTracks = append(allTracks, tracks...)
		
		// 如果返回的记录数少于pageSize，说明已经获取完所有记录
		if len(tracks) < pageSize {
			break
		}
		
		offset += pageSize
	}
	
	return allTracks, nil
}

// GetTracksByArtist 获取特定艺术家的所有曲目
func GetTracksByArtist(ctx context.Context, artist string) ([]*TrackPlayCount, error) {
	var tracks []*TrackPlayCount
	err := GetDB().WithContext(ctx).Where("artist LIKE ?", "%"+artist+"%").Find(&tracks).Error
	if err != nil {
		return nil, err
	}
	return tracks, nil
}
