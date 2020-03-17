package epcoding

import (
	"html/template"
	"io"
)

// TemplatedOutput can be implemented by outputs to render a named sub-template
// in the namespace provided in the HTMLEncoding constructor
type TemplatedOutput interface {
	Template() string
}

// HTMLEncoding will write html by executing html/template.Templates
type HTMLEncoding struct {
	view     *template.Template
	produces string
}

// NewHTMLEncoding provides the ability to render the provided template or any
// named sub-template in the namespace
func NewHTMLEncoding(view *template.Template) *HTMLEncoding {
	return &HTMLEncoding{view, "text/html"}
}

// Produces returns what conten type this encoding produces
func (e HTMLEncoding) Produces() string { return e.produces }

// SetProduces will onvewrite what content type the encoding produces
func (e *HTMLEncoding) SetProduces(p string) *HTMLEncoding { e.produces = p; return e }

// Encoder return an encoder
func (e HTMLEncoding) Encoder(w io.Writer) Encoder { return HTMLEncoder{e.view, w} }

// HTMLEncoder allows for actual encoding
type HTMLEncoder struct {
	view *template.Template
	w    io.Writer
}

// Encode the provided value
func (e HTMLEncoder) Encode(v interface{}) error {
	if tv, ok := v.(TemplatedOutput); ok {
		return e.view.ExecuteTemplate(e.w, tv.Template(), v)
	}

	return e.view.Execute(e.w, v)
}
