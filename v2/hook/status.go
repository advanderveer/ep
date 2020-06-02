package hook

import "net/http"

type statusOutput interface{ Status() int }

// Status is a response hook will assert if the out interface has a status
// method and if it has, write the header with that status
func Status(w http.ResponseWriter, r *http.Request, out interface{}) {
	if outt, ok := out.(statusOutput); ok {
		w.WriteHeader(outt.Status())
	}
}
