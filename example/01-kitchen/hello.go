package main

import (
	"fmt"
	"github.com/advanderveer/ep"
	"net/http"
)

// Hello shows the simplest endpoint possible
type Hello struct{}

func (e Hello) Handle(res *ep.Response, req *http.Request) {
	fmt.Fprintf(res, "hello, %s\n", req.RemoteAddr)
	return
}
