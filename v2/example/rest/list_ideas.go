package rest

import (
	"context"
	"net/http"
)

type (
	ListIdeasInput  struct{ NameFilter string }
	ListIdeasOutput []map[string]string
)

func (h *handler) ListIdeas(
	ctx context.Context,
	in ListIdeasInput,
) (out ListIdeasOutput, err error) {

	h.Lock()
	defer h.Unlock()
	out = make(ListIdeasOutput, 0, len(h.db))
	for _, idea := range h.db {
		if in.NameFilter != "" && idea["name"] != in.NameFilter {
			continue
		}

		out = append(out, idea)
	}

	return
}

func (in *ListIdeasInput) Read(r *http.Request) {
	in.NameFilter = r.FormValue("name")
}

func (out ListIdeasOutput) Empty() bool {
	return out == nil || len(out) < 1
}
