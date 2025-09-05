package exec

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/go-audio/wav"
	"github.com/spf13/cast"
	"go.uber.org/zap"

	alog "github.com/vincenty1ung/lastfm-scrobbler/log"
)

const (
	MRMediaNowPlayingGet              = "get"
	MRMediaNowPlayingAppRoon          = "com.roon.Roon"
	MRMediaNowPlayingAppMusic         = "com.apple.Music"
	MRMediaNowPlayingBundleIdentifier = "bundleIdentifier"
	MRMediaNowPlayingIsPlaying        = "isPlaying"
	MRMediaNowPlayingAlbum            = "album"
	MRMediaNowPlayingTitle            = "title"
	MRMediaNowPlayingArtist           = "artist"
	MRMediaNowPlayingDuration         = "duration"
	MRMediaNowPlayingElapsedTime      = "elapsedTime"
	MRMediaNowPlayingTimestamp        = "timestamp"
	MRMediaNowPlayingMediaType        = "mediaType"
	MRMediaNowPlayingIsMusicApp       = "isMusicApp"
	MRMediaNowPlayingUniqueIdentifier = "uniqueIdentifier"
)

type (
	MataDataHandle interface {
		GetTitle() string
		GetArtists() string
		GetArtist() string
		GetAlbumartist() string
		GetTrackNumber() int64
		GetMusicBrainzTrackId() string
	}

	ExiftoolInfo map[string]any
	WavInfo      struct {
		wav.Metadata
	}
	MRMediaNowPlaying struct {
		Title            string  `json:"title"`
		Artist           string  `json:"artist"`
		Album            string  `json:"album"`
		IsPlaying        bool    `json:"isPlaying"`
		Duration         float64 `json:"duration"`
		ElapsedTime      float64 `json:"elapsed_time"`
		BundleIdentifier string  `json:"bundleIdentifier"`
	}
)

