package ep

import (
	"github.com/advanderveer/ep/coding"
)

// Validator may be implemented and configured to allow automatic validation
// of all endpoint inputs.
type Validator interface {
	Validate(v interface{}) error
}

// Conf builds endpoint configuration
type Conf struct {
	langs []string
	encs  []epcoding.Encoding
	decs  []epcoding.Decoding
	val   Validator
	qdec  epcoding.URLValuesDecoder

	serverErrFactory func(err error) Output
	clientErrFactory func(err error) Output
}

// New inits an empty configuration
func New() *Conf { return &Conf{} }

// Copy will duplicate the configuration
func (c Conf) Copy() (cc *Conf) {
	cc = &Conf{
		langs: make([]string, len(c.langs)),
		encs:  make([]epcoding.Encoding, len(c.encs)),
		decs:  make([]epcoding.Decoding, len(c.decs)),
		val:   c.val,
		qdec:  c.qdec,

		serverErrFactory: c.serverErrFactory,
		clientErrFactory: c.clientErrFactory,
	}

	copy(cc.langs, c.langs)
	copy(cc.encs, c.encs)
	copy(cc.decs, c.decs)
	return
}

// Encodings returns the configured output encodings
func (c Conf) Encodings() []epcoding.Encoding                { return c.encs }
func (c *Conf) SetEncodings(encs ...epcoding.Encoding) *Conf { c.encs = encs; return c }
func (c *Conf) WithEncoding(encs ...epcoding.Encoding) *Conf {
	c.encs = append(c.encs, encs...)
	return c
}

// Decodings returns the configured input decodings
func (c Conf) Decodings() []epcoding.Decoding                { return c.decs }
func (c *Conf) SetDecodings(decs ...epcoding.Decoding) *Conf { c.decs = decs; return c }
func (c *Conf) WithDecoding(decs ...epcoding.Decoding) *Conf {
	c.decs = append(c.decs, decs...)
	return c
}

// Languages return the supported languages for this endpoint
func (c Conf) Languages() []string                 { return c.langs }
func (c *Conf) SetLanguages(langs ...string) *Conf { c.langs = langs; return c }
func (c *Conf) WithLanguage(langs ...string) *Conf { c.langs = append(c.langs, langs...); return c }

// Validator returns the configured input validator
func (c Conf) Validator() Validator            { return c.val }
func (c *Conf) SetValidator(v Validator) *Conf { c.val = v; return c }

// QueryDecoder configures the query to be decoded into the input struct
func (c Conf) QueryDecoder() epcoding.URLValuesDecoder { return c.qdec }
func (c *Conf) SetQueryDecoder(d epcoding.URLValuesDecoder) *Conf {
	c.qdec = d
	return c
}

// SetClientErrFactory configures how client error outputs are created
func (r *Conf) SetClientErrFactory(f func(err error) Output) *Conf {
	r.clientErrFactory = f
	return r
}

// ClientErrFactory returns the current client error factory
func (r Conf) ClientErrFactory() func(err error) Output { return r.clientErrFactory }

// SetServerErrFactory configures a factory for the creation of server error outputs
func (r *Conf) SetServerErrFactory(f func(err error) Output) *Conf {
	r.serverErrFactory = f
	return r
}

// ServerErrFactory returns the configured factory for server errors
func (r Conf) ServerErrFactory() func(err error) Output { return r.serverErrFactory }

// Handler will copy the configuration and make the endpoint as a handler
func (c Conf) Handler(ep Endpoint) *Handler {
	return &Handler{c.Copy(), ep}
}

// HandlerFunc will copy the configuration the func as a handler
func (c Conf) HandlerFunc(epf EndpointFunc) *Handler {
	return c.Handler(epf)
}

type ConfReader interface {
	Encodings() []epcoding.Encoding
	Decodings() []epcoding.Decoding
	Languages() []string
	Validator() Validator
	QueryDecoder() epcoding.URLValuesDecoder
	ClientErrFactory() func(err error) Output
	ServerErrFactory() func(err error) Output
}
