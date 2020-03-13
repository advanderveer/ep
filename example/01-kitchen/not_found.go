package main

import (
	"html/template"
	"net/http"
	"net/url"

	"github.com/advanderveer/ep"
	"github.com/advanderveer/ep/coding"
)

// NotFound page shows how to render a plain html template
type NotFound struct{}

func (e NotFound) Config() *ep.Config {
	cfg := &ep.Config{}
	cfg.Encoders(epcoding.NewHTMLEncoding(NotFoundPageTmpl, ErrorPageTmpl))
	return cfg
}

func (e NotFound) Handle(res *ep.Response, req *http.Request) {
	res.Render(nil, NotFoundPage{req.URL})
	return
}

// NotFoundPage holds data for redering the not found page
type NotFoundPage struct{ Location *url.URL }

func (o NotFoundPage) Head(http.ResponseWriter, *http.Request) (err error) { return }

// NotFoundPageTmpl defines how the output will be rendered
var NotFoundPageTmpl = template.Must(template.New("").Parse(`oops, couldn't find {{.Location}}`))
var ErrorPageTmpl = template.Must(template.New("").Parse(`error: {{.ErrorMessage}}`))
