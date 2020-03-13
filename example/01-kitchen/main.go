package main

import (
	"net/http"

	"github.com/advanderveer/ep"
	"github.com/alexedwards/scs/v2"
	"github.com/gorilla/mux"
)

func main() {
	smgr := scs.New()

	r := mux.NewRouter()
	r.Path("/register").Methods("GET", "POST").Handler(ep.Handler(Register{r, smgr})).Name("register")
	r.Path("/hello").Handler(ep.Handler(Hello{}))
	r.Path("/kitchen").Handler(ep.Handler(Kitchen{}))
	r.PathPrefix("/").Handler(ep.Handler(NotFound{}))

	panic(http.ListenAndServe(":10010", smgr.LoadAndSave(r)))
}
