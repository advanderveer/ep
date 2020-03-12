package epcoding

import (
	"encoding/xml"
	"io"
	"net/http"
)

type XMLDecoding struct{}

func NewXMLDecoding() XMLDecoding { return XMLDecoding{} }
func (d XMLDecoding) Accepts() []string {
	return []string{
		"application/xml",
		"text/xml",
	}
}

func (d XMLDecoding) Decoder(r *http.Request) Decoder { return xml.NewDecoder(r.Body) }

type XMLEncoding struct{}

func NewXMLEncoding() XMLEncoding { return XMLEncoding{} }

func (e XMLEncoding) Produces() string { return "application/xml" }

func (d XMLEncoding) Encoder(r io.Writer) Encoder { return xml.NewEncoder(r) }
