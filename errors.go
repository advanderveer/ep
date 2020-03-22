package ep

import (
	"errors"
	"fmt"
	"net/http"
)

// serverErrOutput is the output that is returned by default when the response
// gets into the server error state.
type serverErrOutput struct {
	ErrorMessage string `json:"ErrorMessage"`
}

func (out serverErrOutput) Template() string { return "error" }
func (out serverErrOutput) Head(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusInternalServerError)
	return nil
}

// clientErrOutput is the output that is returned by default when the response
// gets into the client error state
type clientErrOutput struct {
	ErrorMessage string `json:"ErrorMessage"`
}

func (out clientErrOutput) Template() string { return "error" }
func (out clientErrOutput) Head(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusBadRequest)
	return nil
}

// appErrOutput is the default output that is retured when the application returns its
// own error. It allows for rendering a specific status code.
type appErrOutput struct {
	code         int
	ErrorMessage string `json:"ErrorMessage"`
}

func (out appErrOutput) Template() string { return "error" }
func (out appErrOutput) Head(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(out.code)
	return nil
}

// AppError is an application specific error that can be returned by application
// code to trigger the rendering of error output with a specific status code.
type AppError struct {
	Code int
	Err  error
}

// Error implements the error interface
func (e AppError) Error() string { return e.Err.Error() }

// Unwrap allows this error to be unwrapped
func (e AppError) Unwrap() error { return e.Err }

// Error will create an AppError with at most 1 wrapped error. If the no error
// is provided the resulting output will show the default http status message
// for the provided code.
func Error(code int, err ...error) *AppError {
	if len(err) > 1 {
		panic("ep: only takes 0 or 1 errors")
	} else if len(err) == 1 {
		return &AppError{code, err[0]}
	} else {
		return &AppError{code, errors.New(http.StatusText(code))}
	}
}

// Errorf will create a new AppError that wraps a new formatted error message.
// The f string may contain %w to wrap an error.
func Errorf(code int, f string, args ...interface{}) *AppError {
	return &AppError{code, fmt.Errorf(f, args...)}
}
