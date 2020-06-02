package ep

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/advanderveer/ep/v2/coding"
)

func TestNegotiate(t *testing.T) {
	for i, c := range []struct {
		method string
		body   string
		opts   []Option
		ct     string
		accept string

		expErr    error
		expEnc    coding.Encoder
		expDec    coding.Decoder
		expEncErr error
		expEncCT  string
	}{
		{
			expErr:    nil,
			expEncErr: Err(Op("negotiateEncoder"), ServerError),
		},
		{
			method:    "POST",
			body:      "{}",
			expErr:    Err(Op("negotiateDecoder"), UnsupportedError),
			expEncErr: Err(Op("negotiateEncoder"), ServerError),
		},
		{
			method:   "POST",
			body:     "{}",
			opts:     []Option{RequestDecoding(coding.JSON{}), ResponseEncoding(coding.JSON{})},
			expEnc:   new(json.Encoder),
			expDec:   new(json.Decoder),
			expEncCT: "application/json",
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			r := httptest.NewRequest(c.method, "/", strings.NewReader(c.body))
			r.Header.Set("Content-Type", c.ct)
			r.Header.Set("Accept", c.accept)

			w := httptest.NewRecorder()
			a := New(c.opts...)

			res := newResponse(w, r, a.reqHooks, a.resHooks, a.errHooks, a.decodings, a.encodings)
			if !errors.Is(res.decNegotiateErr, c.expErr) {
				t.Fatalf("expected error %#v, got: %#v", c.expErr, res.decNegotiateErr)
			}

			if !errors.Is(res.encNegotiateErr, c.expEncErr) {
				t.Fatalf("expected error %#v, got: %#v", c.expEncErr, res.encNegotiateErr)
			}

			if c.expEncCT != res.encContentType {
				t.Fatalf("expected encoding ct to be %s, got: %s", c.expEncCT, res.encContentType)
			}

			if c.expEnc == nil && res.enc != nil {
				t.Fatalf("expected nil encoder, got: %#v", res.enc)
			} else if reflect.TypeOf(res.enc) != reflect.TypeOf(c.expEnc) {
				t.Fatalf("expected decoder type %T, got: %T", c.expEnc, res.enc)
			}

			if c.expDec == nil && res.dec != nil {
				t.Fatalf("expected nil encoder, got: %#v", res.dec)
			} else if reflect.TypeOf(res.dec) != reflect.TypeOf(c.expDec) {
				t.Fatalf("expected decoder type %T, got: %T", c.expDec, res.dec)
			}
		})
	}
}

func TestResponseWriting(t *testing.T) {
	hook1 := func(w http.ResponseWriter, r *http.Request, out interface{}) {
		w.WriteHeader(404)
		w.Write([]byte("rab"))
	}

	hook2 := func(w http.ResponseWriter, r *http.Request, out interface{}) {
		w.Write([]byte("bar"))
	}

	hook3 := func(w http.ResponseWriter, r *http.Request, out interface{}) {
		http.Redirect(w, r, "/", 303)
	}

	hook4 := func(w http.ResponseWriter, r *http.Request, out interface{}) {
		w.Write([]byte("foobar"))
		w.WriteHeader(404) // this is ignore, write calls writeheader implicitely
	}

	for i, c := range []struct {
		hook      ResponseHook
		write     string
		expCode   int
		expErr    error
		expBody   string
		expHeader http.Header
	}{
		{
			write:   "foo",
			expCode: 200,
			expBody: "foo",
		},

		{
			hook:    ResponseHook(hook1),
			write:   "foo",
			expCode: 404,
			expBody: "rabfoo",
		},

		{
			hook:    ResponseHook(hook2),
			write:   "foo",
			expCode: 200,
			expBody: "barfoo",
		},

		{
			hook:    ResponseHook(hook3),
			write:   "",
			expCode: 303,
			expBody: `<a href="/">See Other</a>.` + "\n\n",
			expHeader: http.Header{
				"Content-Type": []string{"text/html; charset=utf-8"},
				"Location":     []string{"/"},
			},
		},

		{
			hook:    ResponseHook(hook4),
			write:   "",
			expCode: 200,
			expBody: "foobar",
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			r := httptest.NewRequest("GET", "/", nil)
			w := httptest.NewRecorder()

			var hooks []ResponseHook
			if c.hook != nil {
				hooks = append(hooks, c.hook)
			}

			res := newResponse(w, r, nil, hooks, nil, nil, nil)

			n, err := res.Write([]byte(c.write))
			if len(c.write) != n {
				t.Fatalf("wrote %d, expected %d", n, len(c.write))
			}

			if !errors.Is(err, c.expErr) {
				t.Fatalf("expected error %#v, got: %#v", c.expErr, err)
			}

			if w.Code != c.expCode {
				t.Fatalf("expected code: %d, to: %d", c.expCode, w.Code)
			}

			if w.Body.String() != c.expBody {
				t.Fatalf("expected body: '%s', got: '%s'", c.expBody, w.Body.String())
			}

			if c.expHeader == nil {
				c.expHeader = http.Header{}
			}

			if !reflect.DeepEqual(c.expHeader, w.Header()) {
				t.Fatalf("expected headers: %v, got: %v", c.expHeader, w.Header())
			}
		})
	}
}

