package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"path"
)

type API struct {
	Library *Library
}

func (api *API) Artists(w http.ResponseWriter, r *http.Request) {

	base := path.Base(r.URL.String())
	unescapedBase, err := url.QueryUnescape(base)
	if err != nil {
		log.Println(err)
	}
	log.Println(unescapedBase)

	var response interface{}
	if unescapedBase != "artists" {
		response = api.Library.AlbumsBy(unescapedBase)
	} else {
		response = api.Library.Artists
	}

	b, err := json.Marshal(response)
	if err != nil {
		http.Error(w, fmt.Sprintf("%s\n", err), 400)
	}

	w.Write(b)
}

func (api *API) AllAlbums(w http.ResponseWriter, r *http.Request) {
	b, err := json.Marshal(api.Library.Albums)
	if err != nil {
		http.Error(w, fmt.Sprintf("%s\n", err), 400)
	}

	w.Write(b)
}

func ServeAPI(lib *Library) {
	var api API
	api.Library = lib

	http.HandleFunc("/api/artists/", api.Artists)
	http.HandleFunc("/api/albums", api.AllAlbums)
	http.ListenAndServe(":8080", nil)
}
