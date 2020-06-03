package web

import (
	"log"
	"net/http"
	"os"

	"github.com/advanderveer/ep"
	"github.com/advanderveer/ep/epcoding"
	"github.com/advanderveer/ep/ephook"
)

type handler struct{ app *ep.App }

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/register":
		h.app.Handle(h.Register).ServeHTTP(w, r)
	default:
		h.app.Handle(func(w ep.ResponseWriter, r *http.Request) {
			w.Render(h.NotFound())
		}).ServeHTTP(w, r)
	}
}

func New() http.Handler {
	logs := log.New(os.Stderr, "", 0)
	h := &handler{ep.New(
		ep.RequestHook(ephook.Read),
		ep.ResponseEncoding(epcoding.NewHTML(nil)),
		ep.ResponseHook(ephook.Redirect),
		ep.ResponseHook(ephook.Status),
		ep.ErrorHook(ephook.NewStandardError(logs)),
	)}

	return h
}
