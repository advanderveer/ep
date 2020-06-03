package coding

import (
	"errors"
	"html/template"
	"net/http"
)

var (
	// NoTemplateSpecified is returned when an output that wants to be encoded
	// doesn't specify a template
	NoTemplateSpecified = errors.New("output without method that selects template")
)

// NewHTML initializes an html encoder with the provide template(s)
func NewHTML(t *template.Template) Encoding {
	return &htmlEncoding{t}
}

type (
	htmlEncoding struct {
		t *template.Template
	}
	htmlEncoder struct {
		w http.ResponseWriter
		e *htmlEncoding
	}
)

func (e *htmlEncoding) Produces() string { return "text/html" }

func (e *htmlEncoding) Encoder(w http.ResponseWriter) Encoder {
	return &htmlEncoder{w, e}
}

func (e *htmlEncoder) Encode(v interface{}) (err error) {
	switch vt := v.(type) {
	case interface{ Template() string }:
		return e.e.t.ExecuteTemplate(e.w, vt.Template(), v)
	case interface{ Template() *template.Template }:
		return vt.Template().Execute(e.w, v)
	default:
		return NoTemplateSpecified
	}
}
