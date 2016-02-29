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

// Artist holds an artist's name.
type Artist struct {
	Artist string
	Albums []string
}

// Album holds an album's title.
type Album struct {
	Album string
}

// endPoints describes the available api end points.
var endPoints = map[string]string{
	"all_artists":   "/artists",
	"artist_albums": "/artists/{artist_name}",
	"all_albums":    "/albums",
	"album_tracks":  "/albums/{album_name}",
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

// Artists responds with a json encoded array of all artists.
func (api *API) Artists(w http.ResponseWriter, r *http.Request) {
	Logger(r)
	artists := []Artist{}
	for _, artist := range api.Library.Artists {
		newArtist := Artist{
			Artist: artist,
		}
		albums := api.Library.AlbumsByArtist(artist)
		newArtist.Albums = albums
		artists = append(artists, newArtist)
	}

	if err := json.NewEncoder(w).Encode(artists); err != nil {
		http.Error(w, fmt.Sprintf("%s\n", err), 400)
	}
}

// Albums reponds with a json encoded arary of all albums.
func (api *API) Albums(w http.ResponseWriter, r *http.Request) {
	Logger(r)
	albums := []Album{}
	for _, album := range api.Library.Albums {
		albums = append(albums, Album{album})
	}

	if err := json.NewEncoder(w).Encode(albums); err != nil {
		http.Error(w, fmt.Sprintf("%s\n", err), 400)
	}
}

// AlbumsByArtist respondes with a json encoded array of albums by an
// artist.
func (api *API) AlbumsByArtist(w http.ResponseWriter, r *http.Request) {
	Logger(r)
	query := strings.Replace(r.URL.String(), "/artists/", "", -1)
	unescapedQueryPath, err := url.QueryUnescape(query)
	if err != nil {
		http.Error(w, fmt.Sprintf("%s\n", err), 400)
		return
	}

	if unescapedQueryPath == "" {
		http.Error(w, fmt.Sprintf("no artist specified."), 400)
		return
	}

	albums := []Album{}
	for _, album := range api.Library.AlbumsByArtist(unescapedQueryPath) {
		albums = append(albums, Album{album})
	}

	if err := json.NewEncoder(w).Encode(albums); err != nil {
		http.Error(w, fmt.Sprintf("%s\n", err), 400)
	}
}

// TracksByAlbum respondes with a json encoded array of track objects
// from an album.
func (api *API) TracksByAlbum(w http.ResponseWriter, r *http.Request) {
	Logger(r)
	query := strings.Replace(r.URL.String(), "/albums/", "", -1)
	unescapedQueryPath, err := url.QueryUnescape(query)
	if err != nil {
		http.Error(w, fmt.Sprintf("%s\n", err), 400)
		return
	}

	if unescapedQueryPath == "" {
		http.Error(w, fmt.Sprintf("no album specified."), 400)
		return
	}

	tracks := api.Library.TracksByAlbum(unescapedQueryPath)
	if err := json.NewEncoder(w).Encode(tracks); err != nil {
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

// addCORSHeader is a wrapper function that enables Cross Origin
// Resource Sharing if set is true. It does this by setting the
// Access-Control-Allow-Origin header to the specified origin.
func addCORSHeader(set bool, fn http.HandlerFunc) http.HandlerFunc {
	if set {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:8081")
			fn(w, r)
		}
	}
	return func(w http.ResponseWriter, r *http.Request) {}
}

// ServeAPI takes a pointer to a Library and serves the api.
func ServeAPI(lib *Library, port string, setCORS bool) {
	api := &API{lib}

	http.HandleFunc("/", addCORSHeader(setCORS, api.Index))
	http.HandleFunc("/artists", addCORSHeader(setCORS, api.Artists))
	http.HandleFunc("/artists/", addCORSHeader(setCORS, api.AlbumsByArtist))
	http.HandleFunc("/albums", addCORSHeader(setCORS, api.Albums))
	http.HandleFunc("/albums/", addCORSHeader(setCORS, api.TracksByAlbum))
	http.HandleFunc("/artwork/", addCORSHeader(setCORS, api.Artwork))

	log.Printf("Blaster API server started on port :%s.", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
