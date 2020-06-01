package rest

import (
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/advanderveer/ep/v2"
	"github.com/advanderveer/ep/v2/coding"
	"github.com/advanderveer/ep/v2/hook"
)

type handler struct {
	db  map[string]interface{}
	app *ep.App
	sync.Mutex
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/idea":
		switch r.Method {
		case http.MethodPost:
			h.app.Handle(func(w ep.ResponseWriter, r *http.Request) {
				var in CreateIdeaInput
				if w.Bind(&in) {
					w.Render(h.CreateIdea(r.Context(), in))
				}
			}).ServeHTTP(w, r)
		case http.MethodGet:
			h.app.Handle(func(w ep.ResponseWriter, r *http.Request) {
				var in ListIdeasInput
				if w.Bind(&in) {
					w.Render(h.ListIdeas(r.Context(), in))
				}
			}).ServeHTTP(w, r)
		default:
			h.app.Handle(func(w ep.ResponseWriter, r *http.Request) {
				w.Render(h.MethodNotAllowed(), nil)
			}).ServeHTTP(w, r)
		}
	default:
		h.app.Handle(func(w ep.ResponseWriter, r *http.Request) {
			w.Render(h.NotFound(), nil)
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
	), db: make(map[string]interface{})}
}
