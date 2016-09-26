package stream

type Stream struct {
}

type Transformer interface {
	ID() string
	Transform(Stream) (Stream, error)
}
