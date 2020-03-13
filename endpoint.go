package ep

import (
	"net/http"

	"github.com/advanderveer/ep/coding"
)

// Input describes what is read from a request as input to the endpoint.
type Input interface{}

// CheckerInput can be implemented by Inputs to allow them to validate themselves
type CheckerInput interface {
	Input
	Check() (err error)
}

// ReaderInput is an input that may optionally be implemented by inputs to
// indicate that it has custom logic for reading the request.
type ReaderInput interface {
	Input
	Read(r *http.Request) error
}

// An Output represents one item that results from the endpoint is will be
// send as the response back to the client.
type Output interface{}

// HeaderOutput can be optionally implementedd by outputs to customize the
// response headers
type HeaderOutput interface {
	Output
	Head(http.ResponseWriter, *http.Request) error
}

// ErrorOutput is output data that represents a client or server error. Encoders
// might decide to handle these differently from regular outputs
type ErrorOutput interface {
	epcoding.ErrorEncode
}

// Endpoint can be implemented to descibe an HTTP endpoint
type Endpoint interface {
	Handle(*Response, *http.Request)
}

// ConfigurerEndpoint can be implemented by endpoints to overwrite the default
// configuration
type ConfigurerEndpoint interface {
	Endpoint
	// @TODO rename config to something a bit more description
	// Config is called whenever the endpoint is turned into an http.Handler
	Config(*Config)
}
