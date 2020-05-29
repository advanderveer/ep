package coding

import (
	"encoding/xml"
	"net/http"
)

// XML decoding
type XML struct{}

func (d XML) Produces() string {
	return "application/xml"
}

func (d XML) Encoder(w http.ResponseWriter) Encoder {
	return xml.NewEncoder(w)
}

func (d XML) Accepts() string {
	return "application/xml, text/xml"
}

func (d XML) Decoder(r *http.Request) Decoder {
	return xml.NewDecoder(r.Body)
}
