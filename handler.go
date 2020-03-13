package ep

import (
	"context"
	"net/http"

	"github.com/advanderveer/ep/accept"
	"github.com/advanderveer/ep/coding"
)

// Handler will create a http.Handler from the provided endpoint.
func Handler(ep Endpoint) http.Handler {
	var cfg *Config
	if cep, ok := ep.(ConfigurerEndpoint); ok {
		cfg = cep.Config() //@TODO allow endpoint to inherit configuration
	}

	// @TODO allow default config to be set
	if cfg == nil {
		cfg = &Config{}
	}

	// @TODO make sure it is safe to read config from multiple routines

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		req = Negotiate(*cfg, req)
		ep.Handle(NewResponse(w, req, *cfg), req)
	})
}

// Negotiate will look at the requests' headers and set context values for
// encoding, decoding and language.
func Negotiate(cfg Config, req *http.Request) *http.Request {
	req = req.WithContext(
		context.WithValue(req.Context(), epContextkey("lang"),
			accept.Negotiate("Accept-Language", req.Header, cfg.langs, ""),
		))

	if e := epcoding.NegotiateEncoding(req.Header, cfg.encs); e != nil {
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

	if d := epcoding.NegotiateDecoding(req.Header, cfg.decs); d != nil {
		req = req.WithContext(
			context.WithValue(req.Context(), epContextkey("decoding"), d),
		)
	}

	return req
}
