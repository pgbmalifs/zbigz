### Dependencies

* Go 1.7+
* `go-bindata` (`go get github.com/jteeuwen/go-bindata`)

### Architecture

* `cloudtorrent/module/` contains ancillary features. Modules can optionally implement the:
    * `config.Configurable` interface - allows `App` to get and set a JSON marshallable structure.
    * `fs.FS` interface - an `fs.Node` tree is synced with the front-end and can be downloaded and/or uploaded from/to.
    * `stream.Transformer` interface - allow the implementation of streaming transforms of file transfers from any, to any source.
* `cloudtorrent/static/` contains static front-end files which are embedded into `files.go` using `go generate ./...`.
* `cloudtorrent/` contains the root application class (`cloudtorrent.App`)
    * Contains the root HTTP handler which has routes for:
        * State sync using velox
        * Configuration API
        * Static file serving
