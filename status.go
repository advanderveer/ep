package ep

import (
	"net/http"
)

// StatusCreatedHook allow outputs to be rendered as created responses
func StatusCreatedHook(out Output, w http.ResponseWriter, r *http.Request) error {
	outt, ok := out.(interface{ statusCreated() string })
	if !ok {
		return nil
	}

	loc := outt.statusCreated()
	if loc != "" {
		w.Header().Set("Location", loc)
	}

	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusCreated)
	return nil
}

// StatusCreated can be embedded to automatically set the status code to 201
// when it is rendered as output. If SetLocation is called the Location header
// will also be set.
type StatusCreated struct{ location string }

func (s StatusCreated) statusCreated() string { return s.location }
func (s *StatusCreated) SetLocation(l string) { s.location = l }

// StatusNoContentHook allow outputs to be rendered as created responses
func StatusNoContentHook(out Output, w http.ResponseWriter, r *http.Request) error {
	_, ok := out.(interface{ statusNoContent() })
	if !ok {
		return nil
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

// StatusNoContent can be embedded to automatically set the response status
// code to 204 and prevent any body from being returned.
type StatusNoContent struct{}

func (s StatusNoContent) statusNoContent() {}

// StatusRedirectHook allow redirect outputs to be rendered correctly easily
func StatusRedirectHook(out Output, w http.ResponseWriter, r *http.Request) error {
	outt, ok := out.(interface{ statusRedirect() (string, int) })
	if !ok {
		return nil
	}

	loc, code := outt.statusRedirect()
	if loc != "" && code > 0 {
		http.Redirect(w, r, loc, code)
	}

	return nil
}

// StatusRedirect can be embedded to render a redirect using http.Redirect. It
// will only do so if SetRedirect is called with a non-empty string.
type StatusRedirect struct {
	location string
	code     int
}

func (s StatusRedirect) statusRedirect() (string, int) {
	return s.location, s.code
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
