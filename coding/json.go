package epcoding

import (
	"encoding/json"
	"io"
	"net/http"
)

// JSONDecoding defines a decoding that produces json
type JSONDecoding struct {
	accepts []string
}

// NewJSONDecoding inits the json decoding
func NewJSONDecoding() *JSONDecoding {
	return &JSONDecoding{accepts: []string{
		"application/json",
		"application/vnd.api+json", // https://jsonapi.org/
	}}
}

// Accepts defines what content types this decoder can decode
func (d JSONDecoding) Accepts() []string { return d.accepts }

// SetAccepts allows overwriting what content-types will be decoded
func (d *JSONDecoding) SetAccepts(accs []string) *JSONDecoding { d.accepts = accs; return d }

// Decoder creates the actual decoding
func (d JSONDecoding) Decoder(r *http.Request) Decoder { return json.NewDecoder(r.Body) }

// JSONEncoding encodes json
type JSONEncoding struct {
	produces string
}

// NewJSONEncoding creates a new JSON encoder
func NewJSONEncoding() JSONEncoding { return JSONEncoding{"application/json"} }

// Produces defines what content type the encoding encodes to
func (e JSONEncoding) Produces() string { return e.produces }

// SetProduces overwrites the content-type this encoder will produce
func (e *JSONEncoding) SetProduces(s string) *JSONEncoding { e.produces = s; return e }

// Encoder creates an anctual json encoder
func (d JSONEncoding) Encoder(r io.Writer) Encoder { return json.NewEncoder(r) }
