package ep

import (
	"bufio"
	"io"
)

// PeekReadCloser allows for reading, closing and peeking
type PeekReadCloser interface {
	io.ReadCloser
	Peek(n int) ([]byte, error)
}

// Buffer the provided ReadCloser so that it is possible to peek into the reader
// without influencing the next read. Unlike the original ReadCloser this
// buffered variant might not error immediately when calling Read after Close.
// Instead it will error once the buffered bytes are completely read.
func Buffer(rc io.ReadCloser) PeekReadCloser {
	type peekReadCloser struct {
		io.Closer
		*bufio.Reader
	}

	return &peekReadCloser{rc, bufio.NewReader(rc)}
}
