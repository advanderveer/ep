package rest

import "net/http"

type MethodNotAllowedOutput struct {
	Message string `json:"message"`
}

func (h *handler) MethodNotAllowed() (out MethodNotAllowedOutput) {
	out.Message = http.StatusText(http.StatusMethodNotAllowed)
	return
}

func (_ MethodNotAllowedOutput) Status() int {
	return http.StatusMethodNotAllowed
}
