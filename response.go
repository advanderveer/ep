package ep

import (
	"errors"
	"net/http"

	"github.com/advanderveer/ep/coding"
)

var (
	// InvalidInput can be used explicitely to render the response as an invalid
	// input instead of an server error
	InvalidInput = errors.New("invalid input")
)

// Response is an http.ResponseWriter implementation that comes with
// a host of untility method for common tasks in http handling.
type Response struct {
	wr  http.ResponseWriter
	req *http.Request
	cfg ConfReader
	dec epcoding.Decoder
	enc epcoding.Encoder

	state struct {
		wroteHeader bool
		validErr    error // Validation error
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
		res.enc = e.Encoder(wr)
	}

	if d := Decoding(req.Context()); d != nil {
		res.dec = d.Decoder(req)
	}

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
	case r.state.validErr != nil:
		return r.state.validErr
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

	// it is valid to have no decoder and just rely on the read implementation
	// to populate the input.
	if r.dec == nil {
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
		r.state.validErr = v.Validate(in)
		if r.state.validErr != nil {
			return r.state.validErr
		}
	}

	// inputs may optionally implement a validation method
	if incheck, ok := in.(CheckerInput); ok {
		r.state.validErr = incheck.Check()
	}

	return r.state.validErr
}

// Render will assert the provided error and earlier errors and provide
// appropriate feedback in the response. If 'err' is not the same error
// as returned by Validate() it will be handled as a server error.
func (r *Response) Render(err error, out Output) {
	if err != nil && err != r.state.validErr && err != InvalidInput {
		r.state.serverErr = err
	}

	err = r.render(out) // first pass
	if err != nil {
		err = r.render(nil) // second pass
		if err != nil {
			panic("ep/response: failed to render: " + err.Error())
		}
	}
}

func (r *Response) serverErrorOutput(err error) ErrorOutput {
	f := r.cfg.ServerErrFactory()
	if f == nil {
		return serverErrOutput{http.StatusText(http.StatusInternalServerError)}
	}

	return f(err)
}

func (r *Response) clientErrorOutput(err error) ErrorOutput {
	println("client err:", err.Error())
	f := r.cfg.ClientErrFactory()
	if f == nil {
		return clientErrOutput{http.StatusText(http.StatusBadRequest)}
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
	}

	if out == nil {
		return // nothing to do
	}

	// only call the output's head if the response header was not yet written
	hout, hok := out.(HeaderOutput)
	if !r.state.wroteHeader && hok {
		err = hout.Head(r, r.req)
		if err != nil {
			r.state.serverErr = err
			return
		}
	}

	if r.enc == nil {
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
	r.state.wroteHeader = true // because underlying http.ResponseWriter does so
	return r.wr.Write(b)
}

// WriteHeader implements the http.ResponseWriter's "WriteHeader" method
func (r *Response) WriteHeader(statusCode int) {
	r.state.wroteHeader = true
	r.wr.WriteHeader(statusCode)
}

// serverErrOutput is the output that is returned by default when the response
// gets into the server error state.
type serverErrOutput struct{ ErrorMessage string }

func (out serverErrOutput) Head(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(http.StatusInternalServerError)
	return nil
}

func (out serverErrOutput) IsError() {}

// clientErrOutput is the output that is returned by default when the response
// gets into the client error state
type clientErrOutput struct{ ErrorMessage string }

func (out clientErrOutput) Head(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(http.StatusBadRequest)
	return nil
}

func (out clientErrOutput) IsError() {}
