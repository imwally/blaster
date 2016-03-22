# blaster

Blaster is a lightweight application that indexes your music
collection and exposes an API to retrieve information from
it. Meta-data is organized and stored in a json file which is accessed
through a RESTful API. The end points are described below.

## Usage
```
Usage: blaster [command] [options]
    
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

    Serve an HTTP API server on port 8080 allowing requests from
    http://127.0.0.1:8081.
   
        $ serve ~/Music/blaser.json -origin "http://127.0.0.1:8081"
```

## HTTP End Points

`/artists`

All artists.

`/albums`

All albums.

`/albums?artist={artist}`

Albums by an artist.

`/tracks?album={album}`

Tracks from an album.

`/artwork?track={path}`

Album artwork for a specific track.
