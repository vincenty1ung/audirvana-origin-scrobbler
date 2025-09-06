package musixmatch

import (
	"context"
	"fmt"
	"log"
	"net/http"

	mxm "github.com/milindmadhukar/go-musixmatch"
	"github.com/milindmadhukar/go-musixmatch/params"
)

var (
	client    = http.DefaultClient
	mxmClient *mxm.Client
)

func InitMxmClient(apiKey string) {
	mxmClient = mxm.New(apiKey, client)
}

func SearchArtist(artist string) {
	artists, err := mxmClient.SearchArtist(context.Background(), params.QueryArtist(artist))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(artists)
}
func GetMatcherLyrics(artist, track string) error {
	lyrics, err := mxmClient.GetMatcherLyrics(
		context.Background(), params.QueryTrack(track),
		params.QueryArtist(artist),
	)
	if err != nil {
		return err
	}
	fmt.Println(lyrics)
	return nil
}
