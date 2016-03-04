package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

// API holds the Library which is accessed for api reponses.
type API struct {
	Library *Library
}

// APIError holds an api error message.
type APIError struct {
	Message string
}

// endPoints describes the available api end points.
var endPoints = map[string]string{
	"all_artists":   "/artists",
	"all_albums":    "/albums",
	"artist_albums": "/albums?artist={artist_name}",
	"album_tracks":  "/tracks?album={album_name}",
	"album_artwork": "/artwork?track={track_path}",
}

// Logger is a helper function that prints HTTP request information to
// stdout.
func Logger(r *http.Request) {
	log.Println(r.RemoteAddr, r.Method, r.Proto, r.RequestURI)
}

// JsonResponse JSON encodes i and writes the results to w.
func JsonResponse(w http.ResponseWriter, i interface{}) {
	if err := json.NewEncoder(w).Encode(i); err != nil {
		http.Error(w, fmt.Sprintf("%s\n", err), 400)
	}
}

// Index responds with a map of api end points.
func (api *API) Index(w http.ResponseWriter, r *http.Request) {
	Logger(r)
	JsonResponse(w, endPoints)
}

// Artists responds with a JSON encoded array of all artists.
func (api *API) Artists(w http.ResponseWriter, r *http.Request) {
	Logger(r)
	JsonResponse(w, api.Library.Artists)
}

// Albums responds with a JSON encoded array of all albums or albums
// by an artists if the artist query is specified.
func (api *API) Albums(w http.ResponseWriter, r *http.Request) {
	Logger(r)

	q, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		log.Println(err)
	}

	if len(q) == 0 {
		JsonResponse(w, api.Library.Albums)
		return
	}

	artist := q.Get("artist")
	albums := api.Library.AlbumsByArtist(artist)
	if len(albums) > 0 {
		JsonResponse(w, api.Library.AlbumsByArtist(artist))
		return
	}

	JsonResponse(w, APIError{"no albums found"})
}

// Tracks responds with a JSON encoded array of tracks by an album.
func (api *API) Tracks(w http.ResponseWriter, r *http.Request) {
	Logger(r)

	q, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		log.Println(err)
	}

	album := q.Get("album")
	if len(album) > 0 {
		tracks := api.Library.TracksByAlbum(album)
		if len(tracks) != 0 {
			JsonResponse(w, tracks)
			return
		}
	}

	JsonResponse(w, APIError{"no tracks found"})
}

// Artwork responds with the album cover of a track.
func (api *API) Artwork(w http.ResponseWriter, r *http.Request) {
	Logger(r)

	q, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		log.Println(err)
	}

	track := q.Get("track")
	art, err := AlbumArt(track)
	if err != nil || art == nil {
		JsonResponse(w, APIError{"no artwork found"})
		return
	}

	w.Header().Set("Content-Type", art.MIMEType)
	w.Write(art.Data)
}

// addCORSHeader is a wrapper function that enables Cross Origin
// Resource Sharing if set is true. It does this by setting the
// Access-Control-Allow-Origin header to the specified origin.
func addCORSHeader(set bool, origin string, fn http.HandlerFunc) http.HandlerFunc {
	if set {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			fn(w, r)
		}
	}
	return func(w http.ResponseWriter, r *http.Request) {}
}

// ServeAPI takes a pointer to a Library and serves the api.
func ServeAPI(lib *Library, port string, setCORS bool, origin string) {
	api := &API{lib}

	http.HandleFunc("/", addCORSHeader(setCORS, origin, api.Index))
	http.HandleFunc("/artists", addCORSHeader(setCORS, origin, api.Artists))
	http.HandleFunc("/albums", addCORSHeader(setCORS, origin, api.Albums))
	http.HandleFunc("/tracks", addCORSHeader(setCORS, origin, api.Tracks))
	http.HandleFunc("/artwork", addCORSHeader(setCORS, origin, api.Artwork))

	log.Printf("Blaster API server started on port :%s.", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
