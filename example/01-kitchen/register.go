package main

import (
	"context"
	"errors"
	"html/template"
	"net/http"
	"net/url"

	"github.com/advanderveer/ep"
	"github.com/advanderveer/ep/coding"
	"github.com/alexedwards/scs/v2"
	"github.com/gorilla/mux"
)

type Registration struct {
	Email           string `form:"email"`
	Password        string `form:"password"`
	ConfirmPassword string `form:"confirm_password"`
}

type RegisterInput struct {
	Registration *Registration `form:"registration"`
}

// Read is an optional interface that is called by bind to read values from the
// request before decoding. This is for example usefull to read url parameters
// or authorization headers. When an error is returned it is considered a BadRequest
func (i *RegisterInput) Read(r *http.Request) (err error) {
	// @TODO can we make this method easier to test by formalizing the request
	// interface while providing a mock for it.
	return
}

// Validate the input with any custom logic. Any error indicates that validation
// has failed.
func (i RegisterInput) Check() (verr error) {
	if i.Registration == nil {
		return
	}

	if i.Registration.Password == "" || i.Registration.Email == "" || i.Registration.ConfirmPassword == "" {
		return errors.New("mandatory fields missing")
	}

	if i.Registration.Password != i.Registration.ConfirmPassword {
		return errors.New("password confirm didn't equal password")
	}

	return
}

type Register struct {
	mux  *mux.Router
	sess *scs.SessionManager
}

func (e Register) Config(cfg *ep.Config) {
	cfg.SetDecodings(NewFormDecoding())
	cfg.SetEncodings(epcoding.NewHTMLEncoding(RegisterPageTmpl, ErrorPageTmpl))
	cfg.SetLanguages("nl", "en-GB")
	return
}

func (e Register) Handle(res *ep.Response, req *http.Request) {
	var in RegisterInput
	if res.Bind(&in) {
		res.Render(e.Exec(req.Context(), &in, res.Validate(in)))
	}
}

// Exec is passed the bound input and any validation errors. Validation errors
// are passed in because service (database) interaction might be necessary to
// render the page with validation errors.
func (e Register) Exec(ctx context.Context, in *RegisterInput, verr error) (err error, out RegisterOutput) {
	out = RegisterOutput{Input: in}

	if verr == nil && in.Registration != nil {
		e.sess.Put(ctx, "message", "Success!")
		out.Redirect, err = e.mux.Get("register").URL()
		return
	}

	out.Action, err = e.mux.Get("register").URL()
	if err != nil {
		return
	}

	if out.Input.Registration == nil {
		out.Input.Registration = new(Registration)
	}

	if verr != nil {
		out.Message = verr.Error()
	}

	out.Message += e.sess.PopString(ctx, "message")
	out.Message += ep.Language(ctx)
	return verr, out // invalid input, render as show
}

type RegisterOutput struct {
	Message  string
	Redirect *url.URL
	Action   *url.URL
	Input    *RegisterInput
}

// Head is called until the request header is written. Either explicitely or
// implicitely by the output being encoded. As such, this is the best place
// to set response headers based on the output
func (o RegisterOutput) Head(w http.ResponseWriter, r *http.Request) (err error) {
	// @TODO can we make head easier to unit test by better specing the req, res
	// arguments throug an interface.

	if o.Redirect != nil {
		http.Redirect(w, r, o.Redirect.String(), http.StatusMovedPermanently)
		return
	}

	// @TODO allow for returning special error that prevents any further encoding
	// that way the implementation can completely customize encoding if desired
	return
}

// RegisterPageTmpl defines how the output will be rendered
var RegisterPageTmpl = template.Must(template.New("").Parse(
	`<!doctype html>
<html lang="en">
  <body>
  	<p>Message: {{.Message}}</p>
  	<form method="post" action="{{.Action}}">
  		<input name="registration.email" type="email" autofocus required value="{{.Input.Registration.Email}}"/>
  		<input name="registration.password" type="password" required value="{{.Input.Registration.Password}}"/>
  		<input name="registration.confirm_password" type="password" required value="{{.Input.Registration.ConfirmPassword}}"/>
  		<button type="submit">register</button>
  		<button type="reset">reset</button>
  	</form>
  </body>
</html>`))
