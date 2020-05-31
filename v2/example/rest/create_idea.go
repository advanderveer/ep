package rest

import (
	"context"
	"errors"
)

type (
	CreateIdeaInput  struct{ Name string }
	CreateIdeaOutput struct{}
)

func (h *handler) CreateIdea(
	ctx context.Context,
	in CreateIdeaInput,
) (out *CreateIdeaOutput, err error) {
	err = in.Validate()
	if err != nil {
		return nil, err
	}

	// @TODO validate input

	// @TODO check if the idea doesn't exist yet

	// @TODO set actual idea

	return
}

func (in CreateIdeaInput) Validate() error {
	if in.Name == "" {
		return errors.New("name is empty")
	}

	return nil
}
