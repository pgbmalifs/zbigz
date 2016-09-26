package fs

type FSMode int

const (
	R  FSMode = 1 << 0
	W         = 1 << 1
	RW        = R | W
)

type FS interface {
	Mode() FSMode
	Sync(chan Node) error
}
