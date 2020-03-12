package epcoding

import (
	"html/template"
	"io"
)

type HTMLEncoding struct {
	vt *template.Template // regular value template
	et *template.Template // error value template
}

func NewHTMLEncoding(vt *template.Template, et *template.Template) HTMLEncoding {
	return HTMLEncoding{vt, et}
}

func (e HTMLEncoding) Produces() string            { return "text/html" }
func (e HTMLEncoding) Encoder(w io.Writer) Encoder { return HTMLEncoder{e.vt, e.et, w} }

type HTMLEncoder struct {
	vt *template.Template
	et *template.Template
	w  io.Writer
}

func (e HTMLEncoder) Encode(v interface{}) error {
	if _, ok := v.(ErrorEncode); ok {
		return e.et.Execute(e.w, v)
	} else {
		return e.vt.Execute(e.w, v)
	}
}
