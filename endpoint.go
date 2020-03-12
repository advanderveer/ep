package ep

import (
	"net/http"

	"github.com/advanderveer/ep/coding"
)

// Input describes what is read from a request as input to the endpoint.
type Input interface {
	Check() (err error)
	// @TODO maybe the input.Check method should be optional
}

// ReaderInput is an input that may optionally be implemented by inputs to
// indicate that it has custom logic for reading the request.
type ReaderInput interface {
	Input
	Read(r *http.Request) error
}

// An Output represents one item that results from the endpoint is will be
// send as the response back to the client.
type Output interface {
	Head(http.ResponseWriter, *http.Request) error
}

// ErrorOutput is output data that represents a client or server error. Encoders
// might decide to handle these differently from regular outputs
type ErrorOutput interface {
	Output
	epcoding.ErrorEncode
}

// Endpoint can be implemented to descibe an HTTP endpoint
type Endpoint interface {

	// @TODO rename config to something a bit more description
	// Config is called whenever the endpoint is turned into an http.Handler
	Config() *Config
	// Handle is called upon every request to hit the endpoint
	Handle(*Response, *http.Request)
}

// Handler will create a http.Handler from the provided endpoint.
func Handler(ep Endpoint) http.Handler {
	cfg := ep.Config()
	// @TODO make sure it is ok to read config from multiple routines

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ep.Handle(NewResponse(w, req, *cfg), req)
	})
}
