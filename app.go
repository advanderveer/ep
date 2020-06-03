package ep

import (
	"net/http"
	"reflect"

	"github.com/advanderveer/ep/epcoding"
)

// App holds application wide configuration for binding inputs and rendering
// outputs.
type App struct {
	resHooks []ResponseHook
	reqHooks []RequestHook
	errHooks []ErrorHook

	decodings []epcoding.Decoding
	encodings []epcoding.Encoding
}

// New initiates a new ep application
func New(opts ...Option) (app *App) {
	app = &App{}
	Options(opts...).apply(app)
	return
}

// Handle will initiate an http handler that handles request according to
// the application configuration.
func (a *App) Handle(f interface{}) http.Handler {
	switch ft := f.(type) {
	case func(ResponseWriter, *http.Request):
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			res := NewResponse(w, r, a.reqHooks, a.resHooks, a.errHooks, a.decodings, a.encodings)
			defer res.Recover()
			ft(res, r)
		})
	default:
		clb, err := newCallable(f)
		if err != nil {
			panic("ep: failed to turn argument into a handler: " + err.Error())
		}

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			res := NewResponse(w, r, a.reqHooks, a.resHooks, a.errHooks, a.decodings, a.encodings)
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
