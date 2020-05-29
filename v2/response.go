package ep

import (
	"io"
	"log"
	"net/http"

	"github.com/advanderveer/ep/v2/coding"
)

// ResponseWriter extends the traditional http.ResponseWriter interface with
// functionality that standardizes input decoding and output encoding.
type ResponseWriter interface {
	Negotiate([]coding.Decoding, []coding.Encoding) bool
	Bind(in interface{}) bool
	Render(out interface{}, err error)
	http.ResponseWriter
}

type response struct {
	http.ResponseWriter
	req *http.Request

	reqHooks []RequestHook
	resHooks []ResponseHook
	errHooks []ErrorHook

	enc coding.Encoder
	dec coding.Decoder

	encContentType  string
	encNegotiateErr error
	wroteHeader     bool
	runningReqHooks bool
	currentOutput   interface{}
}

func newResponse(
	w http.ResponseWriter,
	r *http.Request,
	reqh []RequestHook,
	resh []ResponseHook,
	errh []ErrorHook,
) *response {
	return &response{
		ResponseWriter: w,
		req:            r,

		reqHooks: reqh,
		resHooks: resh,
		errHooks: errh,
	}
}

// NewResponse initializes a ResponseWriter
func NewResponse(
	w http.ResponseWriter,
	r *http.Request,
	reqh []RequestHook,
	resh []ResponseHook,
	errh []ErrorHook,
) ResponseWriter {
	return newResponse(w, r, reqh, resh, errh)
}

// Negotiate which encoding and decoding will be used for binding and rendering
func (res *response) Negotiate(
	decs []coding.Decoding,
	encs []coding.Encoding,
) bool {
	err := res.negotiate(decs, encs)
	if err != nil {
		res.Render(nil, err)
		return false
	}

	return true
}

func (res *response) negotiate(
	decs []coding.Decoding,
	encs []coding.Encoding,
) (err error) {
	res.dec, err = negotiateDecoder(res.req, decs)
	if err != nil {
		return err
	}

	// Failing to negotiate an encoder is only important when we know for sure
	// that it will be used during a call to render. The user might decide to
	// write to the response itself, or the API doesn't need encoding at all.
	// So we keep the error in the response to be reported later.
	res.enc, res.encContentType, res.encNegotiateErr = negotiateEncoder(res.req, res, encs)
	return
}

// Write some data to the response body. If the header was not yet written this
// method will make an implicit call to WriteHeader which will call any hooks
// in order.
func (res *response) Write(b []byte) (int, error) {
	if !res.wroteHeader {
		res.WriteHeader(http.StatusOK)
	}

	return res.ResponseWriter.Write(b)
}

// WriteHeader will call any configured hooks and sends the http response header
// with the resulting status code.
func (res *response) WriteHeader(statusCode int) {
	if res.wroteHeader {
		return
	}

	if !res.runningReqHooks {
		res.runningReqHooks = true
		for _, h := range res.resHooks {
			h(
				res,
				res.req,
				res.currentOutput, // might be nil
			)
		}
		res.runningReqHooks = false
	}

	// this check ensures that if any hooks called writeHeader we won't be
	// calling it again.
	if res.wroteHeader {
		return
	}

	res.ResponseWriter.WriteHeader(statusCode)
	res.wroteHeader = true
}

// Bind will decode the next value from the request into the input 'in'
func (res *response) Bind(in interface{}) bool {
	ok, err := res.bind(in)
	if err != nil {
		res.Render(nil, err)
		return false
	}

	return ok
}

func (res *response) bind(in interface{}) (ok bool, err error) {
	const op Op = "response.bind"

	for _, h := range res.reqHooks {
		if err := h(res.req, in); err != nil {
			return false, Err(op, "request hook failed", err, RequestHookError)
		}
	}

	if res.dec == nil {
		return false, nil
	}

	err = res.dec.Decode(in)
	if err == io.EOF {
		return false, nil
	} else if err != nil {
		return false, Err(op, "request body decoder failed", DecoderError)
	}

	return true, nil
}

// Render will encode the output into the response body
func (res *response) Render(out interface{}, err error) {
	if err != nil {
		out = err // error render takes precedence over
	}

	err = res.render(out) // first pass
	if err != nil {
		err = res.render(err) // second pass
		if err != nil {
			panic("ep/response: failed to render: " + err.Error())
		}
	}
}

// render just the output value
func (res *response) render(v interface{}) (err error) {
	const op Op = "response.render"

	if errv, ok := v.(error); ok {

		// If there was an error but no hooks to turn it into an output
		// we log this situation to the default logger so the user knows
		// whats going on.
		if len(res.errHooks) < 1 {
			log.Printf("ep: no error hooks for rendering error: %v", errv)
		}

		// Error hooks are responsible for turning any error into an output
		// that can be rendered by the encoder.
		var foundErrOutput bool
		for _, h := range res.errHooks {
			if eout := h(errv); eout != nil {
				v = eout
				foundErrOutput = true
				break
			}
		}

		// If no output was delivered by the hooks, we set the value to nil
		// so errors are not accidentaly encoded with sensitive information.
		if !foundErrOutput {
			v = nil
		}
	}

	// If the value turns out to be nil, we won't be needing any encoder but
	// still wanna write the header.
	if v == nil {
		res.WriteHeader(http.StatusOK)
		return nil
	}

	// We for sure have a value to encode, so if we had any issues with getting
	// an encoder we will stop here
	if res.encNegotiateErr != nil {
		return res.encNegotiateErr
	}

	// WriteHeader needs access to the output value but the interface refrains
	// us from passing it as an argument so we need to set it as an temporary
	// struct member.
	res.currentOutput = v
	defer func() { res.currentOutput = nil }()

	// The encoder is nil without any enc negotiation error
	if res.enc == nil {
		return Err(op, "no encoder to serialize non-nil output value", ServerError)
	}

	var ctFromEnc bool
	if res.Header().Get("Content-Type") == "" {
		ctFromEnc = true
		res.Header().Set("Content-Type", res.encContentType)
	}

	err = res.enc.Encode(v)
	if err != nil {

		// If we just added the content-type header but the encoding fails we
		// reset it such that a subsequent call to render can set it again.
		if ctFromEnc {
			res.Header().Del("Content-Type")
		}

		return Err(op, "response body encoder failed", err, EncoderError)
	}

	return
}
