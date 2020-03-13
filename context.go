package ep

import (
	"context"

	"github.com/advanderveer/ep/coding"
)

// epContextKey reservers a specific type for our context keys
type epContextkey string

// Language will return the language as negotiated with the client. This will
// only be empty if the server didn't provide any languages
func Language(ctx context.Context) (s string) {
	if v := ctx.Value(epContextkey("language")); v != nil {
		s = v.(string)
	}
	return
}

// Encoding will return the response encoding as negotiated with the client.
// This will only be empty if the server didn't specify any supported encodings
func Encoding(ctx context.Context) (enc epcoding.Encoding) {
	if v := ctx.Value(epContextkey("encoding")); v != nil {
		enc = v.(epcoding.Encoding)
	}
	return
}

// Decoding will return the request decoding that was determined based on
// request headers and MIME sniffing. This will only be empty if the server
// didn't configure any decodings.
func Decoding(ctx context.Context) (dec epcoding.Decoding) {
	if v := ctx.Value(epContextkey("decoding")); v != nil {
		dec = v.(epcoding.Decoding)
	}
	return
}
