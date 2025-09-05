package model

import (
	"context"
	"time"
)

type TrackPlayRecord struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	Artist        string    `gorm:"index" json:"artist"`
	AlbumArtist   string    `json:"album_artist"`
	Track         string    `json:"track"`
	Album         string    `json:"album"`
	Duration      int64     `json:"duration"`
	PlayTime      time.Time `json:"play_time"`
	Scrobbled     bool      `gorm:"index" json:"scrobbled"` // 是否已同步到Last.fm
	MusicBrainzID string    `json:"musicbrainz_id"`
	TrackNumber   int64     `json:"track_number"`
	Source        string    `gorm:"index" json:"source"` // 数据来源：Audirvana 或 Roon
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func InsertTrackPlayRecord(ctx context.Context, record *TrackPlayRecord) error {
	return GetDB().WithContext(ctx).Create(record).Error
}

func UpdateScrobbledStatus(ctx context.Context, id uint, scrobbled bool) error {
	return GetDB().WithContext(ctx).Where("id = ?", id).Update("scrobbled", scrobbled).Error
}

func GetUnscrobbledRecords(ctx context.Context, limit int) ([]*TrackPlayRecord, error) {
	var trackPlayRecords []*TrackPlayRecord
	err := GetDB().WithContext(ctx).Where(
		"scrobbled = ?", false,
	).Order("play_time ASC").Limit(limit).Find(&trackPlayRecords).Error
	if err != nil {
		return nil, err
	}
	return trackPlayRecords, nil
}
