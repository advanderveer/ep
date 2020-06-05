package ep

import (
	"io"
	"log"
	"net/http"
	"reflect"

	"github.com/advanderveer/ep/epcoding"
)

// ResponseWriter extends the traditional http.ResponseWriter interface with
// functionality that standardizes input decoding and output enepcoding.
type ResponseWriter interface {
	Bind(in interface{}) bool
	Render(outs ...interface{})
	Recover()
	http.ResponseWriter
}

type response struct {
	http.ResponseWriter
	req *http.Request

	reqHooks []RequestHook
	resHooks []ResponseHook
	errHooks []ErrorHook

	enc epcoding.Encoder
	dec epcoding.Decoder

	encContentType  string
	encNegotiateErr error
	decNegotiateErr error
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

	decs []epcoding.Decoding,
	encs []epcoding.Encoding,
) *response {
	res := &response{
		ResponseWriter: w,
		req:            r,

		reqHooks: reqh,
		resHooks: resh,
		errHooks: errh,
	}

	// any failure to negotiate is only important if we actually wanna decode
	// something during a call to bind
	res.dec, res.decNegotiateErr = negotiateDecoder(res.req, decs)

	// Failing to negotiate an encoder is only important when we know for sure
	// that it will be used during a call to render. The user might decide to
	// write to the response itself, or the API doesn't need encoding at all.
	// So we keep the error in the response to be reported later.
	res.enc, res.encContentType, res.encNegotiateErr = negotiateEncoder(res.req, res, encs)
	return res
}

// NewResponse initializes a ResponseWriter
func NewResponse(
	w http.ResponseWriter,
	r *http.Request,
	reqh []RequestHook,
	resh []ResponseHook,
	errh []ErrorHook,
	decs []epcoding.Decoding,
	encs []epcoding.Encoding,
) ResponseWriter {
	return newResponse(w, r, reqh, resh, errh, decs, encs)
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
		defer func() { res.runningReqHooks = false }()
		for _, h := range res.resHooks {
			h(
				res,
				res.req,
				res.currentOutput, // might be nil
			)
		}
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

	// if the input is nil or has an SkipDecode() method we skip decoding
	switch vt := in.(type) {
	case nil:
		return true, nil
	case interface{ SkipDecode() bool }:
		if vt.SkipDecode() {
			return true, nil
		}
	}

	if res.decNegotiateErr != nil {
		return false, res.decNegotiateErr
	}

	if res.dec == nil {
		return true, nil
	}

	err = res.dec.Decode(in)
	if err == io.EOF {
		return false, nil
	} else if err != nil {
		return false, Err(op, "request body decoder failed", err, DecoderError)
	}

	return true, nil
}

// Render will encode the first non-nil argument into the response body. If any
// of the arguments is an error, it takes precedence and is rendered instead.
func (res *response) Render(outs ...interface{}) {
	var out interface{}
	for _, o := range outs {

		switch o.(type) {
		case nil:
		case error:
			out = o
		default:
			if out != nil {
				continue
			}

			// as a special exception it should also skip arguments that are nil
			// pointers. Overhead has been benchmarked at 3-4ns
			if rv := reflect.ValueOf(o); rv.Kind() == reflect.Ptr && rv.IsNil() {
				continue
			}

			out = o
		}
	}

	err := res.render(out) // first pass
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
			log.Printf("ep: no error hooks to render error: %v", errv)
		}

		// Error hooks are responsible for turning any error into an output
		// that can be rendered by the encoder.
		// var foundErrOutput bool
		for _, h := range res.errHooks {
			if eout := h(errv); eout != nil {
				v = eout
				break
			}
		}
	}

	// WriteHeader needs access to the output value but the interface refrains
	// us from passing it as an argument so we need to set it as an temporary
	// struct member.
	res.currentOutput = v
	defer func() { res.currentOutput = nil }()

	// If the value turns out to be nil or implements the Empty() method, we won't
	// be needing any encoder but still wanna write the header and call any
	// response hooks.
	switch vt := v.(type) {
	case nil:
		res.WriteHeader(http.StatusOK)
		return nil
	case interface{ Empty() bool }:
		if vt.Empty() {
			res.WriteHeader(http.StatusOK)
			return nil
		}
	}

	// We for sure have a value to encode, so if we had any issues with getting
	// an encoder we will stop here
	if res.encNegotiateErr != nil {
		return res.encNegotiateErr
	}

	// The encoder is nil without any enc negotiation error
	if res.enc == nil {
		return Err(op, "no encoder to serialize non-nil output value", ServerError)
	}

	var ctFromEnc bool
	if res.Header().Get("Content-Type") == "" {
		ctFromEnc = true
		res.Header().Set("Content-Type", res.encContentType)

		// we know the content-type for sure so we can prevent content sniffing
		res.Header().Set("X-Content-Type-Options", "nosniff")
	}

	err = res.enc.Encode(v)
	if err != nil {

		// If we just added the content-type header but the encoding fails we
		// reset it such that a subsequent call to render can set it again.
		if ctFromEnc {
			res.Header().Del("Content-Type")
			res.Header().Del("X-Content-Type-Options")
		}

		return Err(op, "response body encoder failed", err, EncoderError)
	}

	return
}

func (res *response) Recover() {
	r := recover()
	if r == nil {
		return
	}

	var perr error
	switch rt := r.(type) {
	case error:
		perr = Err(Op("response.Recover"), "error", ServerError, rt)
	case string:
		perr = Err(Op("response.Recover"), rt, ServerError)
	default:
		perr = Err(Op("response.Recover"), "unknown panic", ServerError)
		return
	}

	res.Render(nil, perr)
}
