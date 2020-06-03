package rest

import (
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/advanderveer/ep"
	"github.com/advanderveer/ep/epcoding"
	"github.com/advanderveer/ep/ephook"
)

type handler struct {
	db    map[string]map[string]string
	codec *ep.Codec
	sync.Mutex
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/idea":
		switch r.Method {
		case http.MethodPost:
			h.codec.Handle(h.CreateIdea).ServeHTTP(w, r)
		case http.MethodGet:
			h.codec.Handle(h.ListIdeas).ServeHTTP(w, r)
		default:
			h.codec.Handle(h.MethodNotAllowed).ServeHTTP(w, r)
		}
	default:
		h.codec.Handle(func(w ep.ResponseWriter, r *http.Request) {
			w.Render(h.NotFound())
		}).ServeHTTP(w, r)
	}
}

func New() http.Handler {
	logs := log.New(os.Stderr, "", 0)

	return &handler{codec: ep.New(
		ep.ErrorHook(ephook.NewStandardError(logs)),
		ep.RequestDecoding(epcoding.JSON{}),
		ep.ResponseEncoding(epcoding.JSON{}),
		ep.ResponseHook(ephook.Status),
		ep.ResponseHook(ephook.Head),
		ep.RequestHook(ephook.Read),
	), db: map[string]map[string]string{
		"existing": {"name": "existing"},
	}}
}