func TestMultipleWriteHeaderWithHooks(t *testing.T) {
	var i int
	hook := func(w http.ResponseWriter, r *http.Request, out interface{}) {
		i++
	}

	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	res := newResponse(w, r, nil, []ResponseHook{hook}, nil, nil, nil)

	res.WriteHeader(404)
	res.WriteHeader(405)

	if i != 1 {
		t.Fatalf("expected hook to be called only once, got: %d", i)
	}
}

type emptyOutput struct{}

func (_ emptyOutput) Empty() bool { return true }

func TestPrivateRender(t *testing.T) {
	hook1 := func(err error) (out interface{}) {
		return struct {
			M string `json:"message"`
		}{err.Error()}
	}

	hook2 := func(err error) (out interface{}) {
		return struct {
			M string `json:"error"`
		}{err.Error()}
	}

	for i, c := range []struct {
		out   interface{}
		hooks []ErrorHook
		encs  []coding.Encoding

		expErr    error
		expCode   int
		expBody   string
		expHeader http.Header
		expLogs   string
	}{
		{
			expCode: 200,
		},
		{
			out:     emptyOutput{},
			expCode: 200,
		},
		{
			out:     struct{}{},
			expErr:  Err(Op("negotiateEncoder"), ServerError),
			expCode: 200,
		},
		{
			out:     struct{}{},
			encs:    []coding.Encoding{coding.JSON{}},
			expCode: 200,
			expBody: "{}\n",
			expHeader: http.Header{
				"Content-Type":           {"application/json"},
				"X-Content-Type-Options": {"nosniff"},
			},
		},
		{
			out:     make(chan struct{}), //something that cannot be encoded
			encs:    []coding.Encoding{coding.JSON{}},
			expErr:  Err(Op("response.render"), EncoderError),
			expCode: 200,
		},
		{ //without error hooks, the errors are logged to stdlogger and the error
			// is encoded as is
			out:     Err("my error"),
			encs:    []coding.Encoding{coding.JSON{}},
			expCode: 200,
			expLogs: "my error",
			expBody: "{}\n",
			expHeader: http.Header{
				"Content-Type":           {"application/json"},
				"X-Content-Type-Options": {"nosniff"},
			},
		},
		{ // following tests check that the first hook's result takes precedence
			out:     Err("my error"),
			encs:    []coding.Encoding{coding.JSON{}},
			hooks:   []ErrorHook{hook1, hook2},
			expCode: 200,
			expBody: `{"message":"my error"}` + "\n",
			expHeader: http.Header{
				"Content-Type":           {"application/json"},
				"X-Content-Type-Options": {"nosniff"},
			},
		},
		{
			out:     Err("other error"),
			encs:    []coding.Encoding{coding.JSON{}},
			hooks:   []ErrorHook{hook2, hook1},
			expCode: 200,
			expBody: `{"error":"other error"}` + "\n",
			expHeader: http.Header{
				"Content-Type":           {"application/json"},
				"X-Content-Type-Options": {"nosniff"},
			},
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			lbuf := bytes.NewBuffer(nil)
			log.SetOutput(lbuf)
			defer log.SetOutput(os.Stderr)

			r := httptest.NewRequest("GET", "/", nil)
			w := httptest.NewRecorder()

			res := newResponse(w, r, nil, nil, c.hooks, nil, c.encs)

			err := res.render(c.out)
			if !errors.Is(err, c.expErr) {
				t.Fatalf("expected error %#v, got: %#v", c.expErr, err)
			}

			if w.Code != c.expCode {
				t.Fatalf("expected code: %d, to: %d", c.expCode, w.Code)
			}

			if w.Body.String() != c.expBody {
				t.Fatalf("expected body: '%s', got: '%s'", c.expBody, w.Body.String())
			}

			if c.expHeader == nil {
				c.expHeader = http.Header{}
			}

			if !reflect.DeepEqual(c.expHeader, w.Header()) {
				t.Fatalf("expected headers: %v, got: %v", c.expHeader, w.Header())
			}

			if c.expLogs == "" && lbuf.Len() != 0 {
				t.Fatalf("expected logs to be empty, got: %v", lbuf.String())
			} else if !strings.Contains(lbuf.String(), c.expLogs) {
				t.Fatalf("expected logger to contain: '%v', got: '%v'", c.expLogs, lbuf.String())
			}
		})
	}
}

