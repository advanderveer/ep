package ep

import (
	"net/http"
)

// InvalidInputError can be returned to render to write a response that indicates
// that the server understood the request but didn't accept the parameters
type InvalidInputError struct{ e error }

// InvalidInput creates a error that can be returned to indicate that the client
// needs to check its parameters
func InvalidInput(e error) *InvalidInputError {
	return &InvalidInputError{e}
}

// Error implements error interface
func (e *InvalidInputError) Error() string { return e.e.Error() }

// Unwrap returns the wrapped error
func (e *InvalidInputError) Unwrap() error { return e.e }

// serverErrOutput is the output that is returned by default when the response
// gets into the server error state.
type serverErrOutput struct{ ErrorMessage string }

func (out serverErrOutput) Template() string { return "error" }
func (out serverErrOutput) Head(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(http.StatusInternalServerError)
	return nil
}

// clientErrOutput is the output that is returned by default when the response
// gets into the client error state
type clientErrOutput struct{ ErrorMessage string }

func (out clientErrOutput) Template() string { return "error" }
func (out clientErrOutput) Head(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(http.StatusBadRequest)
	return nil
}

// invalidErrOutput is the output that is returned by default when the response
// gets into the client error state
type invalidErrOutput struct{ ErrorMessage string }

func (out invalidErrOutput) Template() string { return "error" }
func (out invalidErrOutput) Head(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(422)
	return nil
}
