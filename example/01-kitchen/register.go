package main

import (
	"context"
	"errors"
	"html/template"
	"net/http"
	"net/url"

	"github.com/advanderveer/ep"
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

func (e Register) Handle(res *ep.Response, req *http.Request) {
	var in RegisterInput
	if res.Bind(&in) {
		res.Render(e.Exec(req.Context(), &in, res.Validate(in)))
	}
}

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

func (o RegisterOutput) Head(w http.ResponseWriter, r *http.Request) (err error) {
	if o.Redirect != nil {
		http.Redirect(w, r, o.Redirect.String(), http.StatusMovedPermanently)
		return
	}

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
