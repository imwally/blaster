# blaster

Blaster is a lightweight application that indexes your music
collection and exposes an API to retrieve information from
it. Meta-data is organized and stored in a json file which is accessed
through a RESTful API. The end points are described below.

## Usage
```
Usage: blaster [command] [options]
    
Commands:
    generate PATH              Generate a new library from a music path.
    serve PATH PORT            Serve HTTP API server on PORT.

Options:
    -origin                    The origin used in the Access-Control-Allow-Origin header.

Examples:

    Create a new blaster.json library in your ~/Music directory.
   
        $ generate ~/Music

    Serve an HTTP API server on port 8080.
   
        $ serve ~/Music/blaser.json 8080 -origin "http://127.0.0.1:8081"
```

## HTTP End Points

*/artists*

All artists.

*/albums*

All albums.

*/albums?artist={artist}*

Albums by an artists.

*/tracks?album={album}*

Tracks from an album.

*/artwork?track={path}*

Album artwork for a specific track.
