package exec

import (
	"context"
	"fmt"
	"testing"

	"github.com/vincenty1ung/lastfm-scrobbler/common"
	"github.com/vincenty1ung/lastfm-scrobbler/core/audirvana"
)

func TestFindMataDataHandleCache(t *testing.T) {
	state, err := audirvana.GetState(context.Background())
	if err != nil {
		t.Error(err)
		return
	}
	if state == common.PlayerStatePlaying {
		audirvanaTrackInfo := audirvana.GetNowPlayingTrackInfo(context.Background())
		handle := FindMataDataHandleCache(context.Background(), audirvanaTrackInfo.Url)
		fmt.Println(handle.GetTitle())
		fmt.Println(handle.GetArtists())
		fmt.Println(handle.GetArtist())
		fmt.Println(handle.GetAlbumartist())
		fmt.Println(handle.GetTrackNumber())
		fmt.Println(handle.GetMusicBrainzTrackId())
		handle = FindMataDataHandleCache(
			context.Background(), "/Users/vincent/Documents/多媒体/音乐/CD/李志/我爱南京/2-05 思念观世音.wav",
		)
		fmt.Println(handle.GetTitle())
		fmt.Println(handle.GetArtists())
		fmt.Println(handle.GetArtist())
		fmt.Println(handle.GetAlbumartist())
		fmt.Println(handle.GetTrackNumber())
		fmt.Println(handle.GetMusicBrainzTrackId())
	}

}
