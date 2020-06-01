package rest

import (
	"context"
	"errors"
	"net/http"
)

type (
	CreateIdeaInput struct {
		Name string `json:"name"`
	}

	CreateIdeaOutput struct{}
)

func (h *handler) CreateIdea(
	ctx context.Context,
	in CreateIdeaInput,
) (out *CreateIdeaOutput, err error) {
	err = in.Validate()
	if err != nil {
		return nil, ErrorOutput{422, err.Error()}
	}

	h.Lock()
	defer h.Unlock()
	if _, ok := h.db[in.Name]; ok {
		return nil, ErrorOutput{409, "Idea already exist"}
	}

	h.db[in.Name] = struct{}{}

	// @TODO return 201 created
	return
}

func (_ CreateIdeaOutput) Head(h http.Header) {
	h.Set("Location", "/ideas")
}

func (_ CreateIdeaOutput) Status() int { return 201 }

func (in CreateIdeaInput) Validate() error {
	if in.Name == "" {
		return errors.New("Name is empty")
	}

	return nil
}
