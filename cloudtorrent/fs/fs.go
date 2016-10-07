package fs

import "io"

type ReadableFS interface {
	ReadStream(path string) (io.ReadCloser, error)
}

type WritableFS interface {
	WriteStream(path string) (io.WriteCloser, error)
}

type FS interface {
	Sync(chan Node) error
}
