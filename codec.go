package ep

import (
	"net/http"
	"reflect"

	"github.com/advanderveer/ep/epcoding"
)

// Codec provides http.Handlers that automatically decode requests and encode
// responses based on input and output structs
type Codec struct {
	resHooks []ResponseHook
	reqHooks []RequestHook
	errHooks []ErrorHook

	decodings []epcoding.Decoding
	encodings []epcoding.Encoding
}

// New initiates a new ep Codec
func New(opts ...Option) (c *Codec) {
	c = &Codec{}
	Options(opts...).apply(c)
	return
}

// Handle will initiate an http handler that handles request according to
// the Codec configuration.
func (c *Codec) Handle(f interface{}) http.Handler {
	switch ft := f.(type) {
	case func(ResponseWriter, *http.Request):
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			res := NewResponse(w, r, c.reqHooks, c.resHooks, c.errHooks, c.decodings, c.encodings)
			defer res.Recover()
			ft(res, r)
		})
	default:
		clb, err := newCallable(f)
		if err != nil {
			panic("ep: failed to turn argument into a handler: " + err.Error())
		}

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			res := NewResponse(w, r, c.reqHooks, c.resHooks, c.errHooks, c.decodings, c.encodings)
			defer res.Recover()

			ok := true
			inv := clb.Input()
			if (inv != reflect.Value{}) {
				ok = res.Bind(inv.Interface())
			}

			if ok {
				res.Render(clb.Call(clb.Args(r, inv))...)
			}
		})
	}
}
