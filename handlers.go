package main

import (
	"TopSongTracks/api"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
)

var response api.Response

func FavouriteArtistHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	country := strings.TrimPrefix(path, "/")

	wg := sync.WaitGroup{}

	wg.Add(1)
	go response.LastFmTrackInfo(&wg, country)
	wg.Add(1)
	go response.MusixmatchSearchTrackID(&wg)
	wg.Add(1)
	go response.SuggestedTracksBasedOnTrackAndArtist(&wg)
	wg.Add(1)
	go response.MusixmatchSearchLyricsByTrackID(&wg)

	wg.Wait()

	jsonData, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		fmt.Println("couldn't encode response data to json format", err)
	}

	statusCode := http.StatusAccepted
	if response.TrackErr == true {
		statusCode = http.StatusNotFound
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(jsonData)

}
