package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/dhowden/tag"
)

// Library is the entry point to all artists, albums and tracks.
type Library struct {
	Artists []*Artist
	Albums  []*Album
	Tracks  []*Track
}

// Artist holds an artist's name.
type Artist struct {
	Name string
}

// Album holds an album's title, artist name, and tracks.
type Album struct {
	Title  string
	Artist string
	Tracks []*Track
}

// Track holds the meta data and path to a a track.
type Track struct {
	FileType    tag.FileType
	Title       string
	Album       string
	Artist      string
	AlbumArtist string
	Composer    string
	Genre       string
	Year        int
	TrackNumber int
	TrackTotal  int
	DiscNumber  int
	DiscTotal   int
	Lyrics      string
	Path        string
}

// Generate reads all artists, albums and tracks into the Library.
func (l *Library) Generate(tracks []*Track) {
	l.Artists = Artists(tracks)
	l.Albums = Albums(tracks)
	l.Tracks = tracks
}

// Artists returns a slice of unique artists from a slice of Tracks.
func Artists(tracks []*Track) []*Artist {
	artists := make(map[string]*Artist)
	for _, track := range tracks {
		artists[track.Artist] = &Artist{
			Name: track.Artist,
		}
	}

	var uniqueArtists []*Artist
	for _, artist := range artists {
		uniqueArtists = append(uniqueArtists, artist)
	}

	return uniqueArtists
}

// Albums returns a slice of unique albums from an artist. If artist
// is nil all albums are returned.
func Albums(tracks []*Track) []*Album {
	albums := make(map[string]*Album)
	albumTracks := make(map[string][]*Track)
	for _, track := range tracks {
		albumTracks[track.Album] = append(albumTracks[track.Album], track)
		albums[track.Album] = &Album{
			Title:  track.Album,
			Artist: track.Artist,
		}
	}

	var uniqueAlbums []*Album
	for _, album := range albums {
		album.Tracks = albumTracks[album.Title]
		uniqueAlbums = append(uniqueAlbums, album)
	}

	return uniqueAlbums
}

// NewTrack attempts to read the meta data from a file into a Track
// structure. The string "Uknown" will be used in place of a blank
// artist or album field.
func NewTrack(path string) (*Track, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	m, err := tag.ReadFrom(f)
	if err != nil {
		return nil, err
	}

	title := m.Title()
	if title == "" {
		title = "Untitled"
	}

	artist := m.Artist()
	if artist == "" {
		artist = "Unknown"
	}

	album := m.Album()
	if album == "" {
		album = "Untitled"
	}

	fmt.Printf(" Adding: %s - %s\r", artist, title)

	trackNumber, trackTotal := m.Track()
	discNumber, discTotal := m.Disc()

	t := &Track{
		m.FileType(),
		title,
		album,
		artist,
		m.AlbumArtist(),
		m.Composer(),
		m.Genre(),
		m.Year(),
		trackNumber,
		trackTotal,
		discNumber,
		discTotal,
		m.Lyrics(),
		"",
	}

	return t, nil
}

// ScanForTracks scans a path for files containing valid meta tag
// information. Valid meta tag formats are ID3, MP4, and Vorbis. It
// returns a slice of pointers to a Track.
func ScanForTracks(path string) ([]*Track, error) {
	var tracks []*Track

	findSong := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and files with unsupported
		// extensions.
		if !info.IsDir() && SupportedExtention(path) {
			track, err := NewTrack(path)
			if err != nil {
				log.Printf("%s: %s", err, path)
			}
			track.Path = path
			tracks = append(tracks, track)
		}

		return nil
	}

	err := filepath.Walk(path, findSong)
	if err != nil {
		log.Println(path, err)
	}

	return tracks, nil
}

// SupportedExtension returns true if the extension of the path is
// supported.
func SupportedExtention(path string) bool {
	supported := []string{".mp3", ".mp4", ".m4a", ".flac", ".aac"}
	for _, supported_ext := range supported {
		ext := strings.EqualFold(filepath.Ext(path), supported_ext)
		if ext {
			return true
		}
	}

	return false
}
