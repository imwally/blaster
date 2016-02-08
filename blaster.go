package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/dhowden/tag"
)

const musicDir = ""

type Library struct {
	Artists []*Artist
	Albums  []*Album
	Songs   []*Song
}

type Artist struct {
	Name   string
	Albums []*Album
	Songs  []*Song
}

type Album struct {
	Title  string
	Artist string
	Year   int
	Genre  string
	Songs  []*Song
}

type Song struct {
	FileType    tag.FileType
	Title       string
	Album       string
	Artist      string
	AlbumArtist string
	Composer    string
	Year        int
	Genre       string
	//	Track
	//	Disc
	//	Picture *tag.Picture
	Lyrics string
	Path   string
}

// SongsByArtist returns all songs by an Artist.
func (lib *Library) SongsByArtist(artistName string) []*Song {
	for _, artist := range lib.Artists {
		if artist.Name == artistName {
			return artist.Songs
		}
	}

	return nil
}

// AlbumsByArtist returns all albums by an Artist.
func (lib *Library) AlbumsByArtist(artistName string) []*Album {
	for _, artist := range lib.Artists {
		if artist.Name == artistName {
			fmt.Println(artist.Albums)
			return artist.Albums
		}
	}

	return nil
}

// Organize organizes the library into an hierarchical structure of
// Artists, Albums, and Songs.
func Organize(lib *Library) {
	artists := make(map[string]*Artist)
	albums := make(map[string]*Album)
	albumSongs := make(map[string][]*Song)

	for _, song := range lib.Songs {
		albums[song.Album] = &Album{
			Title:  song.Album,
			Artist: song.Artist,
			Year:   song.Year,
			Genre:  song.Genre,
		}
		albumSongs[song.Album] = append(albumSongs[song.Album], song)
		
		artists[song.Artist] = &Artist{
			Name: song.Artist,
		}
	}

	for _, artist := range artists {
		lib.Artists = append(lib.Artists, artist)
	}
	
	for albumName, album := range albums {
		album.Songs = albumSongs[albumName]
		lib.Albums = append(lib.Albums, album)
	}
}

// NewSong attempts to read the meta data of file and apply it to a
// Song structure.
func NewSong(path string) (*Song, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	m, err := tag.ReadFrom(f)
	if err != nil {
		return nil, err
	}

	artist := m.Artist()
	if artist == "" {
		artist = "Unknown"
	}

	album := m.Album()
	if album == "" {
		album = "Uknown"
	}

	s := &Song{
		m.FileType(),
		m.Title(),
		m.Album(),
		artist,
		m.AlbumArtist(),
		m.Composer(),
		m.Year(),
		m.Genre(),
		m.Lyrics(),
		path,
	}

	return s, nil
}

// GenerateLibrary expects a slice of valid song file paths and calls
// newSong() on each path. After all songs have been added the library
// is organized into Artists, Albums, and Songs.
func GenerateLibrary(songs []string) *Library {
	lib := new(Library)
	for _, song := range songs {
		s, err := NewSong(song)
		if err != nil {
			log.Println(err)
		}
		// Add song to the library.
		lib.Songs = append(lib.Songs, s)
	}

	// Organize artists, albums, and songs in library.
	Organize(lib)

	return lib
}

// ScanForSongs scans a directory and returns a slice of files that
// contain valid meta data of the format MP3, MP4, Flac, and Ogg.
func ScanForSongs(path string) ([]string, error) {
	var songs []string
	findMusic := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}

		if !info.IsDir() {
			_, err = tag.ReadFrom(f)
			if err != tag.ErrNoTagsFound {
				songs = append(songs, path)
			}
		}

		return nil
	}

	err := filepath.Walk(path, findMusic)
	if err != nil {
		log.Println(path, err)
	}

	return songs, nil
}

func main() {
	songsFound, err := ScanForSongs(musicDir)
	if err != nil {
		fmt.Println(err)
	}

	lib := GenerateLibrary(songsFound)

	b, err := json.MarshalIndent(lib, "", "\t")
	if err != nil {
		fmt.Println(err)
	}

	os.Stdout.Write(b)
}
