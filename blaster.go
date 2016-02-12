package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

const libFile string = "./library.json"

func OpenLib(path string) (*Library, error) {
	f, err := os.Open(libFile)
	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	library := new(Library)
	err = json.Unmarshal(b, &library)
	if err != nil {
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
	output, err := json.Marshal(library)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(libFile, output, 0600)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	if len(os.Args) < 2 {
		return
	}
	cmd := os.Args[1]

	if cmd == "generate" {
		err := GenerateLibrary(os.Args[2])
		if err != nil {
			fmt.Println(err)
		}
	}

	if cmd == "serve" {
		lib, err := OpenLib("./library.json")
		if err != nil {
			fmt.Println(err)
		}

		ServeAPI(lib)
	}
}
