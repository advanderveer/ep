package epcoding

import (
	"io"
	"net/http"
	"strings"

	"github.com/advanderveer/ep/accept"
)

type ErrorEncode interface {
	IsError()
}

type Encoding interface {
	Produces() string
	Encoder(w io.Writer) Encoder
}

type Encoder interface {
	Encode(v interface{}) error
}

type Decoding interface {
	Accepts() []string
	Decoder(r *http.Request) Decoder
}

type Decoder interface {
	Decode(v interface{}) error
}

// NegotiateEncoding will examine the Accept header and select the most
// appropriate encoding. If there are no encodings this function will return nil.
func NegotiateEncoding(h http.Header, encs []Encoding) (enc Encoding) {
	if len(encs) < 1 {
		return
	}

	offers := make([]string, 0, len(encs))
	mapped := make(map[string]Encoding, len(encs))
	for _, enc := range encs {
		offers = append(offers, enc.Produces())
		if _, ok := mapped[enc.Produces()]; ok {
			panic("ep/coding: multiple encodings are registered to produce the same content type")
		}

		mapped[enc.Produces()] = enc
	}

	return mapped[accept.Negotiate("Accept", h, offers, offers[0])]
}

// NegotiateDecoding will examine the Content-Type header and select the most
// appropriate decoder. If there are no decoders this function will return nil.
func NegotiateDecoding(h http.Header, decs []Decoding) (dec Decoding) {
	if len(decs) < 1 {
		return
	}

	var def string
	ac := make(http.Header)
	mapped := make(map[string]Decoding)
	for _, dec := range decs {
		for _, spec := range dec.Accepts() {
			if def == "" {
				def = spec
			}
			ac.Add("Accept", spec)
			mapped[spec] = dec
		}
	}

	// Valid content-type header may include an encoding and boundary parts
	parts := strings.SplitN(h.Get("Content-Type"), ";", 2)

	return mapped[accept.Negotiate("Accept", ac, []string{parts[0]}, def)]
}
