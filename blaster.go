package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

const libName = "blaster.json"
const usage = `Usage: blaster [command] [options]
    
Commands:
    generate PATH              Generate a new library from a music path.
    serve PATH PORT            Serve HTTP API server on PORT.

Options:
    -origin                    The origin used in the Access-Control-Allow-Origin header.

Examples:

    Create a new library containing all the tracks found in your
    ~/Music directory. A new file named blaster.json will be created
    in the same directory.
   
        $ generate ~/Music           

    Serve an HTTP API server on port 8080.
   
        $ serve ~/Music/blaser.json 8080 -origin "http://127.0.0.1:8081"
`

// OpenLib takes a path to a JSON encoded music library and returns a
// pointer to a new Library.
func OpenLib(libPath string) (*Library, error) {
	f, err := os.Open(libPath)
	if err != nil {
		return nil, err
	}

	library := new(Library)
	if err := json.NewDecoder(f).Decode(&library); err != nil {
		return nil, err
	}

	return library, nil
}

// GenerateLibrary takes a path to a music directory containing audio
// files and generates a JSON encoded music library.
func GenerateLibrary(musicPath string) error {
	tracks, err := ScanForTracks(musicPath)
	if err != nil {
		return err
	}

	library := Generate(tracks)
	f, err := os.Create(musicPath + "/" + libName)
	if err != nil {
		return err
	}

	if err := json.NewEncoder(f).Encode(library); err != nil {
		return err
	}

	return nil
}

func main() {
	app := os.Args[0]
	if len(os.Args) < 2 {
		fmt.Printf("%s: no command given.\n", app)
		return
	}

	cmd := os.Args[1]
	if cmd == "help" {
		fmt.Println(usage)
		return
	}

	if cmd == "generate" {
		if len(os.Args) < 3 {
			fmt.Printf("%s: not enough arguments.\n", app)
			fmt.Println(usage)
			return
		}
		musicPath := os.Args[2]
		err := GenerateLibrary(musicPath)
		if err != nil {
			fmt.Println(err)
			return
		}
		return
	}

	var options = flag.NewFlagSet("", flag.ExitOnError)
	var originFlag = options.String("origin", "", "the origin used in the Access-Control-Allow-Origin header")

	if cmd == "serve" {
		if len(os.Args) < 4 {
			fmt.Printf("%s: not enough arguments.\n", app)
			fmt.Println(usage)
			return
		}

		libPath := os.Args[2]
		port := os.Args[3]
		lib, err := OpenLib(libPath)
		if err != nil {
			fmt.Println(err)
			return
		}

		options.Parse(os.Args[4:])
		ServeAPI(lib, port, true, *originFlag)
		return
	}

	fmt.Printf("%s: %s is an invalid command.\n", app, cmd)
	fmt.Println(usage)
}
