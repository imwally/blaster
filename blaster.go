package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
)

const libName = "blaster.json"
const usage = `Usage: blaster [command] [options]
    
Commands:
    generate PATH        Generate a new library from a music path.
    serve PATH           Serve the blaster library as an HTTP API sever.

Options:
    -port                The port the HTTP server should bind to. Default is 8080.
    -origin              The origin used in the Access-Control-Allow-Origin header.

Examples:

    Create a new library containing all the tracks found in your
    ~/Music directory. A new file named blaster.json will be created
    in the same directory.
   
        $ generate ~/Music           

    Serve an HTTP API server on port 8081.
   
        $ serve ~/Music/blaser.json -port 8081 -origin "http://127.0.0.1:8081"
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

// Gen is the sub command that calls GenerateLibrary.
func Gen(args []string) {
	err := CheckArgs(args)
	if err != nil {
		fmt.Println(err)
		return
	}

	musicPath := os.Args[2]
	err = GenerateLibrary(musicPath)
	if err != nil {
		fmt.Println(err)
		return
	}
}

// Serve is the sub command that calls ServeAPI.
func Serve(args []string) {
	err := CheckArgs(args)
	if err != nil {
		fmt.Println(err)
		return
	}

	libPath := os.Args[2]
	lib, err := OpenLib(libPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	var options = flag.NewFlagSet("", flag.ExitOnError)
	var originFlag = options.String("origin", "", "the origin used in the Access-Control-Allow-Origin header")
	var portFlag = options.String("port", "8080", "the port to bind the HTTP server to")
	options.Parse(os.Args[3:])
	ServeAPI(lib, *portFlag, true, *originFlag)
}

// Help is the sub command that prints the usage text.
func Help(args []string) {
	fmt.Println(usage)
}

// CheckArgs is a helper function that checks the length of os.Args.
func CheckArgs(args []string) error {
	if len(os.Args) < 3 {
		return errors.New("blaster: not enough arguments.")
	}

	return nil
}

func main() {
	app := os.Args[0]
	if len(os.Args) < 2 {
		fmt.Fprint(os.Stdout, "blaster: no command given.\n")
		return
	}

	cmds := map[string]func([]string){
		"help":     Help,
		"generate": Gen,
		"serve":    Serve,
	}

	cmd, ok := cmds[os.Args[1]]
	if !ok {
		fmt.Printf("%s: %s is an invalid command.\n", app, cmd)
		fmt.Println(usage)
		return
	}

	cmd(os.Args)
}
