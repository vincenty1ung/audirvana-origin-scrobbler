package exec

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/go-audio/wav"

	"github.com/audirvana-origin-scrobbler/log"
)

func init() {
	_ = log.LogInit("./logs", "info", make(<-chan struct{}))
}

func TestExecExiftoolHandl(t *testing.T) {
	info, err := BuildExiftoolHandle("file:///Users/vincent/Documents/多媒体/音乐/CD/许巍/此时此刻/02 爱情.m4a")
	if err != nil {
		t.Fatal(err)
	}
	// fmt.Println(info)
	fmt.Println(info.GetMusicBrainzTrackId())
	fmt.Println(info.GetTrackNumber())
	fmt.Println(info.GetArtists())
	fmt.Println(info.GetArtist())
	fmt.Println(info.GetAlbumartist())
	fmt.Println("==============================")
	info, err = BuildExiftoolHandle("file:///Users/vincent/Documents/多媒体/音乐/CD/寸铁/aLIVE IN CHINA 2017-2023/3-01 2021.10.20 濟南 雀躍之地.flac")
	if err != nil {
		t.Fatal(err)
	}
	// fmt.Println(info)
	fmt.Println(info.GetMusicBrainzTrackId())
	fmt.Println(info.GetTrackNumber())
	fmt.Println(info.GetTrackNumber())
	fmt.Println(info.GetArtist())
	fmt.Println(info.GetArtists())
	fmt.Println(info.GetAlbumartist())
	fmt.Println("==============================")
	info, err = BuildExiftoolHandle("file:///Users/vincent/Documents/多媒体/音乐/CD/李志/梵高先生/05 广场.wav")
	if err != nil {
		t.Fatal(err)
	}
	// fmt.Println(info)
	fmt.Println(info.GetMusicBrainzTrackId())
	fmt.Println(info.GetTrackNumber())
	fmt.Println(info.GetArtist())
	fmt.Println(info.GetArtists())
	fmt.Println(info.GetAlbumartist())
	info, err = BuildExiftoolHandle("file:///Users/vincent/Documents/多媒体/音乐/CD/万能青年旅店/张洲/01 张洲.wav")
	if err != nil {
		t.Fatal(err)
	}
	// fmt.Println(info)
	fmt.Println(info.GetMusicBrainzTrackId())
	fmt.Println(info.GetTrackNumber())
	fmt.Println(info.GetArtist())
	fmt.Println(info.GetArtists())
	fmt.Println(info.GetAlbumartist())

}

func TestName(t *testing.T) {
	output, err := runCommand("nowplaying-cli", "get", "title", "album", "artist")
	if err != nil {
		t.Fatal(err)
	}
	split := strings.Split(output, "\n")
	fmt.Println(split)
}

func TestWavInfoHandle(t *testing.T) {
	ok, file, err := IsValidPath("file:///Users/vincent/Documents/多媒体/音乐/CD/万能青年旅店/张洲/01 张洲.wav")
	if err != nil {
		t.Fatal(err)
		return
	}
	if ok {
		in, err := os.Open(file)
		defer func(in *os.File) {
			err := in.Close()
			if err != nil {

			}
		}(in)
		if err != nil {
			t.Fatal(err)
		}
		mwav := wav.NewDecoder(in)
		/*buf, err := d.FullPCMBuffer()
		if err != nil {
			t.Fatal(err)
		}*/
		mwav.ReadInfo()
		fmt.Println(mwav)
	}

	handle, err := BuildWavInfoHandle("file:///Users/vincent/Documents/多媒体/音乐/CD/李志/我爱南京/2-05 思念观世音.wav")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(handle.GetTitle())
	fmt.Println(handle.GetArtists())
	fmt.Println(handle.GetArtist())
	fmt.Println(handle.GetAlbumartist())
	fmt.Println(handle.GetTrackNumber())
	fmt.Println(handle.GetMusicBrainzTrackId())
	handle, err = BuildWavInfoHandle("file:///Users/vincent/Documents/多媒体/音乐/CD/万能青年旅店/张洲/01 张洲.wav")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(handle.GetTitle())
	fmt.Println(handle.GetArtists())
	fmt.Println(handle.GetArtist())
	fmt.Println(handle.GetAlbumartist())
	fmt.Println(handle.GetTrackNumber())
	fmt.Println(handle.GetMusicBrainzTrackId())
}
