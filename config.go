package ep

import (
	"github.com/advanderveer/ep/coding"
)

// Validator may be implemented and configured to allow automatic validation
// of all endpoint inputs.
type Validator interface {
	Validate(v interface{}) error
}

// Config allows endpoints to configure response formulation
type Config struct {
	encs  []epcoding.Encoding
	decs  []epcoding.Decoding
	langs []string

	serverErrFactory func(err error) ErrorOutput
	clientErrFactory func(err error) ErrorOutput

	validator Validator
}

// Validator returns the configured validator
func (r Config) Validator() Validator { return r.validator }

// SetValidator configures a validator that is run for each endpoint input
func (r *Config) SetValidator(v Validator) { r.validator = v }

// Languages returns the currently supported languages
func (r Config) Languages() []string { return r.langs }

// SetLanguages will configure what laguages are supported
func (r *Config) SetLanguages(langs ...string) { r.langs = langs }

// SetDecodings will configure the supported input decodings
func (r *Config) SetDecodings(decs ...epcoding.Decoding) *Config { r.decs = decs; return r }

// Decodings returns the configured endpoint decodings
func (r Config) Decodings() []epcoding.Decoding { return r.decs }

// SetEncodings will configure the supported output encodings
func (r *Config) SetEncodings(encs ...epcoding.Encoding) *Config { r.encs = encs; return r }

// Encodings returns the configured output encodings
func (r Config) Encodings() []epcoding.Encoding { return r.encs }

// SetClientErrFactory configures how client error outputs are created
func (r *Config) SetClientErrFactory(f func(err error) ErrorOutput) { r.clientErrFactory = f }

// ClientErrFactory returns the current client error factory
func (r Config) ClientErrFactory() func(err error) ErrorOutput { return r.clientErrFactory }

// SetServerErrFactory configures a factory for the creation of server error outputs
func (r *Config) SetServerErrFactory(f func(err error) ErrorOutput) { r.serverErrFactory = f }

// ServerErrFactory returns the configured factory for server errors
func (r Config) ServerErrFactory() func(err error) ErrorOutput { return r.serverErrFactory }