func TestRenderWithDeliberateNilEncoder(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	res := newResponse(w, r, nil, nil, nil, nil, []coding.Encoding{coding.JSON{}})
	res.enc = nil // this will cause the error we are looking to test

	err := res.render("foo")
	if !errors.Is(err, Err(Op("response.render"), ServerError)) {
		t.Fatalf("expected error, got: %#v", err)
	}
}

func TestPrivateRenderSequentially(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	res := newResponse(w, r, nil, nil, nil, nil, []coding.Encoding{coding.JSON{}})

	err := res.render("foo")
	if err != nil {
		t.Fatalf("unexpected, got: %v", err)
	}

	err = res.render("bar")
	if err != nil {
		t.Fatalf("unexpected, got: %v", err)
	}

	if w.Body.String() != `"foo"`+"\n"+`"bar"`+"\n" {
		t.Fatalf("unexpected, got: %v", w.Body.String())
	}
}

type input1 struct{}

func (_ input1) Empty() bool { return true }

func TestPrivateBind(t *testing.T) {
	errhook := func(r *http.Request, in interface{}) error {
		return errors.New("foo")
	}

	for i, c := range []struct {
		method string
		body   string
		in     interface{}
		hooks  []RequestHook
		decs   []coding.Decoding
		expErr error
		expOK  bool
		expIn  interface{}
	}{
		//empty request should bind as OK as it just means that the handler
		//has to work with the zero value
		{expOK: true},
		{
			in:     nil,
			method: "POST", body: `{"Foo": "bar"}`,
			expOK:  true,
			expErr: nil, // because input is nil
		},
		{
			in:     &struct{}{},
			method: "POST", body: ``, // empty bodies are ignored without error
			expOK:  true,
			expErr: nil,
			expIn:  &struct{}{},
		},
		{
			in:     &input1{},
			method: "POST", body: `{"Foo": "bar"}`,
			expOK:  true,
			expErr: nil, // because input has Empty() method
			expIn:  &input1{},
		},
		{
			in:     &struct{}{},
			method: "POST", body: `{"Foo": "bar"}`,
			expErr: Err(Op("negotiateDecoder"), UnsupportedError),
			expIn:  &struct{}{},
		},
		{
			in:     struct{}{}, // not a pointer so decoder fails
			method: "POST", body: `{"Foo": "bar"}`,
			decs:   []coding.Decoding{coding.JSON{}},
			expErr: Err(Op("response.bind"), DecoderError),
			expIn:  struct{}{},
		},
		{
			in:     &struct{ Foo string }{},
			method: "POST", body: `{"Foo": "bar"}`,
			decs:  []coding.Decoding{coding.JSON{}},
			expOK: true,
			expIn: &struct{ Foo string }{"bar"},
		},
		{
			hooks:  []RequestHook{errhook},
			expOK:  false,
			expErr: Err(Op("response.bind"), RequestHookError),
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			r := httptest.NewRequest(c.method, "/", strings.NewReader(c.body))
			w := httptest.NewRecorder()

			res := newResponse(w, r, c.hooks, nil, nil, c.decs, nil)

			ok, err := res.bind(c.in)
			if !errors.Is(err, c.expErr) {
				t.Fatalf("expected error %#v, got: %#v", c.expErr, err)
			}

			if ok != c.expOK {
				t.Fatalf("expected bind OK to be: %v, got: %v", c.expOK, ok)
			}

			if !reflect.DeepEqual(c.in, c.expIn) {
				t.Fatalf("expected bound input to be: %#v, got: %#v", c.expIn, c.in)
			}
		})
	}
}

