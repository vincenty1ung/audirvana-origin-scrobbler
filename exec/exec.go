package exec

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cast"
	"go.uber.org/zap"

	alog "github.com/audirvana-origin-scrobbler/log"
)

type ExiftoolInfo map[string]any

func ExiftoolHandle(file string) (*ExiftoolInfo, error) {
	infos := make([]*ExiftoolInfo, 0)
	res := new(ExiftoolInfo)
	ok, file, err := isValidPath(file)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("invalid exiftool path: %s", file)
	}
	command, err := runCommand("exiftool", "-json", file)
	if err != nil {
	}
	err = json.Unmarshal([]byte(command), &infos)
	if err != nil {
		return nil, err
	}
	if len(infos) > 0 {
		res = infos[0]
	}
	return res, nil
}

// GetTrackNumber GetTrackNumber
func (receiver ExiftoolInfo) GetTrackNumber() int64 {
	key1, key2, key3 := "TrackNumber", "Tracknumber", "tracknumber"
	//  "TrackNumber": "1 of 12",
	//   "TrackNumber": 1,
	var val any
	val, ok := receiver[key1]
	if ok {
		return castToInt64(val)
	}
	val, ok = receiver[key2]
	if ok {
		return castToInt64(val)
	}
	val, ok = receiver[key3]
	if ok {
		return castToInt64(val)
	}
	return 0
}

func castToInt64(val any) int64 {
	switch v := val.(type) {
	case string:
		split := strings.Split(v, "of")
		if len(split) > 0 {
			return cast.ToInt64(strings.TrimSpace(split[0]))
		}
	case int, int16, int32, int64, float32, float64, uint, uint8, uint16, uint32, uint64:
		return cast.ToInt64(v)
	}
	return 0
}

// GetMusicBrainzTrackID GetMusicBrainzTrackID
func (receiver ExiftoolInfo) GetMusicBrainzTrackId() string {
	key1, key2 := "MusicbrainzTrackid", "MusicBrainzTrackId"
	var val any
	val, ok := receiver[key1]
	if ok {
		return cast.ToString(val)
	}
	val, ok = receiver[key2]
	if ok {
		return cast.ToString(val)
	}
	return ""
}

func runCommand(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("Error executing command %s: %v\n%s", command, err, output)
	}
	return string(output), nil
}

func isValidPath(path string) (bool, string, error) {
	// 确保路径不是空字符串
	if path == "" {
		return false, "", fmt.Errorf("empty or undefined path")
	}
	path, _ = strings.CutPrefix(path, "file://")
	// 使用 filepath.Clean() 处理任何多余的斜杠或其他非法字符，确保路径整洁。
	path = filepath.Clean(path)
	// 解析符号链接和相对路径到绝对路径
	resolvedPath, err := filepath.EvalSymlinks(path)
	if err != nil {
		alog.Logger.Warn("Cannot resolve symlinks:", zap.Error(err))
		// 检查是否因为文件不存在而失败，如果是则返回。
		var pathError *os.PathError
		isNotExist := errors.As(err, &pathError)
		if isNotExist || os.IsNotExist(err) {
			return false, "", fmt.Errorf("path does not exist: %s", path)
		}
		// 如果不是因为路径不存在导致的错误，则记录并返回
		alog.Logger.Warn("Unknown error occurred while resolving symlinks:", zap.Error(err))

		return false, "", err
	}
	fileInfo, err := os.Stat(resolvedPath)
	if err != nil {
		alog.Logger.Info(
			"Cannot stat the path: %s - Error: ", zap.String("resolvedPath", resolvedPath), zap.Error(err),
		)
		return false, "", fmt.Errorf("error while checking file existence: %s", err)
	}
	// 根据fileInfo.IsDir()判断是文件还是目录
	isDirectory := fileInfo.IsDir()
	alog.Logger.Info(fmt.Sprintf("The path [%s] exists and [%v] a directory", resolvedPath, isDirectory))
	return true, resolvedPath, nil
}
