package hook

import "net/http"

// Status is a response hook will assert if the out interface has a status
// method and if it has, write the header with that status
func Status(w http.ResponseWriter, r *http.Request, out interface{}) {
	if outt, ok := out.(interface{ Status() int }); ok {
		w.WriteHeader(outt.Status())
	}
}
