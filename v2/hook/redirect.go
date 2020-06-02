package hook

import "net/http"

// Redirect is a response hook will check if the output implements a Redirect
// method to return a non-empty string as a location to redirect to. If the
// output also implements status it is called to determine the status code to
// use for redirection.
func Redirect(w http.ResponseWriter, r *http.Request, out interface{}) {
	status := http.StatusSeeOther
	if outt, ok := out.(statusOutput); ok {
		status = outt.Status()
	}

	outt, ok := out.(interface{ Redirect() string })
	if !ok {
		return
	}

	dest := outt.Redirect()
	if dest == "" {
		return
	}

	http.Redirect(w, r, dest, status)
}
