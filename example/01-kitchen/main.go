package main

import (
	"html/template"
	"net/http"

	"github.com/advanderveer/ep"
	"github.com/advanderveer/ep/coding"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"github.com/gorilla/mux"
)

func main() {
	smgr := scs.New()
	urld := form.NewDecoder()

	view := template.New("view")
	view.New("register").Parse(RegisterPageTmpl)
	view.New("not_found").Parse(NotFoundPageTmpl)
	view.New("error").Parse(`error: {{.ErrorMessage}}`)

	r := mux.NewRouter()
	r.Use(smgr.LoadAndSave)

	r.Path("/register").Methods("GET", "POST").Name("register").Handler(ep.New().
		WithLanguage("nl", "en-GB").
		WithDecoding(epcoding.NewFormDecoding(urld)).
		WithEncoding(epcoding.NewHTMLEncoding(view)).
		Handler(Register{r, smgr}))

	r.Path("/hello").Handler(ep.New().Handler(Hello{}))

	r.Path("/kitchen").Handler(ep.New().
		SetQueryDecoder(urld).
		WithEncoding(epcoding.NewJSONEncoding()).
		HandlerFunc(HandleKitchen))

	r.PathPrefix("/").Handler(ep.New().
		WithEncoding(epcoding.NewHTMLEncoding(view)).
		Handler(NotFound{}))

	panic(http.ListenAndServe(":10010", r))
}
