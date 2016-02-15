package main

import (
	"encoding/json"
	"fmt"
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
	"albums":        "/albums",
	"album_artwork": "/artwork/{track_path}",
}

// Index responds with a map of api end points.
func (api *API) Index(w http.ResponseWriter, r *http.Request) {
	e := json.NewEncoder(w)
	if err := e.Encode(endPoints); err != nil {
		http.Error(w, fmt.Sprintf("%s\n", err), 400)
	}
}

// Artists responds with json encoded Library.Artists struct.
func (api *API) Artists(w http.ResponseWriter, r *http.Request) {
	e := json.NewEncoder(w)
	if err := e.Encode(api.Library.Artists); err != nil {
		http.Error(w, fmt.Sprintf("%s\n", err), 400)
	}
}

// Albums reponds with a json encoded Library.Albums struct.
func (api *API) Albums(w http.ResponseWriter, r *http.Request) {
	e := json.NewEncoder(w)
	if err := e.Encode(api.Library.Albums); err != nil {
		http.Error(w, fmt.Sprintf("%s\n", err), 400)
	}
}

// Artwork takes the path after /artwork/ and responds with the album
// cover of the track found at that path.
func (api *API) Artwork(w http.ResponseWriter, r *http.Request) {
	location := strings.Trim(r.URL.String(), "/artwork/")
	unescapedLocation, err := url.QueryUnescape(location)
	if err != nil {
		http.Error(w, fmt.Sprintf("%s\n", err), 400)
	}

	art, err := AlbumArt("/" + unescapedLocation)
	if err != nil {
		http.Error(w, fmt.Sprintf("no artwork for: %s\n", unescapedLocation), 400)
	}

	w.Write(art)
}

// ServeAPI takes a pointer to a Library and serves the api.
func ServeAPI(lib *Library) {
	api := &API{lib}

	http.HandleFunc("/", api.Index)
	http.HandleFunc("/artists", api.Artists)
	http.HandleFunc("/albums", api.Albums)
	http.HandleFunc("/artwork/", api.Artwork)
	http.ListenAndServe(":8080", nil)
}
