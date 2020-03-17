package epcoding

import (
	"html/template"
	"io"
)

type TemplatedOutput interface {
	Template() string
}

type HTMLEncoding struct{ view *template.Template }

func NewHTMLEncoding(view *template.Template) HTMLEncoding {
	return HTMLEncoding{view}
}

func (e HTMLEncoding) Produces() string            { return "text/html" }
func (e HTMLEncoding) Encoder(w io.Writer) Encoder { return HTMLEncoder{e.view, w} }

type HTMLEncoder struct {
	view *template.Template
	w    io.Writer
}

func (e HTMLEncoder) Encode(v interface{}) error {
	if tv, ok := v.(TemplatedOutput); ok {
		return e.view.ExecuteTemplate(e.w, tv.Template(), v)
	}

	return e.view.Execute(e.w, v)
}
