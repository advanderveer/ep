package ep

import (
	"net/http"

	"github.com/advanderveer/ep/v2/coding"
)

// App holds application wide configuration for binding inputs and rendering
// outputs.
type App struct {
	resHooks []ResponseHook
	reqHooks []RequestHook
	errHooks []ErrorHook

	decodings []coding.Decoding
	encodings []coding.Encoding
}

// New initiates a new ep application
func New(opts ...Option) (app *App) {
	app = &App{}
	Options(opts...).apply(app)
	return
}

// Handle will initiate an http handler that handles request according to
// the application configuration.
func (a *App) Handle(h func(ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		res := NewResponse(w, r, a.reqHooks, a.resHooks, a.errHooks)
		if res.Negotiate(a.decodings, a.encodings) {
			h(res, r)
		}
	})
}
