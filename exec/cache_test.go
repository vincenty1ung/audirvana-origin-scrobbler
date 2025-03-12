package exec

import (
	"fmt"
	"testing"

	"github.com/audirvana-origin-scrobbler/audirvana"
	"github.com/audirvana-origin-scrobbler/common"
)

func TestFindMataDataHandleCache(t *testing.T) {
	state, err := audirvana.GetState()
	if err != nil {
		t.Error(err)
		return
	}
	if state == common.PlayerStatePlaying {
		audirvanaTrackInfo := audirvana.GetNowPlayingTrackInfo()
		handle := FindMataDataHandleCache(audirvanaTrackInfo.Url)
		fmt.Println(handle.GetTitle())
		fmt.Println(handle.GetArtists())
		fmt.Println(handle.GetArtist())
		fmt.Println(handle.GetAlbumartist())
		fmt.Println(handle.GetTrackNumber())
		fmt.Println(handle.GetMusicBrainzTrackId())
		handle = FindMataDataHandleCache("/Users/vincent/Documents/多媒体/音乐/CD/李志/我爱南京/2-05 思念观世音.wav")
		fmt.Println(handle.GetTitle())
		fmt.Println(handle.GetArtists())
		fmt.Println(handle.GetArtist())
		fmt.Println(handle.GetAlbumartist())
		fmt.Println(handle.GetTrackNumber())
		fmt.Println(handle.GetMusicBrainzTrackId())
	}

}
