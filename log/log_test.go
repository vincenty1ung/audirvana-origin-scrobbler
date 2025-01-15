package log

import (
	"testing"
	"time"

	"go.uber.org/zap"
)

func TestLogInit(t *testing.T) {
	init := LogInit("./logs", "info", make(<-chan struct{}))
	init.Info("haha", zap.String("time", time.Now().Format("2006-01-02 15:04:05")))
}
