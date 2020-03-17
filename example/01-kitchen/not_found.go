package main

import (
	"net/http"
	"net/url"

	"github.com/advanderveer/ep"
)

// NotFound page shows how to render a plain html template
type NotFound struct{}

func (e NotFound) Handle(res *ep.Response, req *http.Request) {
	res.Render(nil, NotFoundPage{req.URL})
	return
}

// NotFoundPage holds data for redering the not found page
type NotFoundPage struct{ Location *url.URL }

func (o NotFoundPage) Template() string                                    { return "not_found" }
func (o NotFoundPage) Head(http.ResponseWriter, *http.Request) (err error) { return }

var NotFoundPageTmpl = `oops, couldn't find {{.Location}}`
