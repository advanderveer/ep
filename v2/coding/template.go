package coding

import (
	"errors"
	"io"
	"net/http"
)

var (
	// NoTemplateSpecified is returned when an output that wants to be encoded
	// doesn't specify a template
	NoTemplateSpecified = errors.New("output without method that selects template")
)

// Template provides the the only method the encoder requires. It implemented
// by html/template.Template or text/template.Template
type Template interface {
	ExecuteTemplate(w io.Writer, name string, v interface{}) error
}

// NewTemplate initializes an template encoding
func NewTemplate(t Template) Encoding {
	return &templateEncoding{t}
}

type (
	templateEncoding struct {
		t Template
	}
	templateEncoder struct {
		w http.ResponseWriter
		e *templateEncoding
	}
)

func (e *templateEncoding) Produces() string { return "text/html" }

func (e *templateEncoding) Encoder(w http.ResponseWriter) Encoder {
	return &templateEncoder{w, e}
}

func (e *templateEncoder) Encode(v interface{}) (err error) {
	name, ok := v.(interface{ Template() string })
	if !ok {
		return NoTemplateSpecified
	}

	err = e.e.t.ExecuteTemplate(e.w, name.Template(), v)
	if err != nil {
		return err
	}

	return nil
}
