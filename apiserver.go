package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// API holds the Library which is accessed for api reponses.
type API struct {
	Library *Library
}

// endPoints describes api end points.
var endPoints = map[string]string{
	"all_artists":   "/artists",
	"albums":        "/albums{/artist_name}",
	"album_artwork": "/artwork/{track_path}",
}

// Logger is a helper function that prints HTTP request information to
// stdout.
func Logger(r *http.Request) {
	log.Println(r.RemoteAddr, r.Method, r.Proto, r.RequestURI)
}

// Index responds with a map of api end points.
func (api *API) Index(w http.ResponseWriter, r *http.Request) {
	Logger(r)
	if err := json.NewEncoder(w).Encode(endPoints); err != nil {
		http.Error(w, fmt.Sprintf("%s\n", err), 400)
	}
}

// Artists responds with a json encoded Library.Artists struct
// containing all artists.
func (api *API) Artists(w http.ResponseWriter, r *http.Request) {
	Logger(r)
	if err := json.NewEncoder(w).Encode(api.Library.Artists); err != nil {
		http.Error(w, fmt.Sprintf("%s\n", err), 400)
	}
}

// Albums reponds with a json encoded Library.Albums struct containing
// all albums. If an artist name procedes the /albums/ path then only
// albums from that artist will be returned.
func (api *API) Albums(w http.ResponseWriter, r *http.Request) {
	Logger(r)
	artistName := strings.Replace(r.URL.String(), "/albums/", "", -1)
	unescapedArtistName, err := url.QueryUnescape(artistName)
	if err != nil {
		http.Error(w, fmt.Sprintf("%s\n", err), 400)
	}

	var response interface{}
	if unescapedArtistName != "" {
		response = api.Library.AlbumsBy(unescapedArtistName)
	} else {
		response = api.Library.Albums
	}

	e := json.NewEncoder(w)
	if err := e.Encode(response); err != nil {
		http.Error(w, fmt.Sprintf("%s\n", err), 400)
	}
}

// Artwork takes the path after /artwork/ and responds with the album
// cover of the track found at that path.
func (api *API) Artwork(w http.ResponseWriter, r *http.Request) {
	Logger(r)
	trackPath := strings.Trim(r.URL.String(), "/artwork/")
	unescapedTrackPath, err := url.QueryUnescape(trackPath)
	if err != nil {
		http.Error(w, fmt.Sprintf("%s\n", err), 400)
	}

	art, err := AlbumArt("/" + unescapedTrackPath)
	if err != nil {
		http.Error(w, fmt.Sprintf("no artwork for: %s\n", unescapedTrackPath), 400)
	}

	w.Write(art)
}

// ServeAPI takes a pointer to a Library and serves the api.
func ServeAPI(lib *Library) {
	api := &API{lib}

	http.HandleFunc("/", api.Index)
	http.HandleFunc("/artists", api.Artists)
	http.HandleFunc("/artists/", api.Artists)
	http.HandleFunc("/albums", api.Albums)
	http.HandleFunc("/albums/", api.Albums)
	http.HandleFunc("/artwork/", api.Artwork)

	log.Println("Blaster API server started on port :8080.")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
