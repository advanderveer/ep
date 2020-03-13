package ep

import (
	"bufio"
	"io"
	"net/http"
)

// Reader is a buffered reader that keeps track of the number of bytes that
// have been read.
type Reader struct {
	io.Closer
	*bufio.Reader
}

// NewReader turns an unbuffered reader into a buffered read that keeps track
// of reading progress. The buffer size is 512 bytes to allow MIME sniffing of
// the readers content
func NewReader(r io.ReadCloser) *Reader {
	return &Reader{r, bufio.NewReaderSize(r, 512)}
}

// Sniff will use peek into the reader and tries to determine the content type
func (r *Reader) Sniff() (ct string) {
	b, _ := r.Reader.Peek(512) // in case of error, just work with whats available

	ct = http.DetectContentType(b)
	if ct == "text/plain; charset=utf-8" && len(b) > 0 {

		// in our usecase of detecting request bodies we can be a bit more
		// liberal and assume that the client tries to send either JSON or XML
		switch {
		case b[0] == '{':
			fallthrough
		case b[0] == '"':
			fallthrough
		case b[0] == '[':
			return "application/json; charset=utf-8"
		case b[0] == '<':
			return "text/xml; charset=utf-8"
		}
	}

	return
}