func BuildExiftoolHandle(file string) (MataDataHandle, error) {
	infos := make([]*ExiftoolInfo, 0)
	res := new(ExiftoolInfo)
	/*ok, file, err := IsValidPath(file)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("invalid exiftool path: %s", file)
	}*/
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

func BuildWavInfoHandle(file string) (MataDataHandle, error) {
	wavInfo := new(WavInfo)
	/*if ok, file, err := IsValidPath(file); err != nil {
		return nil, err
	} else if ok {
		in, err := os.Open(file)
		defer func(in *os.File) {
			err := in.Close()
			if err != nil {
				alog.Error(context.Background(), "Failed to close file", zap.String("file", file), zap.Error(err))
			}
		}(in)
		if err != nil {
			return nil, err
		}
		if mwav := wav.NewDecoder(in); mwav.IsValidFile() {
			mwav.ReadMetadata()
			wavInfo.Metadata = *mwav.Metadata
		}
	}*/
	in, err := os.Open(file)
	defer func(in *os.File) {
		err := in.Close()
		if err != nil {
			alog.Error(context.Background(), "Failed to close file", zap.String("file", file), zap.Error(err))
		}
	}(in)
	if err != nil {
		return nil, err
	}
	if mwav := wav.NewDecoder(in); mwav.IsValidFile() {
		mwav.ReadMetadata()
		wavInfo.Metadata = *mwav.Metadata
	}
	return wavInfo, nil
}

// GetTrackNumber GetTrackNumber
func (receiver ExiftoolInfo) GetTitle() string {
	key1, key2 := "Artists", "artists"
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

// GetTrackNumber GetTrackNumber
func (receiver ExiftoolInfo) GetArtists() string {
	key1, key2 := "Artists", "artists"
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
func (receiver ExiftoolInfo) GetArtist() string {
	key1, key2 := "Artist", "artist"
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
func (receiver ExiftoolInfo) GetAlbumartist() string {
	key1, key2, key3 := "Albumartist", "albumArtist", "AlbumArtist"
	var val any
	val, ok := receiver[key1]
	if ok {
		return cast.ToString(val)
	}
	val, ok = receiver[key2]
	if ok {
		return cast.ToString(val)
	}
	val, ok = receiver[key3]
	if ok {
		return cast.ToString(val)
	}
	return ""
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

// GetTrackNumber GetTrackNumber
func (receiver *WavInfo) GetTitle() string {
	return receiver.Product
}

// GetTrackNumber GetTrackNumber
func (receiver *WavInfo) GetArtists() string {
	return receiver.Artist
}
func (receiver *WavInfo) GetArtist() string {
	return receiver.Artist
}
func (receiver *WavInfo) GetAlbumartist() string {
	return receiver.Artist
}

// GetTrackNumber GetTrackNumber
func (receiver *WavInfo) GetTrackNumber() int64 {
	return castToInt64(receiver.TrackNbr)
}

// GetMusicBrainzTrackID GetMusicBrainzTrackID
func (receiver *WavInfo) GetMusicBrainzTrackId() string {
	return ""
}

func GetMRMediaNowPlaying() (*MRMediaNowPlaying, error) {
	// nowplaying-cli  get album title artist duration elapsedTime timestamp mediaType isMusicApp  uniqueIdentifier
	args := []string{
		MRMediaNowPlayingGet,
		MRMediaNowPlayingAlbum,
		MRMediaNowPlayingTitle,
		MRMediaNowPlayingArtist,
		MRMediaNowPlayingDuration,
		MRMediaNowPlayingElapsedTime,
		MRMediaNowPlayingTimestamp,
		MRMediaNowPlayingMediaType,
		MRMediaNowPlayingIsMusicApp,
		MRMediaNowPlayingUniqueIdentifier,
	}
	curList := map[string]int{
		MRMediaNowPlayingBundleIdentifier: 0,
		MRMediaNowPlayingIsPlaying:        1,
		MRMediaNowPlayingAlbum:            2,
		MRMediaNowPlayingTitle:            3,
		MRMediaNowPlayingArtist:           4,
		MRMediaNowPlayingDuration:         5,
		MRMediaNowPlayingElapsedTime:      6,
		MRMediaNowPlayingTimestamp:        7,
		MRMediaNowPlayingMediaType:        8,
		MRMediaNowPlayingIsMusicApp:       9,
		MRMediaNowPlayingUniqueIdentifier: 10,
	}
	output, err := runCommand(
		"nowplaying-cli-mac", args...,
	)
	if err != nil {
		return nil, err
	}
	MRMediaNowPlayingList := strings.Split(output, "\n")
	var np MRMediaNowPlaying
	if len(MRMediaNowPlayingList) > 10 {
		artists := cast.ToString(MRMediaNowPlayingList[curList[MRMediaNowPlayingArtist]])
		artist := artists
		if artistList := strings.Split(artists, ","); len(artistList) > 0 {
			artist = artistList[0]
		}
		np = MRMediaNowPlaying{
			Title:            cast.ToString(MRMediaNowPlayingList[curList[MRMediaNowPlayingTitle]]),
			Artist:           artist,
			Album:            cast.ToString(MRMediaNowPlayingList[curList[MRMediaNowPlayingAlbum]]),
			IsPlaying:        cast.ToString(MRMediaNowPlayingList[curList[MRMediaNowPlayingIsPlaying]]) == "YES",
			Duration:         cast.ToFloat64(MRMediaNowPlayingList[curList[MRMediaNowPlayingDuration]]),
			ElapsedTime:      cast.ToFloat64(MRMediaNowPlayingList[curList[MRMediaNowPlayingElapsedTime]]),
			BundleIdentifier: cast.ToString(MRMediaNowPlayingList[curList[MRMediaNowPlayingBundleIdentifier]]),
		}
	}
	return &np, nil
}

func castToInt64(val any) int64 {
	switch v := val.(type) {
	case string:
		var toInt64 int64 = 0
		if strings.Contains(v, "/") {
			tmp := strings.Split(v, "/")
			if len(tmp) > 0 {
				toInt64 = cast.ToInt64(strings.TrimSpace(tmp[0]))
			}
		} else if strings.Contains(v, "-") {
			tmp := strings.Split(v, "-")
			if len(tmp) > 0 {
				toInt64 = cast.ToInt64(strings.TrimSpace(tmp[0]))
			}
		} else if strings.Contains(v, "of") {
			tmp := strings.Split(v, "of")
			if len(tmp) > 0 {
				toInt64 = cast.ToInt64(strings.TrimSpace(tmp[0]))
			}
		} else {
			toInt64 = cast.ToInt64(strings.TrimSpace(v))
		}
		return toInt64
	case int, int16, int32, int64, float32, float64, uint, uint8, uint16, uint32, uint64:
		return cast.ToInt64(v)
	}
	return 0
}

func runCommand(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("Error executing command %s: %v\n%s", command, err, output)
	}
	return string(output), nil
}

func IsValidPath(ctx context.Context, path string) (bool, string, error) {
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
		alog.Warn(ctx, "Cannot resolve symlinks:", zap.Error(err))
		// 检查是否因为文件不存在而失败，如果是则返回。
		var pathError *os.PathError
		isNotExist := errors.As(err, &pathError)
		if isNotExist || os.IsNotExist(err) {
			return false, "", fmt.Errorf("path does not exist: %s", path)
		}
		// 如果不是因为路径不存在导致的错误，则记录并返回
		alog.Warn(ctx, "Unknown error occurred while resolving symlinks:", zap.Error(err))

		return false, "", err
	}
	fileInfo, err := os.Stat(resolvedPath)
	if err != nil {
		alog.Info(
			ctx,
			"Cannot stat the path: %s - Error: ", zap.String("resolvedPath", resolvedPath), zap.Error(err),
		)
		return false, "", fmt.Errorf("error while checking file existence: %s", err)
	}
	// 根据fileInfo.IsDir()判断是文件还是目录
	isDirectory := fileInfo.IsDir()
	alog.Info(
		ctx,
		fmt.Sprintf("checkValidPath:The path [%s] exists and [%v] a directory", resolvedPath, isDirectory),
	)
	return true, resolvedPath, nil
}
func GetFilePathExt(path string) string {
	return filepath.Ext(path)
}
