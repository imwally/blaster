package main

import (
	"log"
	"os"
	"strings"
	"path/filepath"

	"github.com/dhowden/tag"
)

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
	Lyrics string
	Path   string
}

// NewTrack attempts to read the meta data from a file into a Track
// structure. The string "Uknown" will be used in place of a blank
// artist or album field.
func NewTrack(path string) (*Track, error) {
	log.Println(path)
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
		album = "Uknown"
	}

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
// returns a slice of pointers to a Song.
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
