package ep

import (
	"net/http"
)

// Hook gets triggered before the first byte is written to the response
type Hook func(out Output, w http.ResponseWriter, r *http.Request) error

// HeadHook is a hook that runs a head function on the output if it has one.
func HeadHook(out Output, w http.ResponseWriter, r *http.Request) error {
	hout, ok := out.(HeaderOutput)
	if ok {
		return hout.Head(w, r)
	}

	return nil
}
