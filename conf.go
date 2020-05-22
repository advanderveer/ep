package ep

import (
	"net/http"

	"github.com/advanderveer/ep/coding"
)

// Hook gets triggered before the first byte is written to the response
type Hook func(out Output, w http.ResponseWriter, r *http.Request) error

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
	hooks []Hook
	errh  func(isClient bool, err error) Output
}

// New inits an empty configuration
func New() *Conf { return &Conf{} }

// Copy will duplicate the configuration
func (c Conf) Copy() (cc *Conf) {
	cc = &Conf{
		hooks: make([]Hook, len(c.hooks)),
		langs: make([]string, len(c.langs)),
		encs:  make([]epcoding.Encoding, len(c.encs)),
		decs:  make([]epcoding.Decoding, len(c.decs)),
		val:   c.val,
		qdec:  c.qdec,
		errh:  c.errh,
	}

	copy(cc.hooks, c.hooks)
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

// Handler will copy the configuration and make the endpoint as a handler
func (c Conf) Handler(ep Endpoint) *Handler {
	return &Handler{c.Copy(), ep}
}

// HandlerFunc will copy the configuration the func as a handler
func (c Conf) HandlerFunc(epf EndpointFunc) *Handler {
	return c.Handler(epf)
}

// Hooks returns any configured hooks
func (c *Conf) WithHook(hooks ...Hook) *Conf { c.hooks = append(c.hooks, hooks...); return c }
func (c Conf) Hooks() []Hook                 { return c.hooks }

// OnErrorRender determines how the response rendering handles error
func (c *Conf) SetOnErrorRender(h func(isClient bool, err error) Output) *Conf { c.errh = h; return c }
func (c Conf) OnErrorRender() func(isClient bool, err error) Output {
	return c.errh
}

type ConfReader interface {
	Hooks() []Hook
	Encodings() []epcoding.Encoding
	Decodings() []epcoding.Decoding
	Languages() []string
	Validator() Validator
	QueryDecoder() epcoding.URLValuesDecoder
	OnErrorRender() func(isClient bool, err error) Output
}
