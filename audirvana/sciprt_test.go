package audirvana

import (
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
