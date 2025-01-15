package exec

import (
	"fmt"
	"testing"
)

func TestExecExiftoolHandl(t *testing.T) {
	info, err := ExiftoolHandle("file:///Users/vincent/Documents/多媒体/音乐/CD/Low Roar/Once in a Long, Long While…/01 Low Roar - Don’t Be so Serious.m4a")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(info)
	fmt.Println(info.GetMusicBrainzTrackId())
	fmt.Println(info.GetTrackNumber())
	fmt.Println("==============================")
	info, err = ExiftoolHandle("file:///Users/vincent/Documents/多媒体/音乐/CD/寸铁/aLIVE IN CHINA 2017-2023/3-01 2021.10.20 濟南 雀躍之地.flac")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(info)
	fmt.Println(info.GetMusicBrainzTrackId())
	fmt.Println(info.GetTrackNumber())
	fmt.Println("==============================")
	info, err = ExiftoolHandle("file:///Users/vincent/Documents/多媒体/音乐/CD/李志/梵高先生/05 广场.wav")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(info)
	fmt.Println(info.GetMusicBrainzTrackId())
	fmt.Println(info.GetTrackNumber())
}
