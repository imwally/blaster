package main

import (
	"encoding/json"
	"html"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// API holds the main entry point into the library.
type API struct {
	Library *Library
}

// APIError holds an api error message.
type APIError struct {
	Status  string
	Message string
}

// endPoints describes the available api end points.
var endPoints = map[string]string{
	"all_artists":   "/artists",
	"all_albums":    "/albums",
	"artist_albums": "/albums?artist={artist_name}",
	"album_tracks":  "/tracks?album={album_name}",
	"album_artwork": "/artwork?track={track_path}",
	"audio_file":    "/file/{track_path}",
}

// Logger is a helper function that prints HTTP request information to
// stdout.
func Logger(r *http.Request) {
	log.Println(r.RemoteAddr, r.Method, r.Proto, r.RequestURI)
}

// JSONResponse JSON encodes i and writes the results to w.
func JSONResponse(w http.ResponseWriter, i interface{}) {
	if err := json.NewEncoder(w).Encode(i); err != nil {
		JSONResponse(w, APIError{"error", err.Error()})
		return
	}
}

// GetQuery unescapes and parses a URL query. If successful it returns
// the map of queries as url.Values.
func GetQuery(u string) (url.Values, error) {
	q, err := url.ParseQuery(html.UnescapeString(u))
	if err != nil {
		return nil, err
	}

	return q, nil
}

// Index responds with a map of api end points.
func (api *API) Index(w http.ResponseWriter, r *http.Request) {
	Logger(r)
	JSONResponse(w, endPoints)
}

// Artists responds with a JSON encoded array of all artists.
func (api *API) Artists(w http.ResponseWriter, r *http.Request) {
	Logger(r)
	JSONResponse(w, api.Library.Artists)
}

// Albums responds with a JSON encoded array of all albums or albums
// by an artists if the artist query is specified.
func (api *API) Albums(w http.ResponseWriter, r *http.Request) {
	Logger(r)
	q, err := GetQuery(r.URL.RawQuery)
	if err != nil {
		JSONResponse(w, APIError{"error", err.Error()})
		return
	}

	if len(q) == 0 {
		JSONResponse(w, api.Library.Albums)
		return
	}

	artist := q.Get("artist")
	albums := api.Library.AlbumsByArtist(artist)
	if len(albums) > 0 {
		JSONResponse(w, albums)
		return
	}

	JSONResponse(w, APIError{"error", "no albums found"})
}

// Tracks responds with a JSON encoded array of tracks by an album.
func (api *API) Tracks(w http.ResponseWriter, r *http.Request) {
	Logger(r)
	q, err := GetQuery(r.URL.RawQuery)
	if err != nil {
		JSONResponse(w, APIError{"error", err.Error()})
		return
	}

	album := q.Get("album")
	tracks := api.Library.TracksByAlbum(album)
	if len(tracks) != 0 {
		JSONResponse(w, tracks)
		return
	}

	JSONResponse(w, APIError{"error", "no tracks found"})
}

// Artwork responds with the album cover of a track.
func (api *API) Artwork(w http.ResponseWriter, r *http.Request) {
	Logger(r)
	q, err := GetQuery(r.URL.RawQuery)
	if err != nil {
		JSONResponse(w, APIError{"error", err.Error()})
		return
	}

	track := q.Get("track")
	art, err := AlbumArt(track)
	if err != nil || art == nil {
		JSONResponse(w, APIError{"error", "no artwork found"})
		return
	}

	w.Header().Set("Content-Type", art.MIMEType)
	w.Write(art.Data)
}

// ServeAudio serve's the audio file specified in the path after
// /file. It expects the path after /file to begin with the location
// where the music library was generated
// (i.e. /file/home/user/Music/artist/album/song.mp3). After /file is
// stripped the remaining path is passed to http.ServeFile.
func (api *API) ServeAudio(w http.ResponseWriter, r *http.Request) {
	Logger(r)

	track := strings.TrimPrefix(r.URL.String(), "/file"+api.Library.Path)
	track1, err := url.QueryUnescape(track)
	if err != nil {
		JSONResponse(w, APIError{"error", err.Error()})
	}

	http.ServeFile(w, r, api.Library.Path+track1)
}

// addCORSHeader is a wrapper function that enables Cross Origin
// Resource Sharing if set is true. It does this by setting the
// Access-Control-Allow-Origin header to the specified origin.
func addCORSHeader(set bool, origin string, fn http.HandlerFunc) http.HandlerFunc {
	if !set {
		return func(w http.ResponseWriter, r *http.Request) {}
	}

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		fn(w, r)
	}
}

// ServeAPI takes a pointer to a Library and serves the api.
func ServeAPI(lib *Library, port string, setCORS bool, origin string) {
	api := &API{lib}

	http.HandleFunc("/", addCORSHeader(setCORS, origin, api.Index))
	http.HandleFunc("/artists", addCORSHeader(setCORS, origin, api.Artists))
	http.HandleFunc("/albums", addCORSHeader(setCORS, origin, api.Albums))
	http.HandleFunc("/tracks", addCORSHeader(setCORS, origin, api.Tracks))
	http.HandleFunc("/artwork", addCORSHeader(setCORS, origin, api.Artwork))
	http.HandleFunc("/file/", addCORSHeader(setCORS, origin, api.ServeAudio))

	log.Printf("Blaster API server started on port %s.", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
