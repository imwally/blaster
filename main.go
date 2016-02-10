package main

import (
	"fmt"
	"os"
)

func main() {

	db := "./library"
	
	if len(os.Args) < 2 {
		return
	}

	var appdb appDB
	err := appdb.Open(db)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer appdb.Close()
	
	cmd := os.Args[1]
	if cmd == "scan" {
		err := appdb.Initialize(db)
		if err != nil {
			fmt.Println(err)
		}
		
		tracks, err := ScanForTracks(os.Args[2])
		if err != nil {
			fmt.Println(err)
		}

		for _, track := range tracks {
			err = appdb.AddTrack(track)
			if err != nil {
				fmt.Println(err)
			}
		}
		
	}

	if cmd == "byalbum" {
		tracks, err := appdb.Query("Album", os.Args[2])
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(tracks)
	}

	if cmd == "byartist" {
		tracks, err := appdb.Query("Artist", os.Args[2])
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(tracks)
	}
}
