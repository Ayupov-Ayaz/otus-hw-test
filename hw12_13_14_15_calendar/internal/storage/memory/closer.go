package memory

type Closer struct{}

func NewCloser() Closer {
	return Closer{}
}

func (c Closer) Close() error { return nil }
