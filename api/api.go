package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"reflect"
	"strings"
	"sync"
)

type Response struct {
	TopTrackOfTheRegion          string          `json:"name"`
	Lyrics                       string          `json:"lyrics"`
	ArtistInfoAndImage           Artist          `json:"artist_info"`
	SuggestedTracksBasedOnSearch SuggestedTracks `json:"suggested_tracks"`
	TrackErr                     bool            `json:"track_error"`
	TrackMessage                 string          `json:"track_msg"`
	SuggestedTracksErr           bool            `json:"suggested_track_error"`
	SuggestedTracksMessage       string          `json:"suggested_track_msg"`
}

type TopTracks struct {
	Tracks struct {
		Track []struct {
			Name   string `json:"name"`
			URL    string `json:"url"`
			Artist struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"artist"`
			Image []struct {
				Text string `json:"#text"`
			} `json:"image"`
		} `json:"track"`
	} `json:"tracks"`
}

type Artist struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

var artistInfo Artist

var topTrackchan = make(chan TopTracks)

func (res *Response) LastFmTrackInfo(wg *sync.WaitGroup, countryName string) {
	defer wg.Done()

	// country := "spain"
	country := countryName
	lastFmkey := os.Getenv("lastFm")

	limit := 1 // number of tracks to be fetched per page. We're trying to fetch top result of the country
	// url := "https://ws.audioscrobbler.com/2.0/?method=geo.gettoptracks&country=spain&api_key=YOUR_API_KEY&format=json"

	resp, err := http.Get(fmt.Sprintf("https://ws.audioscrobbler.com/2.0/?method=geo.gettoptracks&country=%s&limit=%d&api_key=%s&format=json", country, limit, lastFmkey))
	if err != nil {
		log.Println("couldn't fetch data from API:", err)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("couldn't read response body:", err)
	}

	var track TopTracks

	err = json.Unmarshal(data, &track)
	if err != nil {
		fmt.Println("couldn't decode json data:", err)
	}

	// fmt.Println("track:", track)

	if reflect.ValueOf(track).IsZero() {
		res.TrackErr = true
		res.TrackMessage = "Country name should be as defined by the ISO 3166-1 country names standard"

		artistInfo.Name = ""
		artistInfo.URL = ""

		res.ArtistInfoAndImage = artistInfo

		res.TopTrackOfTheRegion = ""
	} else {
		res.TrackErr = false
		res.TrackMessage = ""

		artistInfo.Name = track.Tracks.Track[0].Artist.Name
		artistInfo.URL = track.Tracks.Track[0].Artist.URL

		res.ArtistInfoAndImage = artistInfo

		res.TopTrackOfTheRegion = track.Tracks.Track[0].Name
	}

	// fmt.Println(track)
	topTrackchan <- track

}

type SuggestedTracks struct {
	TrackList []struct {
		Track struct {
			TrackID       int    `json:"track_id"`
			TrackName     string `json:"track_name"`
			TrackShareURL string `json:"track_share_url"`
		} `json:"track"`
	} `json:"track_list"`
}

var suggestedTrackschan = make(chan SuggestedTracks)

type MusixmatchSearchTrackIDs struct {
	Message struct {
		Header struct {
			StatusCode  int     `json:"status_code"`
			ExecuteTime float64 `json:"execute_time"`
			Available   int     `json:"available"`
		} `json:"header"`
		Body struct {
			TrackList []struct {
				Track struct {
					TrackID       int    `json:"track_id"`
					TrackName     string `json:"track_name"`
					TrackShareURL string `json:"track_share_url"`
				} `json:"track"`
			} `json:"track_list"`
		} `json:"body"`
	} `json:"message"`
}

var topTrackMusicSearchTrackID = make(chan int)

func (res *Response) MusixmatchSearchTrackID(wg *sync.WaitGroup) {
	defer wg.Done()

	topTrackInfo := <-topTrackchan

	var trackName string
	var artistName string
	if reflect.ValueOf(topTrackInfo).IsZero() {
		trackName = ""
		artistName = ""
	} else {
		trackName = topTrackInfo.Tracks.Track[0].Name
		artistName = topTrackInfo.Tracks.Track[0].Artist.Name
	}
	// trackName := "Mr. Brightside"
	// trackName := topTrackInfo.Tracks.Track[0].Name
	trackNameTemp := strings.Split(trackName, " ")
	trackName = strings.Join(trackNameTemp, "%20")
	// artistName := "The Killers"
	// artistName := topTrackInfo.Tracks.Track[0].Artist.Name
	// artistURL := topTrackInfo.Tracks.Track[0].Artist.URL   // Or artistInfo
	artistNameTemp := strings.Split(artistName, " ")
	artistName = strings.Join(artistNameTemp, "%20")
	// lyrics := "Mr. Brightside"
	lyrics := trackName
	musixMatchkey := os.Getenv("musixMatch")

	// musixmatchURL := "https://api.musixmatch.com/ws/1.1/track.search?q_track=Mr.%20Brightside&q_artist=The%20Killers&q_lyrics=Mr.%20Brightside&apikey=35b9cf320b8da2e1927a573f902c4e7f"
	resp, err := http.Get(fmt.Sprintf("https://api.musixmatch.com/ws/1.1/track.search?q_track=%s&q_artist=%s&q_lyrics=%s&apikey=%s", trackName, artistName, lyrics, musixMatchkey))
	// resp, err := http.Get(musixmatchURL)
	if err != nil {
		log.Println("couldn't fetch data from API:", err)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("couldn't read response body:", err)
	}

	var trackIDs MusixmatchSearchTrackIDs

	err = json.Unmarshal(data, &trackIDs)
	if err != nil {
		fmt.Println("couldn't decode json data:", err)
	}

	var suggestedTracks SuggestedTracks

	if trackIDs.Message.Header.Available == 0 {
		topTrackID := -1 // error, the searched track is not available in Musixmatch API
		topTrackMusicSearchTrackID <- topTrackID
		res.SuggestedTracksErr = true
		res.SuggestedTracksMessage = "the searched track is not available in Musixmatch API Repository"
		suggestedTrackschan <- suggestedTracks

	} else {
		if trackName == "" && artistName == "" {
			res.SuggestedTracksMessage = "Although provided country name is invalid!!! still some tracks are suggested"
		} else {
			res.SuggestedTracksErr = false
			res.SuggestedTracksMessage = ""
		}
		// res.SuggestedTracksErr = false
		// res.SuggestedTracksMessage = ""

		topTrackID := trackIDs.Message.Body.TrackList[0].Track.TrackID
		// fmt.Println("The first track id:", topTrackID)
		topTrackMusicSearchTrackID <- topTrackID

		fmt.Println("===================== Musixmatch track id =============================")
		// fmt.Println(trackIDs)

		// var suggestedTracks SuggestedTracks

		for _, val := range trackIDs.Message.Body.TrackList {
			suggestedTracks.TrackList = append(suggestedTracks.TrackList, val)
		}

		// fmt.Println("Suggested tracks are:", suggestedTracks.TrackList)

		suggestedTrackschan <- suggestedTracks
	}

}

func (res *Response) SuggestedTracksBasedOnTrackAndArtist(wg *sync.WaitGroup) {
	defer wg.Done()
	allSuggestedTracks := <-suggestedTrackschan
	res.SuggestedTracksBasedOnSearch = allSuggestedTracks
	// fmt.Println("============ Suggested Tracks Based on Searched Track And Artist=========")
	// fmt.Println(allSuggestedTracks)
}

type MusixmatchSearchLyrics struct {
	Message struct {
		Header struct {
			StatusCode  int     `json:"status_code"`
			ExecuteTime float64 `json:"execute_time"`
		} `json:"header"`
		Body struct {
			Lyrics struct {
				LyricsID   int    `json:"lyrics_id"`
				Explicit   int    `json:"explicit"`
				LyricsBody string `json:"lyrics_body"`
				// ScriptTrackingURL string    `json:"script_tracking_url"`
				// PixelTrackingURL  string    `json:"pixel_tracking_url"`
				// LyricsCopyright   string    `json:"lyrics_copyright"`
				// UpdatedTime       time.Time `json:"updated_time"`
			} `json:"lyrics"`
		} `json:"body"`
	} `json:"message"`
}

func (res *Response) MusixmatchSearchLyricsByTrackID(wg *sync.WaitGroup) {
	defer wg.Done()

	// trackID := "172729233"
	trackIDtemp := <-topTrackMusicSearchTrackID
	trackID := fmt.Sprintf("%d", trackIDtemp)
	musixMatchkey := os.Getenv("musixMatch")
	// musixmatchURL := "https://api.musixmatch.com/ws/1.1/track.lyrics.get?format=json&track_id=<string>&api_key=<string>"

	// "https://api.musixmatch.com/ws/1.1/track.lyrics.get?format=json&track_id=172729233&apikey=35b9cf320b8da2e1927a573f902c4e7f"

	resp, err := http.Get(fmt.Sprintf("https://api.musixmatch.com/ws/1.1/track.lyrics.get?format=json&track_id=%s&apikey=%s", trackID, musixMatchkey))
	if err != nil {
		log.Println("couldn't fetch data from API:", err)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("couldn't read response body:", err)
	}

	var lyrics MusixmatchSearchLyrics

	err = json.Unmarshal(data, &lyrics)
	if err != nil {
		fmt.Println("couldn't decode json data:", err)
	}
	// fmt.Println("===================== Musixmatch Lyrics of the given track id =============================")
	// fmt.Println(lyrics)
	// fmt.Println(lyrics.Message.Body.Lyrics.LyricsBody)
	res.Lyrics = lyrics.Message.Body.Lyrics.LyricsBody
}
