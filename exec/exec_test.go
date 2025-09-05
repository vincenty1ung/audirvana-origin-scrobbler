package exec

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/go-audio/wav"

	"github.com/vincenty1ung/lastfm-scrobbler/log"
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
	ok, file, err := IsValidPath(nil, "file:///Users/vincent/Documents/多媒体/音乐/CD/万能青年旅店/张洲/01 张洲.wav")
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

func TestNam1(t *testing.T) {
	fmt.Println("公司系统显示入职前已经有100个月的工龄")
	date := time.Date(2023, 8, 30, 0, 0, 0, 5e8, time.Local)
	fmt.Println(date)
	fmt.Println("再加20月到10年")
	date = date.AddDate(0, 20, 0)
	fmt.Println(date)
	fmt.Println("倒推公司显示的工龄开始日期")
	date = date.AddDate(0, -120, 0)
	fmt.Println(date)

	date = time.Date(2016, 5, 1, 0, 0, 0, 5e8, time.Local)
	addDate := date.AddDate(0, 120, 0)
	fmt.Print("深圳：2016年5月开始工作：工作十年到期是？")
	fmt.Println(addDate)

	date = time.Date(2015, 1, 1, 0, 0, 0, 5e8, time.Local)
	addDate = date.AddDate(0, 120, 0)
	fmt.Print("青岛：2015开始工作：工作十年到期是？")
	fmt.Println(addDate)

}

func TestGetMRMediaNowPlaying(t *testing.T) {
	nowPlaying, err := GetMRMediaNowPlaying()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(nowPlaying)
}

func TestEnv(t *testing.T) {
	getenv := os.Getenv("PATH")
	fmt.Println(getenv)
	err := os.Setenv("PATH", getenv+":./shell/bin")
	if err != nil {
		return
	}
	getenv = os.Getenv("PATH")
	fmt.Println(getenv)
}
