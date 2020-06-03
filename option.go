package ep

import (
	"net/http"

	"github.com/advanderveer/ep/epcoding"
)

// Option configures the codec
type Option interface{ apply(c *Codec) }

type options []Option

func (o options) apply(c *Codec) {
	for _, opt := range o {
		if opt == nil {
			continue
		}

		opt.apply(c)
	}
}

func Options(opts ...Option) Option {
	return options(opts)
}

// ResponseEncoding option adds an additional supported response encoding for
// this Codec
func ResponseEncoding(enc epcoding.Encoding) Option {
	return responseEncoding{enc}
}

type responseEncoding struct{ epcoding.Encoding }

func (o responseEncoding) apply(c *Codec) {
	c.encodings = append(c.encodings, o.Encoding)
}

// RequestDecoding option adds an additional supported request decoding for this
// Codec
func RequestDecoding(dec epcoding.Decoding) Option {
	return requestDecoding{dec}
}

type requestDecoding struct{ epcoding.Decoding }

func (o requestDecoding) apply(c *Codec) {
	c.decodings = append(c.decodings, o.Decoding)
}

// ResponseHook option provides arbitrary modification to the response header just
// before is send for the response. It is provided with the output that is
// currently rendered but if the response is called without using the render
// method this argument might be nil.
type ResponseHook func(w http.ResponseWriter, r *http.Request, out interface{})

func (o ResponseHook) apply(c *Codec) {
	c.resHooks = append(c.resHooks, o)
}

// RequestHook option allows reading arbitrary properties on the request to
// decorate the input value before the request body is decoded in the bind.
//
// If an error is returned the bind call will return false and the error will
// rendered.
type RequestHook func(r *http.Request, in interface{}) error

func (o RequestHook) apply(c *Codec) {
	c.reqHooks = append(c.reqHooks, o)
}

// ErrorHook can be provided as an option to be called whenever an error is
// about to be rendered. The error can be logged or an output type can be
// returend to customize how the error will be turned into a response.
//
// Multiple error hooks can be configured and the first that returns a non-nil
// output takes precedence.
//
// The returned output will be subjected to the same hooks that any other
// output would be.
type ErrorHook func(err error) (out interface{})

func (o ErrorHook) apply(c *Codec) {
	c.errHooks = append(c.errHooks, o)
}
