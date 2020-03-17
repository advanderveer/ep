package main

import (
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

	r := mux.NewRouter()
	r.Use(smgr.LoadAndSave)

	r.Path("/register").Methods("GET", "POST").Name("register").Handler(ep.New().
		WithLanguage("nl", "en-GB").
		WithDecoding(epcoding.NewFormDecoding(urld)).
		WithEncoding(epcoding.NewHTMLEncoding(RegisterPageTmpl, ErrorPageTmpl)).
		Handler(Register{r, smgr}))

	r.Path("/hello").Handler(ep.New().Handler(Hello{}))

	r.Path("/kitchen").Handler(ep.New().
		SetQueryDecoder(urld).
		WithEncoding(epcoding.NewJSONEncoding()).
		HandlerFunc(HandleKitchen))

	r.PathPrefix("/").Handler(ep.New().
		WithEncoding(epcoding.NewHTMLEncoding(NotFoundPageTmpl, ErrorPageTmpl)).
		Handler(NotFound{}))

	panic(http.ListenAndServe(":10010", r))
}
