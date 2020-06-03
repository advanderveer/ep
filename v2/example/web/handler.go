package web

import (
	"log"
	"net/http"
	"os"

	"github.com/advanderveer/ep/v2"
	"github.com/advanderveer/ep/v2/coding"
	"github.com/advanderveer/ep/v2/hook"
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
		ep.RequestHook(hook.Read),
		ep.ResponseEncoding(coding.NewHTML(nil)),
		ep.ResponseHook(hook.Redirect),
		ep.ResponseHook(hook.Status),
		ep.ErrorHook(hook.NewStandardError(logs)),
	)}

	return h
}
