package epcoding

import (
	"encoding/xml"
	"io"
	"net/http"
)

type XMLDecoding struct{ accepts []string }

func NewXMLDecoding() *XMLDecoding {
	return &XMLDecoding{[]string{
		"application/xml",
		"text/xml",
	}}
}

func (d *XMLDecoding) SetAccepts(s []string) *XMLDecoding { d.accepts = s; return d }
func (d XMLDecoding) Accepts() []string {
	return d.accepts
}

func (d XMLDecoding) Decoder(r *http.Request) Decoder { return xml.NewDecoder(r.Body) }

type XMLEncoding struct{ produces string }

func NewXMLEncoding() *XMLEncoding { return &XMLEncoding{"application/xml"} }

func (e XMLEncoding) Produces() string                   { return e.produces }
func (e *XMLEncoding) SetProduces(p string) *XMLEncoding { e.produces = p; return e }

func (d XMLEncoding) Encoder(r io.Writer) Encoder { return xml.NewEncoder(r) }
