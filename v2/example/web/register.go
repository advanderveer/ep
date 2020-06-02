package web

import (
	"context"
	"html/template"
	"net/http"
	"net/url"
)

type (
	Registration struct {
		Email    string
		Password string
	}

	RegisterInput struct{ *Registration }

	RegisterOutput struct {
		RedirectTo string
		Validation url.Values
	}
)

func (h *handler) Register(
	ctx context.Context,
	in RegisterInput,
) (out RegisterOutput, err error) {
	if in.Registration == nil {
		return // just show empty form
	}

	if out.Validation = in.Validate(); len(out.Validation) > 0 {
		return // show form with validation feedback
	}

	// success
	out.RedirectTo = "/"
	return
}

func (_ RegisterInput) Empty() bool { return true }

func (in *RegisterInput) Read(r *http.Request) {
	if r.Method != http.MethodPost {
		return
	}

	in.Registration = new(Registration)
	in.Email = r.FormValue("email")
	in.Password = r.FormValue("password")
}

func (in RegisterInput) Validate() (v url.Values) {
	v = make(url.Values)
	if in.Registration.Email == "" {
		v.Add("email", "Email is required")
	}

	if in.Registration.Password == "" {
		v.Add("password", "Password is required")
	}

	return
}

func (out RegisterOutput) Empty() bool {
	if out.RedirectTo != "" {
		return true
	}

	return false
}

func (out RegisterOutput) Status() int {
	if len(out.Validation) > 0 {
		return 422
	}

	if out.RedirectTo != "" {
		return 301
	}

	return 200
}

func (out RegisterOutput) Redirect() string {
	return out.RedirectTo
}

var RegisterTmpl = template.Must(
	template.New("").Parse(`form{{ if .Validation }}invalid{{end}}`))

func (_ RegisterOutput) Template() *template.Template {
	return RegisterTmpl
}
