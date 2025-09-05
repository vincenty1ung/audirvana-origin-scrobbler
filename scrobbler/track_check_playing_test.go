package scrobbler

import (
	"testing"
	"time"

	"github.com/vincenty1ung/lastfm-scrobbler/audirvana"
	"github.com/vincenty1ung/lastfm-scrobbler/common"
)

// mockLastfm a mock implementation of the Last.fm API for testing.
type mockLastfm struct {
	nowPlayingCalled bool
	scrobbleCalled   bool
	lastTrack        string
}

func (m *mockLastfm) UpdateNowPlaying(req *TrackUpdateNowPlayingReq) error {
	m.nowPlayingCalled = true
	m.lastTrack = req.Track
	return nil
}

func (m *mockLastfm) Scrobble(req *PushTrackScrobbleReq) error {
	m.scrobbleCalled = true
	m.lastTrack = req.Track
	return nil
}

// processAudirvanaState is a refactored, testable version of the core logic.
// NOTE: This function is a simplified representation of the logic in AudirvanaCheckPlayingTrack for testing purposes.
// It's designed to be testable in isolation.
func processAudirvanaState(
	lastfm *mockLastfm,
	trackInfo *audirvana.TrackInfo,
	playerState common.PlayerState,
	previousTrack *string,
	scrobbledTracks map[string]bool,
) {
	if playerState != common.PlayerStatePlaying || trackInfo == nil {
		return
	}

	currentTrackKey := trackInfo.Url + trackInfo.Title
	if currentTrackKey != *previousTrack {
		// New track logic
		*previousTrack = currentTrackKey
		// Clear previous track from scrobbled map
		// In real code, this would be more complex
		for k := range scrobbledTracks {
			delete(scrobbledTracks, k)
		}

		lastfm.UpdateNowPlaying(&TrackUpdateNowPlayingReq{
			Artist: trackInfo.Artist,
			Track:  trackInfo.Title,
			Album:  trackInfo.Album,
		})
	}

	// Scrobble logic
	hasBeenScrobbled := scrobbledTracks[currentTrackKey]
	if !hasBeenScrobbled && (trackInfo.Position/float64(trackInfo.Duration)) > percentScrobble {
		lastfm.Scrobble(&PushTrackScrobbleReq{
			Artist:    trackInfo.Artist,
			Track:     trackInfo.Title,
			Album:     trackInfo.Album,
			Timestamp: time.Now().Unix(),
		})
		scrobbledTracks[currentTrackKey] = true
	}
}

func TestProcessAudirvanaState(t *testing.T) {
	// Test cases
	tests := []struct {
		name                   string
		playerState            common.PlayerState
		trackInfo              *audirvana.TrackInfo
		initialPreviousTrack   string
		initialScrobbledTracks map[string]bool
		expectNowPlaying       bool
		expectScrobble         bool
		expectedTrack          string
	}{
		{
			name:                   "Player is stopped",
			playerState:            common.PlayerStateStopped,
			trackInfo:              nil,
			initialPreviousTrack:   "",
			initialScrobbledTracks: make(map[string]bool),
			expectNowPlaying:       false,
			expectScrobble:         false,
		},
		{
			name:        "New track starts playing",
			playerState: common.PlayerStatePlaying,
			trackInfo: &audirvana.TrackInfo{
				TrackBase: audirvana.TrackBase{
					Title:    "New Song",
					Artist:   "Artist",
					Album:    "Album",
					Url:      "file://newsong",
					Duration: 200,
					Position: 10,
				},
			},
			initialPreviousTrack:   "",
			initialScrobbledTracks: make(map[string]bool),
			expectNowPlaying:       true,
			expectScrobble:         false,
			expectedTrack:          "New Song",
		},
		{
			name:        "Existing track continues, not ready to scrobble",
			playerState: common.PlayerStatePlaying,
			trackInfo: &audirvana.TrackInfo{
				TrackBase: audirvana.TrackBase{
					Title:    "Existing Song",
					Artist:   "Artist",
					Album:    "Album",
					Url:      "file://existingsong",
					Duration: 200,
					Position: 50, // Less than 55%
				},
			},
			initialPreviousTrack:   "file://existingsongExisting Song",
			initialScrobbledTracks: make(map[string]bool),
			expectNowPlaying:       false,
			expectScrobble:         false,
		},
		{
			name:        "Existing track is ready to scrobble",
			playerState: common.PlayerStatePlaying,
			trackInfo: &audirvana.TrackInfo{
				TrackBase: audirvana.TrackBase{
					Title:    "Scrobble Song",
					Artist:   "Artist",
					Album:    "Album",
					Url:      "file://scrobblesong",
					Duration: 200,
					Position: 120, // More than 55%
				},
			},
			initialPreviousTrack:   "file://scrobblesongScrobble Song",
			initialScrobbledTracks: make(map[string]bool),
			expectNowPlaying:       false,
			expectScrobble:         true,
			expectedTrack:          "Scrobble Song",
		},
		{
			name:        "Track ready to scrobble but already scrobbled",
			playerState: common.PlayerStatePlaying,
			trackInfo: &audirvana.TrackInfo{
				TrackBase: audirvana.TrackBase{
					Title:    "Scrobble Song",
					Artist:   "Artist",
					Album:    "Album",
					Url:      "file://scrobblesong",
					Duration: 200,
					Position: 130, // More than 55%
				},
			},
			initialPreviousTrack: "file://scrobblesongScrobble Song",
			initialScrobbledTracks: map[string]bool{
				"file://scrobblesongScrobble Song": true,
			},
			expectNowPlaying: false,
			expectScrobble:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockAPI := &mockLastfm{}
			previousTrack := tt.initialPreviousTrack
			scrobbledTracks := tt.initialScrobbledTracks

			// Execute
			processAudirvanaState(mockAPI, tt.trackInfo, tt.playerState, &previousTrack, scrobbledTracks)

			// Assert
			if mockAPI.nowPlayingCalled != tt.expectNowPlaying {
				t.Errorf("Expected nowPlayingCalled to be %v, but got %v", tt.expectNowPlaying, mockAPI.nowPlayingCalled)
			}

			if mockAPI.scrobbleCalled != tt.expectScrobble {
				t.Errorf("Expected scrobbleCalled to be %v, but got %v", tt.expectScrobble, mockAPI.scrobbleCalled)
			}

			if (tt.expectNowPlaying || tt.expectScrobble) && mockAPI.lastTrack != tt.expectedTrack {
				t.Errorf("Expected track to be '%s', but got '%s'", tt.expectedTrack, mockAPI.lastTrack)
			}
		})
	}
}
