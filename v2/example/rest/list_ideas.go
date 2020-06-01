package rest

import (
	"context"
)

type (
	ListIdeasInput  struct{}
	ListIdeasOutput []map[string]string
)

func (h *handler) ListIdeas(
	ctx context.Context,
	in ListIdeasInput,
) (out ListIdeasOutput, err error) {
	return
}

func (out ListIdeasOutput) Empty() bool { return out == nil || len(out) < 1 }
