package ep

import (
	"errors"
	"net/http"
	"net/url"

	"github.com/advanderveer/ep/coding"
)

var (
	// SkipEncode can be retured by the output head to prevent any further
	// decoding
	SkipEncode = errors.New("skip encode")
)

// Response is an http.ResponseWriter implementation that comes with
// a host of untility method for common tasks in http handling.
type Response struct {
	wr  http.ResponseWriter
	req *http.Request
	cfg ConfReader
	dec epcoding.Decoder
	enc epcoding.Encoder

	// @TODO clean this up
	responseContentType string

	state struct {
		wroteHeader int
		invalidErr  error
		clientErr   error // BadRequest status
		serverErr   error // InternalServerError
	}
}

// NewResponse initializes a new Response
func NewResponse(
	wr http.ResponseWriter,
	req *http.Request,
	cfg ConfReader,
) (res *Response) {
	res = &Response{
		wr:  wr,
		req: req,
		cfg: cfg,
	}

	if e := Encoding(req.Context()); e != nil {
		res.responseContentType = e.Produces()
		res.enc = e.Encoder(wr)
	}

	if d := Decoding(req.Context()); d != nil {
		res.dec = d.Decoder(req)
	}

	res.state.wroteHeader = -1
	return
}

// Error will return any server, client or validation error that was
// encountered while formulating the response.
func (r *Response) Error() error {
	switch {
	case r.state.serverErr != nil:
		return r.state.serverErr
	case r.state.clientErr != nil:
		return r.state.clientErr
	case r.state.invalidErr != nil:
		return r.state.invalidErr
	default:
		return nil
	}
}

// Bind will use the negotiated decoder to populate the input.
func (r *Response) Bind(in Input) (ok bool) {

	// without input, binding succeeds because endpoints without input should
	// always render
	if in == nil {
		return true
	}

	// input may implement a reader function that takes the raw request and
	// should initialize the input. Userfull for header reading, url params
	// and setting default values for input.
	if reqr, reqrok := in.(ReaderInput); reqrok {
		err := reqr.Read(r.req)
		if err != nil {
			r.state.clientErr = err
			r.render(nil)
			return
		}
	}

	// if there is a query decoder configured and the url has a valid query
	// it will be decoded into the input
	if qdec := r.cfg.QueryDecoder(); qdec != nil {
		qvals, err := url.ParseQuery(r.req.URL.RawQuery)
		if err != nil {
			r.state.clientErr = err
			r.render(nil)
			return
		}

		err = qdec.Decode(in, qvals)
		if err != nil {
			r.state.clientErr = err
			r.render(nil)
			return
		}
	}

	// it is valid to have no decoder and just rely on the read implementation
	// to populate the input.
	if r.dec == nil || r.req.ContentLength == 0 || r.req.Body == nil {
		return true
	}

	// with a decoder and input we ask the decoder to deserialize
	err := r.dec.Decode(in)
	if err != nil {
		r.state.clientErr = err // includes io.EOF
		r.render(nil)
		return
	}

	return true
}

// Validate will validate the input and return any error. It will first
// use any struct validator first before using the input's check method.
func (r *Response) Validate(in Input) (verr error) {
	if in == nil {
		return // no input is always valid
	}

	// call any non-custom validation, if configured
	if v := r.cfg.Validator(); v != nil {
		verr = v.Validate(in)
		if verr != nil {
			return verr
		}
	}

	// inputs may optionally implement a validation method
	if incheck, ok := in.(CheckerInput); ok {
		verr = incheck.Check()
	}

	return
}

// Render will assert the provided error and earlier errors and provide
// appropriate feedback in the response. If 'err' is not the same error
// as returned by Validate() it will be handled as a server error.
func (r *Response) Render(out Output, err error) {
	if err != nil {

		// if the error is marked as invalid input we will render
		// it as an parameter error response. Mainly usefull for
		// REST endpoints
		var iierr *InvalidInputError
		if errors.As(err, &iierr) {
			r.state.invalidErr = iierr.Unwrap()
		} else {
			r.state.serverErr = err
		}
	}

	err = r.render(out) // first pass
	if err != nil {
		err = r.render(nil) // second pass
		if err != nil {
			panic("ep/response: failed to render: " + err.Error())
		}
	}
}

func (r *Response) serverErrorOutput(err error) Output {
	f := r.cfg.ServerErrFactory()
	if f == nil {
		return serverErrOutput{http.StatusText(http.StatusInternalServerError)}
	}

	return f(err)
}

func (r *Response) clientErrorOutput(err error) Output {
	f := r.cfg.ClientErrFactory()
	if f == nil {
		return clientErrOutput{http.StatusText(http.StatusBadRequest)}
	}

	return f(err)
}

func (r *Response) invalidErrorOutput(err error) Output {
	f := r.cfg.InvalidErrFactory()
	if f == nil {
		return invalidErrOutput{err.Error()}
	}

	return f(err)
}

// render solely based on the internal state of the response.
func (r *Response) render(out Output) (err error) {

	// if there are any client or server errors they will be turned into
	// outputs.
	switch {
	case r.state.serverErr != nil:
		out = r.serverErrorOutput(r.state.serverErr)
	case r.state.clientErr != nil:
		out = r.clientErrorOutput(r.state.clientErr)
	case r.state.invalidErr != nil:
		out = r.invalidErrorOutput(r.state.invalidErr)
	}

	if out == nil {
		return // nothing to do
	}

	// if there is a content type for the response, set it before header written
	if r.responseContentType != "" && r.state.wroteHeader < 0 {
		r.Header().Set("Content-Type", r.responseContentType)
	}

	// only call the output's head if the response header was not yet written
	hout, hok := out.(HeaderOutput)
	if r.state.wroteHeader < 0 && hok {
		err = hout.Head(r, r.req)
		if err == SkipEncode {
			return
		} else if err != nil {
			r.state.serverErr = err
			return
		}
	}

	if r.enc == nil || !bodyAllowedForStatus(r.state.wroteHeader) {
		return
	}

	err = r.enc.Encode(out)
	if err != nil {
		r.state.serverErr = err
		return
	}

	return
}

// Header implements the http.ResponseWriter's "Header" method
func (r *Response) Header() http.Header {
	return r.wr.Header()
}

// Write implements the http.ResponseWriter's "Write" method
func (r *Response) Write(b []byte) (int, error) {
	r.state.wroteHeader = http.StatusOK // because underlying http.ResponseWriter does so
	return r.wr.Write(b)
}

// WriteHeader implements the http.ResponseWriter's "WriteHeader" method
func (r *Response) WriteHeader(statusCode int) {
	r.state.wroteHeader = statusCode
	r.wr.WriteHeader(statusCode)
}

// bodyAllowedForStatus reports whether a given response status code
// permits a body. See RFC 7230, section 3.3.
func bodyAllowedForStatus(status int) bool {
	switch {
	case status >= 100 && status <= 199:
		return false
	case status == 204:
		return false
	case status == 304:
		return false
	}
	return true
}
