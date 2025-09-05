package exec

import (
	"context"

	"github.com/vincenty1ung/yeung-go-study/lru"
	"go.uber.org/zap"

	alog "github.com/vincenty1ung/lastfm-scrobbler/log"

	"github.com/vincenty1ung/lastfm-scrobbler/common"
)

var lruCache = lru.Constructor[string](200)

func FindMataDataHandleCache(ctx context.Context, key string) MataDataHandle {
	var (
		mataDataHandle MataDataHandle
		err            error
	)

	if exiftoolInfo := lruCache.Get(key); exiftoolInfo != nil {
		mataDataHandle = exiftoolInfo.(MataDataHandle)
	} else {
		if ok, path, _ := IsValidPath(ctx, key); ok {
			if GetFilePathExt(path) == common.FileExtWav1 || GetFilePathExt(path) == common.FileExtWav2 {
				mataDataHandle, err = BuildWavInfoHandle(path)
				if err != nil {
					alog.Warn(ctx, "exec BuildExiftoolHandle", zap.Error(err))
					return mataDataHandle
				}
				if mataDataHandle != nil {
					lruCache.Put(key, mataDataHandle)
				}
			} else {
				mataDataHandle, err = BuildExiftoolHandle(path)
				if err != nil {
					alog.Warn(ctx, "exec BuildExiftoolHandle", zap.Error(err))
					return mataDataHandle
				}
				if mataDataHandle != nil {
					lruCache.Put(key, mataDataHandle)
				}
			}
		}
	}
	return mataDataHandle
}
