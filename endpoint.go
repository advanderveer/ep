package ep

import (
	"net/http"
)

// Input describes what is read from a request as input to the endpoint.
type Input interface{}

// CheckerInput can be implemented by Inputs to allow them to validate themselves
type CheckerInput interface {
	Input
	Validate() (err error)
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

// Endpoint can be implemented to descibe an HTTP endpoint
type Endpoint interface {
	Handle(*Response, *http.Request)
}

// EndpointFunc implements endpoint by providng just the Handle func
type EndpointFunc func(*Response, *http.Request)

func (f EndpointFunc) Handle(res *Response, req *http.Request) { f(res, req) }
