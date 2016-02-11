package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		return
	}

	tracks, err := ScanForTracks(os.Args[1])
	if err != nil {
		fmt.Println(err)
	}

	cmd := os.Args[2]

	if cmd == "artists" {
		artists := Artists(tracks)
		for _, artist := range artists {
			fmt.Println(artist.Name)
		}
	}

	if cmd == "albums" {
		albums := Albums(tracks, nil)
		for _, album := range albums {
			fmt.Printf("%s - %s\n", album.Artist, album.Title)
		}
	}

	if cmd == "albumsby" {
		albums := Albums(tracks, &Artist{os.Args[3]})
		for _, album := range albums {
			fmt.Printf("%s - %s\n", album.Artist, album.Title)
			for _, track := range album.Tracks {
				fmt.Printf("\t%d %s\n", track.TrackNumber, track.Title)
			}
		}
	}
}
