package epcoding

import (
	"io"
)

// Template is implemented by text/template.Template or an html/template.Template
type Template interface {
	ExecuteTemplate(w io.Writer, name string, v interface{}) error
	Execute(w io.Writer, v interface{}) error
}

// TemplatedOutput can be implemented by outputs to render a named sub-template
// in the namespace provided in the TemplateEncoding constructor
type TemplatedOutput interface {
	Template() string
}

// TemplateEncoding will write html by executing html/template.Templates
type TemplateEncoding struct {
	view     Template
	produces string
}

// NewTemplateEncoding provides the ability to render the provided template or any
// named sub-template in the namespace
func NewTemplateEncoding(view Template) *TemplateEncoding {
	return &TemplateEncoding{view, "text/html"}
}

// Produces returns what conten type this encoding produces
func (e TemplateEncoding) Produces() string { return e.produces }

// SetProduces will onvewrite what content type the encoding produces
func (e *TemplateEncoding) SetProduces(p string) *TemplateEncoding { e.produces = p; return e }

// Encoder return an encoder
func (e TemplateEncoding) Encoder(w io.Writer) Encoder { return TemplateEncoder{e.view, w} }

// TemplateEncoder allows for actual encoding
type TemplateEncoder struct {
	view Template
	w    io.Writer
}

// Encode the provided value
func (e TemplateEncoder) Encode(v interface{}) error {
	if tv, ok := v.(TemplatedOutput); ok {
		return e.view.ExecuteTemplate(e.w, tv.Template(), v)
	}

	return e.view.Execute(e.w, v)
}
