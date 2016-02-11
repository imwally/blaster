package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		return
	}
	cmd := os.Args[1]

	if cmd == "generate" {
		tracks, err := ScanForTracks(os.Args[2])
		if err != nil {
			fmt.Println(err)
		}

		library := new(Library)
		library.Generate(tracks)

		output, err := json.Marshal(library)
		if err != nil {
			fmt.Println(err)
		}
		err = ioutil.WriteFile("library.json", output, 0600)
		if err != nil {
			fmt.Println(err)
		}
	}

	if cmd == "artists" {
		f, err := os.Open("library.json")
		if err != nil {
			fmt.Println(err)
		}

		b, err := ioutil.ReadAll(f)
		if err != nil {
			fmt.Println(err)
		}

		library := new(Library)
		err = json.Unmarshal(b, &library)
		if err != nil {
			fmt.Println(err)
		}

		for _, artist := range library.Artists {
			fmt.Printf("%s\n", artist.Name)
		}
	}

	if cmd == "albums" {
		f, err := os.Open("library.json")
		if err != nil {
			fmt.Println(err)
		}

		b, err := ioutil.ReadAll(f)
		if err != nil {
			fmt.Println(err)
		}

		library := new(Library)
		err = json.Unmarshal(b, &library)
		if err != nil {
			fmt.Println(err)
		}

		for _, album := range library.Albums {
			fmt.Printf("%s - %s\n", album.Artist, album.Title)
			for _, track := range album.Tracks {
				fmt.Printf("\t%d %s\n", track.TrackNumber, track.Title)
			}

		}
	}

}
