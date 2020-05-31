package rest

import (
	"context"
)

type (
	ListIdeasInput  struct{}
	ListIdeasOutput struct{}
)

func (h *handler) ListIdeas(
	ctx context.Context,
	in ListIdeasInput,
) (out *ListIdeasOutput, err error) {
	return
}
