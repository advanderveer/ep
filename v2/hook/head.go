package hook

import "net/http"

// Head is response hook that checks has a Head method that wants full control
// over the response headers before the first output is encoded.
func Head(w http.ResponseWriter, r *http.Request, out interface{}) {
	switch outt := out.(type) {
	case interface {
		Head(http.Header)
	}:
		outt.Head(w.Header())
	case interface {
		Head(http.ResponseWriter)
	}:
		outt.Head(w)
	case interface {
		Head(http.ResponseWriter, *http.Request)
	}:
		outt.Head(w, r)
	}
}
