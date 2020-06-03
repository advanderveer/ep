package coding

import (
	"encoding/json"
	"net/http"
)

// JSON encoding and decoding
type JSON struct{}

func (_ JSON) Produces() string {
	return "application/json"
}

func (_ JSON) Encoder(w http.ResponseWriter) Encoder {
	return json.NewEncoder(w)
}

func (_ JSON) Accepts() string {
	return "application/json, application/vnd.api+json"
}

func (_ JSON) Decoder(r *http.Request) Decoder {
	return json.NewDecoder(r.Body)
}
