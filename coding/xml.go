package epcoding

import (
	"encoding/xml"
	"io"
	"net/http"
)

// XML Decoding allows deserialization of xml request bodies
type XMLDecoding struct{ accepts []string }

// NewXMLDecoding will init the a new xml decoding
func NewXMLDecoding() *XMLDecoding {
	return &XMLDecoding{[]string{
		"application/xml",
		"text/xml",
	}}
}

// SetAccepts changes what this decoding accepts
func (d *XMLDecoding) SetAccepts(s []string) *XMLDecoding { d.accepts = s; return d }

// Accepts returns what content types this decoding accepts
func (d XMLDecoding) Accepts() []string {
	return d.accepts
}

// Decoder creates an actual decoder that can be used
func (d XMLDecoding) Decoder(r *http.Request) Decoder { return xml.NewDecoder(r.Body) }

// XMLEncoding describes response encoding using xml
type XMLEncoding struct{ produces string }

// NewXMLEncoding inits xml encoding
func NewXMLEncoding() *XMLEncoding { return &XMLEncoding{"application/xml"} }

// Produces returns what content type this encoding produces
func (e XMLEncoding) Produces() string { return e.produces }

// SetProduces will set what is produced
func (e *XMLEncoding) SetProduces(p string) *XMLEncoding { e.produces = p; return e }

// Encoder will encode from the provides
func (d XMLEncoding) Encoder(r io.Writer) Encoder { return xml.NewEncoder(r) }
