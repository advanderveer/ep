package ep

import (
	"net/http"
)

// StatusCreated can be embedded to automatically set the setatus code to 201
// when it is rendered as output. If SetLocation is called the Location header
// will also be set.
type StatusCreated struct{ location string }

func (s *StatusCreated) SetLocation(l string) { s.location = l }
func (s StatusCreated) Head(w http.ResponseWriter, r *http.Request) error {
	if s.location != "" {
		w.Header().Set("Location", s.location)
	}

	w.WriteHeader(http.StatusCreated)
	return nil
}

// StatusNoContent can be embedded to automatically set the response status
// code to 204 and prevent any body from being returned.
type StatusNoContent struct{}

func (s StatusNoContent) Head(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(http.StatusNoContent)
	return nil
}

// StatusRedirect can be embedded to render a redirect using http.Redirect. It
// will only do so if SetRedirect is called with a non-empty string.
type StatusRedirect struct {
	location string
	code     int
}

// SetRedirect will cause the response to be a redirect. By default it will
// be a '303 See Other' redirect unless the code is used set it to another
// 3xx status code.
func (s *StatusRedirect) SetRedirect(l string, code ...int) {
	if len(code) > 1 {
		panic("ep: must be 0 or 1 code provided")
	} else if len(code) == 1 {
		s.code = code[0]
	} else {
		s.code = http.StatusSeeOther
	}

	s.location = l
}

func (s StatusRedirect) Head(w http.ResponseWriter, r *http.Request) (err error) {
	if s.location != "" && s.code > 0 {
		http.Redirect(w, r, s.location, s.code)
		return
	}

	return
}
