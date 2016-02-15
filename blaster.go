package main

import (
	"encoding/json"
	"fmt"
	"os"
)

const libFile = "./library.json"

func OpenLib(path string) (*Library, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	library := new(Library)
	if err := json.NewDecoder(f).Decode(&library); err != nil {
		return nil, err
	}
	
	return library, nil
}

func GenerateLibrary(path string) error {
	tracks, err := ScanForTracks(path)
	if err != nil {
		return err
	}

	library := Generate(tracks)
	f, err := os.Create(libFile)
	if err != nil {
		return err
	}

	if err := json.NewEncoder(f).Encode(library); err != nil {
		return err
	}
		
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("%s: no command given.\n", os.Args[0])
		return
	}
	
	cmd := os.Args[1]

	if cmd == "generate" {
		err := GenerateLibrary(os.Args[2])
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	if cmd == "serve" {
		lib, err := OpenLib(libFile)
		if err != nil {
			fmt.Println(err)
			return
		}

		ServeAPI(lib)
	}

}
