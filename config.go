package ep

import (
	"github.com/advanderveer/ep/coding"

	"net/http"
)

type Validator interface {
	Validate(v interface{}) error
}

type DefaultServerError struct {
	ErrorMessage string
}

func (out DefaultServerError) Head(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(http.StatusInternalServerError)
	return nil
}

func (out DefaultServerError) IsError() {}

type DefaultClientError struct {
	ErrorMessage string
}

func (out DefaultClientError) Head(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(http.StatusBadRequest)
	return nil
}

func (out DefaultClientError) IsError() {}

type Config struct {
	encs []epcoding.Encoding
	decs []epcoding.Decoding

	onServerError func(err error) ErrorOutput
	onClientError func(err error) ErrorOutput

	validator Validator
}

func (r Config) Validator() Validator { return r.validator }

func (r *Config) SetValidator(v Validator) { r.validator = v }

// Decoders specifies what sort of content the server is willing to decode from.
// It will look at the Content-Type header to determine the ddecoder
func (r *Config) Decoders(decs ...epcoding.Decoding) *Config { r.decs = decs; return r }

// Encoders will determine to what content the server is willing to encode. It
// will attempt to satisfy the clients 'Accept' header but may fallback to
// a default encoding.
func (r *Config) Encoders(encs ...epcoding.Encoding) *Config { r.encs = encs; return r }

func (r *Config) OnClientError(f func(err error) ErrorOutput) { r.onClientError = f }
func (r *Config) OnServerError(f func(err error) ErrorOutput) { r.onServerError = f }

func (r *Config) ServerErrorOutput(err error) ErrorOutput {
	if r.onServerError == nil {
		return DefaultServerError{http.StatusText(http.StatusInternalServerError)}
	}

	return r.onServerError(err)
}
func (r *Config) ClientErrorOutput(err error) ErrorOutput {
	if r.onClientError == nil {
		return DefaultClientError{http.StatusText(http.StatusBadRequest)}
	}

	return r.onClientError(err)
}
