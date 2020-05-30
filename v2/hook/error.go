package hook

import (
	"encoding/xml"
	"errors"
	"html/template"
	"log"
	"net/http"

	ep "github.com/advanderveer/ep/v2"
)

// NewPrivateError creates an error hook that logs errors to the provided logger
// and doesn't reveal any info to the client except for a status code and the
// standard HTTP text that is associated with that code. It comes with default
// outputs for the XML, JSON and HTML encoders.
func NewPrivateError(logs *log.Logger) func(err error) interface{} {
	return func(err error) interface{} {
		if logs != nil {
			logs.Print(err)
		}

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
