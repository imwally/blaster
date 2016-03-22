package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/dhowden/tag"
)

// Library is the entry point to all artists, albums and tracks.
type Library struct {
	Artists []string
	Albums  []string
	Tracks  []*Track
	Path    string
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

// AlbumsByArtist returns a slice of album titles by an artist.
func (l *Library) AlbumsByArtist(artist string) []string {
	albums := make(map[string]int)
	for _, track := range l.Tracks {
		artistFound := strings.EqualFold(track.Artist, artist)
		if artistFound {
			albums[track.Album] = 0
		}
	}

	albumsBy := []string{}
	for album := range albums {
		albumsBy = append(albumsBy, album)
	}

	return albumsBy
}

// TracksByAlbum returns a slice of pointers to Tracks that appear on
// an album.
func (l *Library) TracksByAlbum(album string) []*Track {
	tracks := []*Track{}
	for _, track := range l.Tracks {
		albumFound := strings.EqualFold(track.Album, album)
		if albumFound {
			tracks = append(tracks, track)
		}
	}

	return tracks
}

// Generate reads all artists, albums and tracks into the Library.
func Generate(tracks []*Track) *Library {
	artists := Artists(tracks)
	albums := Albums(tracks)
	sort.Strings(artists)
	sort.Strings(albums)

	l := new(Library)
	l.Artists = artists
	l.Albums = albums
	l.Tracks = tracks

	return l
}

// Artists returns a slice of unique artist names from a slice of
// tracks.
func Artists(tracks []*Track) []string {
	uniqueArtists := make(map[string]int)
	for _, v := range tracks {
		uniqueArtists[v.Artist] = 0
	}

	artists := []string{}
	for artist := range uniqueArtists {
		artists = append(artists, artist)
	}

	return artists
}

// Albums returns a slice of unique album titles from a slice of
// tracks.
func Albums(tracks []*Track) []string {
	uniqueAlbums := make(map[string]int)
	for _, v := range tracks {
		uniqueAlbums[v.Album] = 0
	}

	albums := []string{}
	for album := range uniqueAlbums {
		albums = append(albums, album)
	}

	return albums
}

// AlbumArt takes a path to an audio file and returns the album
// artwork if found.
func AlbumArt(path string) (*tag.Picture, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	m, err := tag.ReadFrom(f)
	if err != nil {
		return nil, err
	}

	if m.Picture() == nil {
		return nil, errors.New("No artwork found.")
	}

	return m.Picture(), nil
}

// NewTrack attempts to read the meta data from a file into a Track
// data structure. The string "Unknown" will be used in place of a
// blank artist name and "Untitled" in place of blank album and track
// titles.
func NewTrack(path string) (*Track, error) {
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

	title := m.Title()
	if title == "" {
		title = "Untitled"
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
// information. Valid meta tag formats are ID3, MP4, FLAC, and
// Vorbis. It returns a slice of pointers to a Track.
func ScanForTracks(path string) ([]*Track, error) {
	var tracks []*Track

	findSong := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and files with unsupported
		// extensions.
		if !info.IsDir() && SupportedExtension(path) {
			track, err := NewTrack(path)
			if err != nil {
				log.Printf("%s: %s", err, path)
			}
			absPath, err := filepath.Abs(path)
			if err != nil {
				return err
			}
			track.Path = absPath
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
func SupportedExtension(path string) bool {
	supported := []string{
		".aac",
		".flac",
		".mp3",
		".m4a",
		".mp4",
		".ogg",
		".ogv",
	}
	for _, supportedExt := range supported {
		ext := strings.EqualFold(filepath.Ext(path), supportedExt)
		if ext {
			return true
		}
	}

	return false
}
