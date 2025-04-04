package exec

import (
	"github.com/vincenty1ung/yeung-go-study/lru"
	"go.uber.org/zap"

	"github.com/audirvana-origin-scrobbler/common"
	alog "github.com/audirvana-origin-scrobbler/log"
)

var lruCache = lru.Constructor[string](200)

func FindMataDataHandleCache(key string) MataDataHandle {
	var (
		mataDataHandle MataDataHandle
		err            error
	)

	if exiftoolInfo := lruCache.Get(key); exiftoolInfo != nil {
		mataDataHandle = exiftoolInfo.(MataDataHandle)
	} else {
		if ok, path, _ := IsValidPath(key); ok {
			if GetFilePathExt(path) == common.FileExtWav1 || GetFilePathExt(path) == common.FileExtWav2 {
				mataDataHandle, err = BuildWavInfoHandle(path)
				if err != nil {
					alog.Logger.Warn("exec BuildExiftoolHandle", zap.Error(err))
					return mataDataHandle
				}
				if mataDataHandle != nil {
					lruCache.Put(key, mataDataHandle)
				}
			} else {
				mataDataHandle, err = BuildExiftoolHandle(path)
				if err != nil {
					alog.Logger.Warn("exec BuildExiftoolHandle", zap.Error(err))
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
