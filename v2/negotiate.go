package ep

import (
	"bytes"
	"net/http"

	"github.com/advanderveer/ep/v2/accept"
	"github.com/advanderveer/ep/v2/coding"
)

// detectContentType attempt to guess the content type by looking at the
// provided bytes.
func detectContentType(b []byte) (ct string) {
	ct = http.DetectContentType(b)
	if ct != "text/plain; charset=utf-8" {
		return ct
	}

	b = bytes.TrimSpace(b)
	if len(b) > 0 {

		// in our usecase of detecting request bodies we can be a bit more
		// liberal and assume that the client tries to send either JSON or XML
		switch {
		case b[0] == '{':
			fallthrough
		case b[0] == '"':
			fallthrough
		case b[0] == '[':
			return "application/json; charset=utf-8"
		}
	}

	return
}

// negotiateDecoder will inspect the request and available body decoders
// to figure out which should be used. It returns nil when the request's body
// should not be considered by the server and an error if this is unexpected.
func negotiateDecoder(r *http.Request, decs []coding.Decoding) (coding.Decoder, error) {
	const op Op = "negotiateDecoder"

	// If the request has a content-type it has explicitely indicated to have
	// content and of a certain type. Else we sniff to make sure it's empty
	ct := r.Header.Get("Content-Type")
	if ct == "" {
		prc := Buffer(r.Body)
		peek, _ := prc.Peek(512)
		if len(peek) < 1 {
			return nil, nil // it's for sure empty. That's fine
		}

		r.Body = prc

		// attempt to detect it
		ct = detectContentType(peek)
	}

	// At this point we know the client really has sent something to us so if
	// we can't handle it the negotiation fails
	if len(decs) < 1 {
		return nil, Err(op,
			"client sent request body but no decoders configured",
			UnsupportedError,
		)
	}

	// Parse the content header, we only care about the value
	value, _ := accept.ParseValueAndParams(ct)

	// Turn the decodings into asks for the negotiation algorithm
	asks := make([]string, 0, len(decs))
	for _, dec := range decs {
		asks = append(asks, dec.Accepts())
	}

	// finally, negotiate what is necessary for the content type
	_, aski := accept.Negotiate(asks, []string{value})
	if aski < 0 {
		return nil, Err(op,
			"non-empty request body no configured decoder accepts it",
			UnsupportedError,
		)
	}

	return decs[aski].Decoder(r), nil
}

func negotiateEncoder(
	r *http.Request,
	w http.ResponseWriter,
	encs []coding.Encoding,
) (coding.Encoder, string, error) {
	const op Op = "negotiateEncoder"

	if len(encs) < 1 {
		return nil, "", Err(op, "no encoders configured", ServerError)
	}

	asks := r.Header.Values("Accept")
	if len(asks) < 1 || r.Header.Get("Accept") == "" {
		return encs[0].Encoder(w), encs[0].Produces(), nil
	}

	offers := make([]string, 0, len(encs))
	for _, enc := range encs {
		offers = append(offers, enc.Produces())
	}

	offeri, _ := accept.Negotiate(asks, offers)
	if offeri < 0 {
		return nil, "", Err(op,
			"no configured encoder produces what the client accepts",
			UnacceptableError,
		)
	}

	return encs[offeri].Encoder(w), encs[offeri].Produces(), nil
}
