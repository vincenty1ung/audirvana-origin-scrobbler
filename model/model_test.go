package model

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/vincenty1ung/lastfm-scrobbler/config"
	"github.com/vincenty1ung/lastfm-scrobbler/log"
	"github.com/vincenty1ung/lastfm-scrobbler/telemetry"
)

func setupTestDB(t *testing.T) *gorm.DB {
	// Create custom logger
	logger := log.LogInit("./logs", "debug", make(<-chan struct{}))

	customLogger := NewCustomLogger(logger)
	// Create a temporary in-memory database for testing
	db, err := gorm.Open(
		sqlite.Open(":memory:"), &gorm.Config{
			Logger: customLogger, // Disable logging for tests
		},
	)
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}

	// Auto migrate the schemas
	err = db.AutoMigrate(&TrackPlayRecord{}, &TrackPlayCount{})
	if err != nil {
		t.Fatalf("Failed to auto migrate: %v", err)
	}

	return db
}

func TestTrackPlayRecordCRUD(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	GlobalDB = db

	ctx := context.Background()
	if err := telemetry.Init(
		config.TelemetryConfig{
			Name:           "test",
			Endpoint:       "",
			Sampler:        0,
			Batcher:        "",
			OtlpHeaders:    nil,
			OtlpHttpPath:   "",
			OtlpHttpSecure: false,
			Disabled:       false,
		},
	); err != nil {
		log.Error(ctx, "failed to initialize telemetry")
		return
	}
	defer func(ctx context.Context) {
		err := telemetry.Shutdown(ctx)
		if err != nil {
		}
	}(ctx)

	ctx, span := telemetry.StartSpan(ctx, "TestTrackPlayRecordCRUD")
	defer span.End()

	// Test InsertTrackPlayRecord
	record := &TrackPlayRecord{
		Artist:        "Test Artist",
		AlbumArtist:   "Test Album Artist",
		Track:         "Test Track",
		Album:         "Test Album",
		Duration:      180,
		PlayTime:      time.Now(),
		Scrobbled:     true,
		MusicBrainzID: "test-mbid",
		TrackNumber:   1,
		Source:        "Audirvana",
	}
	err := InsertTrackPlayRecord(ctx, record)
	log.Warn(ctx, "adding record", zap.String("TraceIDFromContext", telemetry.TraceIDFromContext(ctx)))
	assert.NoError(t, err)
	assert.NotZero(t, record.ID)

	// Test GetUnscrobbledRecords
	// First, insert a record that is not scrobbled
	record2 := &TrackPlayRecord{
		Artist:    "Test Artist 2",
		Track:     "Test Track 2",
		Album:     "Test Album 2",
		Duration:  240,
		PlayTime:  time.Now(),
		Scrobbled: false,
		Source:    "Roon",
	}

	err = InsertTrackPlayRecord(ctx, record2)
	assert.NoError(t, err)

	// Get unscrobbled records
	records, err := GetUnscrobbledRecords(ctx, 10)
	assert.NoError(t, err)
	assert.Len(t, records, 1)
	assert.Equal(t, "Test Artist 2", records[0].Artist)
	assert.Equal(t, "Test Track 2", records[0].Track)
	assert.Equal(t, "Roon", records[0].Source)

	// Test UpdateScrobbledStatus
	err = UpdateScrobbledStatus(ctx, records[0].ID, true)
	assert.NoError(t, err)

	// Verify the record is now scrobbled
	records, err = GetUnscrobbledRecords(ctx, 10)
	assert.NoError(t, err)
	assert.Len(t, records, 0)
}

func TestTrackPlayCountCRUD(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	GlobalDB = db

	ctx := context.Background()
	if err := telemetry.Init(
		config.TelemetryConfig{
			Name:           "test",
			Endpoint:       "",
			Sampler:        0,
			Batcher:        "",
			OtlpHeaders:    nil,
			OtlpHttpPath:   "",
			OtlpHttpSecure: false,
			Disabled:       false,
		},
	); err != nil {
		log.Error(ctx, "failed to initialize telemetry")
		return
	}
	defer func(ctx context.Context) {
		err := telemetry.Shutdown(ctx)
		if err != nil {
		}
	}(ctx)

	ctx, span := telemetry.StartSpan(ctx, "TestTrackPlayCountCRUD")
	defer span.End()

	// Test IncrementTrackPlayCount
	artist := "Test Artist"
	album := "Test Album"
	track := "Test Track"

	// Increment play count for the first time
	err := IncrementTrackPlayCount(ctx, artist, album, track)
	assert.NoError(t, err)

	// Check the play count
	record, err := GetTrackPlayCount(ctx, artist, album, track)
	assert.NoError(t, err)
	assert.Equal(t, 1, record.PlayCount)
	assert.Equal(t, artist, record.Artist)
	assert.Equal(t, album, record.Album)
	assert.Equal(t, track, record.Track)

	// Increment play count again
	err = IncrementTrackPlayCount(ctx, artist, album, track)
	assert.NoError(t, err)

	// Check the play count is now 2
	record, err = GetTrackPlayCount(ctx, artist, album, track)
	assert.NoError(t, err)
	assert.Equal(t, 2, record.PlayCount)

	// Test GetTrackPlayCounts
	// Add another track
	err = IncrementTrackPlayCount(ctx, "Another Artist", "Another Album", "Another Track")
	assert.NoError(t, err)

	// Get play counts
	records, err := GetTrackPlayCounts(ctx, 10, 0)
	assert.NoError(t, err)
	assert.Len(t, records, 2)

	// The first record should have the higher play count
	if records[0].PlayCount < records[1].PlayCount {
		records[0], records[1] = records[1], records[0]
	}

	assert.Equal(t, artist, records[0].Artist)
	assert.Equal(t, 2, records[0].PlayCount)
	assert.Equal(t, "Another Artist", records[1].Artist)
	assert.Equal(t, 1, records[1].PlayCount)
}
