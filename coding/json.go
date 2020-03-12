package epcoding

import (
	"encoding/json"
	"io"
	"net/http"
)

type JSONDecoding struct{}

func NewJSONDecoding() JSONDecoding { return JSONDecoding{} }
func (d JSONDecoding) Accepts() []string {
	return []string{
		"application/json",
	}
}

func (d JSONDecoding) Decoder(r *http.Request) Decoder { return json.NewDecoder(r.Body) }

type JSONEncoding struct{}

func NewJSONEncoding() JSONEncoding { return JSONEncoding{} }

func (e JSONEncoding) Produces() string { return "application/json" }

func (d JSONEncoding) Encoder(r io.Writer) Encoder { return json.NewEncoder(r) }
