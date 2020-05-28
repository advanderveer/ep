package coding

import (
	"io"
	"net/http"
)

// Encoder serializes a value into the response body
type Encoder interface {
	Encode(v interface{}) error
}

// Decoder deserializes a value from the request body
type Decoder interface {
	Decode(v interface{}) error
}

// Decoding describes deocding that accepts a certain content type
type Decoding interface {
	Accepts() string
	Decoder(r *http.Request) Decoder
}

// Encoding describes encoding that procoduces a certain content type
type Encoding interface {
	Produces() string
	Encoder(w io.Writer) Encoder
}
