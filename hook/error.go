package hook

import (
	"encoding/xml"
	"errors"
	"html/template"
	"log"
	"net/http"

	"github.com/advanderveer/ep"
)

// NewStandardError creates an error hook for handling ep.Error errors. It logs
// errors to the provided logger creates sensible status codes and only revels
// standard HTTP text that is associated with that code. It comes with default
// outputs for the XML, JSON and HTML encoders.
func NewStandardError(logs *log.Logger) func(err error) interface{} {
	return func(err error) interface{} {
		if logs != nil {
			logs.Print(err)
		}

		// we only create outputs for ep.Error types
		var eperr *ep.Error
		if !errors.As(err, &eperr) {
			return nil
		}

		out := errorOutput{status: http.StatusInternalServerError}
		switch {
		case errors.Is(eperr, ep.Err(ep.UnacceptableError)):
			out.status = http.StatusNotAcceptable
		case errors.Is(eperr, ep.Err(ep.UnsupportedError)):
			out.status = http.StatusUnsupportedMediaType
		case errors.Is(eperr, ep.Err(ep.DecoderError)):
			out.status = http.StatusBadRequest
		}

		out.Message = http.StatusText(out.status)
		return out
	}
}

var errorTemplate = template.Must(template.New("").Parse(`{{.Message}}`))

type errorOutput struct {
	status int

	Message string   `json:"message"`
	XMLName xml.Name `json:"-" xml:"Error"`
}

func (out errorOutput) Status() int { return out.status }

func (out errorOutput) Template() *template.Template {
	return errorTemplate
}
