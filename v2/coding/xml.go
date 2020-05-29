package coding

import (
	"encoding/xml"
	"net/http"
)

// XML encoding and decoding
type XML struct{}

func (_ XML) Produces() string {
	return "application/xml"
}

func (_ XML) Encoder(w http.ResponseWriter) Encoder {
	return xml.NewEncoder(w)
}

func (_ XML) Accepts() string {
	return "application/xml, text/xml"
}

func (_ XML) Decoder(r *http.Request) Decoder {
	return xml.NewDecoder(r.Body)
}
