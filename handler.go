package ep

import (
	"context"
	"net/http"

	"github.com/advanderveer/ep/accept"
	"github.com/advanderveer/ep/coding"
)

// Handler handles http for certain endpoint
type Handler struct {
	*Conf
	ep Endpoint
}

// ServeHTTP allows the endpoint to serve HTTP
func (h Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	req = Negotiate(h.Conf, req)
	h.ep.Handle(NewResponse(w, req, *h.Conf), req)
}

// Negotiate will look at the requests' headers and set context values for
// encoding, decoding and language.
func Negotiate(cfg ConfReader, req *http.Request) *http.Request {
	req = req.WithContext(
		context.WithValue(req.Context(), epContextkey("language"),
			accept.Negotiate("Accept-Language", req.Header, cfg.Languages(), ""),
		))

	if e := epcoding.NegotiateEncoding(req.Header, cfg.Encodings()); e != nil {
		req = req.WithContext(
			context.WithValue(req.Context(), epContextkey("encoding"), e),
		)
	}

	// if there is a request body we will turn it into a small buffered reader
	// that allows us to sniff the content type and keep progress
	if req.Body != nil {
		body := NewReader(req.Body)
		req.Body = body

		if req.Header.Get("Content-Type") == "" {
			req.Header.Set("Content-Type", body.Sniff())
		}
	}

	if d := epcoding.NegotiateDecoding(req.Header, cfg.Decodings()); d != nil {
		req = req.WithContext(
			context.WithValue(req.Context(), epContextkey("decoding"), d),
		)
	}

	return req
}
