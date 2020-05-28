package ep

import (
	"net/http"

	"github.com/advanderveer/ep/v2/coding"
)

// Option configures the app
type Option interface{ apply(a *App) }

type options []Option

func (o options) apply(a *App) {
	for _, opt := range o {
		if opt == nil {
			continue
		}

		opt.apply(a)
	}
}

func Options(opts ...Option) Option {
	return options(opts)
}

// ResponseEncoding option adds an additional supported response encoding for
// this app
func ResponseEncoding(enc coding.Encoding) Option {
	return responseEncoding{enc}
}

type responseEncoding struct{ coding.Encoding }

func (o responseEncoding) apply(a *App) {
	a.encodings = append(a.encodings, o.Encoding)
}

// RequestDecoding option adds an additional supported request decoding for this
// app
func RequestDecoding(dec coding.Decoding) Option {
	return requestDecoding{dec}
}

type requestDecoding struct{ coding.Decoding }

func (o requestDecoding) apply(a *App) {
	a.decodings = append(a.decodings, o.Decoding)
}

// ResponseHook option provides arbitrary modification to the response header just
// before is send for the response. It is provided with the output that is
// currently rendered but if the response is called without using the render
// method this argument might be nil.
type ResponseHook func(w http.ResponseWriter, r *http.Request, out interface{})

func (o ResponseHook) apply(a *App) {
	a.resHooks = append(a.resHooks, o)
}

// RequestHook option allows reading arbitrary properties on the request to
// decorate the input value before the request body is decoded in the bind.
//
// If an error is returned the bind call will return false and the error will
// rendered.
type RequestHook func(r *http.Request, in interface{}) error

func (o RequestHook) apply(a *App) {
	a.reqHooks = append(a.reqHooks, o)
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

func (o ErrorHook) apply(a *App) {
	a.errHooks = append(a.errHooks, o)
}
