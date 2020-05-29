package coding

import (
	"encoding/json"
	"net/http"
)

// JSON decoding
type JSON struct{}

func (d JSON) Produces() string {
	return "application/json"
}

func (d JSON) Encoder(w http.ResponseWriter) Encoder {
	return json.NewEncoder(w)
}

func (d JSON) Accepts() string {
	return "application/json, application/vnd.api+json"
}

func (d JSON) Decoder(r *http.Request) Decoder {
	return json.NewDecoder(r.Body)
}
