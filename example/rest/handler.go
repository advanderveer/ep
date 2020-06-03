package rest

import (
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/advanderveer/ep"
	"github.com/advanderveer/ep/coding"
	"github.com/advanderveer/ep/hook"
)

type handler struct {
	db  map[string]map[string]string
	app *ep.App
	sync.Mutex
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/idea":
		switch r.Method {
		case http.MethodPost:
			h.app.Handle(h.CreateIdea).ServeHTTP(w, r)
		case http.MethodGet:
			h.app.Handle(h.ListIdeas).ServeHTTP(w, r)
		default:
			h.app.Handle(h.MethodNotAllowed).ServeHTTP(w, r)
		}
	default:
		h.app.Handle(func(w ep.ResponseWriter, r *http.Request) {
			w.Render(h.NotFound())
		}).ServeHTTP(w, r)
	}
}

func New() http.Handler {
	logs := log.New(os.Stderr, "", 0)

	return &handler{app: ep.New(
		ep.ErrorHook(hook.NewStandardError(logs)),
		ep.RequestDecoding(coding.JSON{}),
		ep.ResponseEncoding(coding.JSON{}),
		ep.ResponseHook(hook.Status),
		ep.ResponseHook(hook.Head),
		ep.RequestHook(hook.Read),
	), db: map[string]map[string]string{
		"existing": {"name": "existing"},
	}}
}
