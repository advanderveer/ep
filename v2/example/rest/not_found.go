package rest

import "net/http"

type NotFoundOutput struct {
	Message string `json:"message"`
}

func (h *handler) NotFound() (out NotFoundOutput) {
	out.Message = http.StatusText(http.StatusNotFound)
	return
}

func (_ NotFoundOutput) Status() int { return http.StatusNotFound }
