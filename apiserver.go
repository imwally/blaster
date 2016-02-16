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
// all albums.
func (api *API) Albums(w http.ResponseWriter, r *http.Request) {
	Logger(r)
	if err := json.NewEncoder(w).Encode(api.Library.Albums); err != nil {
		http.Error(w, fmt.Sprintf("%s\n", err), 400)
	}
}

// AlbumsBy responds with a json encoded Library.Albums struct
// containing all albums by the artist specified after the /albums/
// path (i.e. /albums/Queen).
func (api *API) AlbumsBy(w http.ResponseWriter, r *http.Request) {
	Logger(r)
	artistName := strings.Replace(r.URL.String(), "/albums/", "", -1)
	unescapedArtistName, err := url.QueryUnescape(artistName)
	if err != nil {
		http.Error(w, fmt.Sprintf("%s\n", err), 400)
		return
	}

	if unescapedArtistName == "" {
		http.Error(w, fmt.Sprintf("no artist specified."), 400)
		return
	}

	response := api.Library.AlbumsBy(unescapedArtistName)
	if response == nil {
		http.Error(w, fmt.Sprintf("no artists by the name %s found.\n", unescapedArtistName), 400)
		return
	}
			
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, fmt.Sprintf("%s\n", err), 400)
	}
}

// Artwork responds with the album cover of the track found at that
// path after /artwork/.
func (api *API) Artwork(w http.ResponseWriter, r *http.Request) {
	Logger(r)
	trackPath := strings.Replace(r.URL.String(), "/artwork/", "", -1)
	unescapedTrackPath, err := url.QueryUnescape(trackPath)
	if err != nil {
		http.Error(w, fmt.Sprintf("%s\n", err), 400)
		return
	}

	if unescapedTrackPath == "" {
		http.Error(w, fmt.Sprintf("no track path specified.\n"), 400)
		return
	}

	art, err := AlbumArt("/" + unescapedTrackPath)
	if err != nil || art == nil {
		http.Error(w, fmt.Sprintf("no artwork for: /%s\n", unescapedTrackPath), 400)
		return
	}

	w.Header().Set("Content-Type", art.MIMEType)
	w.Write(art.Data)
}

// ServeAPI takes a pointer to a Library and serves the api.
func ServeAPI(lib *Library) {
	api := &API{lib}

	http.HandleFunc("/", api.Index)
	http.HandleFunc("/artists", api.Artists)
	http.HandleFunc("/artists/", api.Artists)
	http.HandleFunc("/albums", api.Albums)
	http.HandleFunc("/albums/", api.AlbumsBy)
	http.HandleFunc("/artwork/", api.Artwork)

	log.Println("Blaster API server started on port :8080.")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
