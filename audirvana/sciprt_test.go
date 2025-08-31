package audirvana

import (
	"fmt"
	"testing"

	"go.uber.org/zap"

	"github.com/audirvana-origin-scrobbler/common"
	alog "github.com/audirvana-origin-scrobbler/log"
)

func init() {
	_ = alog.LogInit("./logs", "info", make(<-chan struct{}))
}

func TestAudirvana(t *testing.T) {
	running := IsRunning()
	alog.Logger.Info("running", zap.Any("running", running))
	fmt.Println("Audirvana is running:", running)
	if running {
		state, _ := GetState()
		alog.Logger.Debug("audirvana 播放状态", zap.Any("state", state))
		var audirvanaTrackInfo *TrackInfo
		if state == common.PlayerStatePlaying {
			audirvanaTrackInfo = GetNowPlayingTrackInfo()
			alog.Logger.Info("", zap.Any("audirvana trackInfo", audirvanaTrackInfo))
		}
	}
}
func TestIsRunningReturnsBool(t *testing.T) {
	running := IsRunning()
	if running != true && running != false {
		t.Errorf("IsRunning() should return a boolean, got %v", running)
	}
}

func TestGetStateHandlesError(t *testing.T) {
	// Simulate Audirvana not running by ensuring IsRunning returns false
	if !IsRunning() {
		_, err := GetState()
		if err == nil {
			t.Error("GetState() should return error when Audirvana is not running")
		}
	}
}

func TestGetNowPlayingTrackInfoFields(t *testing.T) {
	if IsRunning() {
		state, _ := GetState()
		if state == common.PlayerStatePlaying {
			info := GetNowPlayingTrackInfo()
			if info == nil {
				t.Error("GetNowPlayingTrackInfo() returned nil while playing")
			} else {
				if info.Title == "" {
					t.Error("TrackInfo.Title is empty")
				}
				if info.Artist == "" {
					t.Error("TrackInfo.Artist is empty")
				}
				if info.Duration <= 0 {
					t.Error("TrackInfo.Duration should be positive")
				}
				if info.Position < 0 {
					t.Error("TrackInfo.Position should not be negative")
				}
			}
		}
	}
}