func TestPrivateBindSequential(t *testing.T) {
	r := httptest.NewRequest("POST", "/", strings.NewReader(`{"Foo": "bar"}`+"\n"+`{"Foo": "rab"}`))
	w := httptest.NewRecorder()

	var n int
	hook1 := func(r *http.Request, in interface{}) error {
		n++
		return nil
	}

	res := newResponse(w, r, []RequestHook{hook1}, nil, nil, []coding.Decoding{coding.JSON{}}, nil)

	for i := 0; i < 100; i++ {
		var in struct{ Foo string }
		ok, err := res.bind(&in)
		if err != nil {
			t.Fatalf("unexpected, got: %v", err)
		}

		if !ok {
			break
		}

		switch i {
		case 0:
			if in.Foo != "bar" {
				t.Fatalf("unexpected, got: %v", in.Foo)
			}

		case 1:
			if in.Foo != "rab" {
				t.Fatalf("unexpected, got: %v", in.Foo)
			}

		default:
			t.Fatalf("unexpected iteration: %d", i)
		}
	}

	// the hook will be called before the decoder encounters EOF, so the input
	// might be partially populated from the hooks
	if n != 3 {
		t.Fatalf("unexpected, got: %d", n)
	}
}

func TestRenderErrorPrecedence(t *testing.T) {
	e := errors.New("my error")
	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	h := shouldRenderErr(t, e)
	res := newResponse(w, r, nil, nil, []ErrorHook{h}, nil, []coding.Encoding{coding.JSON{}})

	res.Render(struct{}{}, e)
}

func TestRenderDoublePassPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	// error hook will be called after the first pass fails to render but will
	// return an output that can also not be rendered. It should give up then
	hook := func(err error) (out interface{}) {
		return make(chan struct{})
	}

	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	res := newResponse(w, r, nil, nil, []ErrorHook{hook}, nil, []coding.Encoding{coding.JSON{}})

	res.Render(make(chan struct{}), nil)
}

func TestBindSuccess(t *testing.T) {
	r := httptest.NewRequest("POST", "/", strings.NewReader(`{"Foo": "bar"}`))
	w := httptest.NewRecorder()

	res := newResponse(w, r, nil, nil, nil, []coding.Decoding{coding.JSON{}}, nil)

	var in struct{ Foo string }
	ok := res.Bind(&in)
	if !ok {
		t.Fatalf("unexpected, got: %v", ok)
	}

	if in.Foo != "bar" {
		t.Fatalf("unexpected, got: %v", in.Foo)
	}
}

func TestBindError(t *testing.T) {
	r := httptest.NewRequest("POST", "/", strings.NewReader(`{}`))
	w := httptest.NewRecorder()

	h := shouldRenderErr(t, Err(Op("response.bind")))
	res := newResponse(w, r, nil, nil, []ErrorHook{h},
		[]coding.Decoding{coding.JSON{}}, []coding.Encoding{coding.JSON{}})

	ok := res.Bind(struct{}{})
	if ok {
		t.Fatalf("unexpected, got: %v", ok)
	}
}

func TestRecover(t *testing.T) {
	e1 := errors.New("bar")

	for i, c := range []struct {
		panic  interface{}
		expErr error
	}{
		{panic: "foo", expErr: Err(Op("response.Recover"), ServerError)},
		{panic: e1, expErr: Err(Op("response.Recover"), ServerError, e1)},
		{panic: 1, expErr: Err(Op("response.Recover"), ServerError, "unknown panic")},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			r := httptest.NewRequest("GET", "/", nil)
			w := httptest.NewRecorder()

			h := shouldRenderErr(t, c.expErr)
			res := newResponse(w, r, nil, nil, []ErrorHook{h}, nil, []coding.Encoding{coding.JSON{}})

			func() {
				defer res.Recover()
				panic(c.panic)
			}()
		})
	}
}

func TestPanicInResponseHook(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	var panicedAlready bool
	h := func(w http.ResponseWriter, r *http.Request, out interface{}) {
		if !panicedAlready {
			panicedAlready = true
			panic("foo")
		}
	}

	sh := shouldRenderErr(t, Err(Op("response.Recover")))
	res := newResponse(w, r, nil, []ResponseHook{h}, []ErrorHook{sh}, nil, []coding.Encoding{coding.JSON{}})

	func() {
		defer res.Recover()
		res.Render(nil, nil)
	}()

	if res.runningReqHooks == true {
		t.Fatalf("should have reset runningReqHooks state, got: %v", res.runningReqHooks)
	}

}

func shouldRenderErr(t *testing.T, target error) func(error) interface{} {
	return func(err error) interface{} {
		if !errors.Is(err, target) {
			t.Fatalf("expected error %#v, got: %#v", target, err)
		}

		return nil
	}
}
