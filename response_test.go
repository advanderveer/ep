package ep

import (
	"bytes"
	"encoding/json"
	"errors"
	"html/template"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/advanderveer/ep/coding"
)

var errIn1 = errors.New("bar")

type in1 struct{ Foo string }

func (in in1) Check() (err error) {
	if in.Foo != "" {
		return nil
	}

	return errIn1
}

type qdec1 struct{}

func (d qdec1) Decode(v interface{}, vals url.Values) error {
	json.Unmarshal([]byte(`{"Foo": "`+vals.Get("Foo")+`"}`), v)
	return nil
}

type qdec2 struct{}

func (d qdec2) Decode(v interface{}, vals url.Values) error {
	return errors.New("fail")
}

func TestResponseBinding(t *testing.T) {
	cfg := New()

	t.Run("bind without any input", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/", nil)
		req = Negotiate(*cfg, req)
		res := NewResponse(nil, req, *cfg)
		ok := res.Bind(nil)
		if !ok {
			t.Fatalf("unexpected, got: %v", ok)
		}
	})

	cfg = New().WithDecoding(epcoding.NewJSONDecoding())

	t.Run("bind with input and decoder", func(t *testing.T) {
		var v in1

		req, _ := http.NewRequest("GET", "/", strings.NewReader(`{"Foo": "bar"}`))
		req = Negotiate(*cfg, req)
		res := NewResponse(nil, req, *cfg)
		ok := res.Bind(&v)
		if !ok {
			t.Fatalf("unexpected, got: %v", ok)
		}

		if v.Foo != "bar" {
			t.Fatalf("unexpected, got: %v", v.Foo)
		}
	})

	t.Run("bind with syntax error", func(t *testing.T) {
		var v in1

		rec := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/", strings.NewReader(`{"Foo: "bar"}`))
		req = Negotiate(*cfg, req)
		res := NewResponse(rec, req, *cfg)
		ok := res.Bind(&v)
		if ok {
			t.Fatalf("unexpected, got: %v", ok)
		}

		if v.Foo != "" {
			t.Fatalf("unexpected, got: %v", v.Foo)
		}

		if !strings.Contains(res.Error().Error(), "invalid") {
			t.Fatalf("unexpected, got %v", res.Error())
		}

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("unexpected, got: %v", rec.Code)
		}
	})

	cfg = New().WithDecoding(epcoding.NewJSONDecoding()).WithEncoding(epcoding.NewJSONEncoding())

	t.Run("bind with syntax error, and JSON encoder to render", func(t *testing.T) {
		var v in1

		rec := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/", strings.NewReader(`{"Foo: "bar"}`))
		req = Negotiate(*cfg, req)
		res := NewResponse(rec, req, *cfg)
		ok := res.Bind(&v)
		if ok {
			t.Fatalf("unexpected, got: %v", ok)
		}

		if rec.Body.String() != `{"ErrorMessage":"Bad Request"}`+"\n" {
			t.Fatalf("expected client error encoded, got: %v", rec.Body.String())
		}

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected bad request status, got: %v", rec.Code)
		}
	})

	cfg = cfg.SetQueryDecoder(qdec1{})

	t.Run("bind with valid query decoder", func(t *testing.T) {
		var v in1

		rec := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/?Foo=bar", nil)
		req = Negotiate(*cfg, req)
		res := NewResponse(rec, req, *cfg)
		ok := res.Bind(&v)
		if !ok {
			t.Fatalf("unexpected, got: %v", ok)
		}

		if v.Foo != "bar" {
			t.Fatalf("unexpected, got :%v", v.Foo)
		}
	})

	t.Run("bind with valid invalid query to decoder", func(t *testing.T) {
		var v in1

		rec := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/?Foo=%?bar", nil)
		req = Negotiate(*cfg, req)
		res := NewResponse(rec, req, *cfg)
		ok := res.Bind(&v)
		if ok {
			t.Fatalf("unexpected, got: %v", ok)
		}

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected bad request status, got: %v", rec.Code)
		}
	})

	cfg = cfg.SetQueryDecoder(qdec2{})

	t.Run("bind with valid invalid query and error decoder", func(t *testing.T) {
		var v in1

		rec := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/?Foo=bar", nil)
		req = Negotiate(*cfg, req)
		res := NewResponse(rec, req, *cfg)
		ok := res.Bind(&v)
		if ok {
			t.Fatalf("unexpected, got: %v", ok)
		}

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected bad request status, got: %v", rec.Code)
		}
	})
}

type in2 struct{ Foo string }

func (in in2) Check() error { return nil }
func (in *in2) Read(r *http.Request) error {
	in.Foo = "barr"
	if r.URL.Path == "/bogus" {
		return errors.New("fail")
	}

	return nil
}

func TestResponseBindingWithReaderInput(t *testing.T) {
	cfg := New()

	t.Run("with valid read", func(t *testing.T) {
		var in in2

		req, _ := http.NewRequest("GET", "/", nil)
		req = Negotiate(*cfg, req)
		res := NewResponse(nil, req, *cfg)
		ok := res.Bind(&in)
		if !ok {
			t.Fatalf("should be ok, got: %v", ok)
		}

		if in.Foo != "barr" {
			t.Fatalf("unexpected, got: %v", in.Foo)
		}
	})

	lbuf := bytes.NewBuffer(nil)
	logs := log.New(lbuf, "", 0)
	cfg = cfg.SetLogger(NewStdLogger(logs))

	t.Run("with read error", func(t *testing.T) {
		var in in2

		rec := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/bogus", nil)
		req = Negotiate(*cfg, req)
		res := NewResponse(rec, req, *cfg)
		ok := res.Bind(&in)
		if ok {
			t.Fatalf("unexpected, got: %v", ok)
		}

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("unexpected, got: %v", rec.Code)
		}

		if !strings.Contains(lbuf.String(), "fail") {
			t.Fatalf("log should show error, got: %v", lbuf.String())
		}
	})
}

type val1 struct{}

var val1err = errors.New("invalid")

func (v val1) Validate(interface{}) error { return val1err }

func TestResponseValidation(t *testing.T) {
	cfg := New()

	t.Run("validate without a request body", func(t *testing.T) {
		var v in1

		req, _ := http.NewRequest("GET", "/", nil)
		req = Negotiate(*cfg, req)
		res := NewResponse(nil, req, *cfg)
		err := res.Validate(v)
		if err != errIn1 {
			t.Fatalf("unexpected, got: %v", err)
		}
	})

	t.Run("validate without input", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/", strings.NewReader(`{}`))
		req = Negotiate(*cfg, req)
		res := NewResponse(nil, req, *cfg)
		err := res.Validate(nil)
		if err != nil {
			t.Fatalf("unexpected, got: %v", err)
		}

		if res.Error() != nil {
			t.Fatalf("unexpected, got: %v", res.Error())
		}
	})

	t.Run("validate with input's own logic", func(t *testing.T) {
		var v in1
		req, _ := http.NewRequest("GET", "/", strings.NewReader(`{}`))
		req = Negotiate(*cfg, req)
		res := NewResponse(nil, req, *cfg)

		err := res.Validate(&v)
		if err != errIn1 {
			t.Fatalf("unexpected, got: %v", err)
		}
	})

	cfg.SetValidator(val1{})

	t.Run("with system validator", func(t *testing.T) {
		var in in1
		req, _ := http.NewRequest("GET", "/", strings.NewReader(`{}`))
		req = Negotiate(*cfg, req)
		res := NewResponse(nil, req, *cfg)

		err := res.Validate(&in)
		if err != val1err {
			t.Fatalf("unexpected, got: %v", err)
		}
	})
}

type out1 struct{ Bar string }

var out1Err = errors.New("out1err")

func (o out1) Head(w http.ResponseWriter, r *http.Request) error {
	return out1Err
}

type out2 struct{ Bar chan bool }

func (o out2) Head(w http.ResponseWriter, r *http.Request) (err error) {
	return
}

func TestResponseRendering(t *testing.T) {
	cfg := New()

	t.Run("render without any output or error", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/", nil)
		req = Negotiate(*cfg, req)
		res := NewResponse(rec, req, *cfg)
		res.Render(nil, nil)

		if rec.Code != http.StatusOK {
			t.Fatalf("unexpected, got: %v", rec.Code)
		}
	})

	lbuf := bytes.NewBuffer(nil)
	logs := log.New(lbuf, "", 0)
	cfg = cfg.SetLogger(NewStdLogger(logs))

	t.Run("rendering an non-validation error", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/", nil)
		req = Negotiate(*cfg, req)
		res := NewResponse(rec, req, *cfg)
		res.Render(nil, errors.New("foo"))

		if rec.Code != http.StatusInternalServerError {
			t.Fatalf("unexpected, got: %v", rec.Code)
		}

		if !strings.Contains(lbuf.String(), "foo") {
			t.Fatalf("log should show error, got: %v", lbuf.String())
		}
	})

	t.Run("rendering an non-validation error", func(t *testing.T) {
		var v in1
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/", strings.NewReader(`{}`))
		req = Negotiate(*cfg, req)
		res := NewResponse(rec, req, *cfg)
		verr := res.Validate(&v)
		res.Render(verr, nil)

		if rec.Code != http.StatusOK {
			t.Fatalf("unexpected, got: %v", rec.Code)
		}
	})

	t.Run("rendering output with failing Head()", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/", nil)
		req = Negotiate(*cfg, req)
		res := NewResponse(rec, req, *cfg)

		res.Render(out1{}, nil)

		if rec.Code != http.StatusInternalServerError {
			t.Fatalf("unexpected, got: %v", rec.Code)
		}

		if res.Error() != out1Err {
			t.Fatalf("should have registered as error")
		}
	})

	cfg = New().WithEncoding(epcoding.NewJSONEncoding())

	t.Run("rendering output that cannot be encoded", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/", nil)
		req = Negotiate(*cfg, req)
		res := NewResponse(rec, req, *cfg)
		res.Render(out2{}, nil)

		if rec.Code != http.StatusInternalServerError {
			t.Fatalf("unexpected, got: %v", rec.Code)
		}
	})

	t.Run("rendering InvalidInputError", func(t *testing.T) {
		e := errors.New("foo")

		rec := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/", nil)
		req = Negotiate(*cfg, req)
		res := NewResponse(rec, req, *cfg)
		res.Render(nil, InvalidInput(e))

		if res.Error() != e {
			t.Fatalf("unexpected, got: %v", res.Error())
		}

		if rec.Code != 422 {
			t.Fatalf("unexpected, got: %v", rec.Code)
		}

		if rec.Body.String() != `{"ErrorMessage":"foo"}`+"\n" {
			t.Fatalf("unexpected, got: %v", rec.Body.String())
		}

	})
}

var _ http.ResponseWriter = &Response{}

func TestResponseWriting(t *testing.T) {
	cfg := New()
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	req = Negotiate(*cfg, req)
	res := NewResponse(rec, req, cfg)
	res.Header().Set("Foo", "BAR")
	res.Write([]byte("foo"))

	if rec.Body.String() != "foo" {
		t.Fatalf("unexpected, got: %v", rec.Body.String())
	}

	if rec.Header().Get("Foo") != "BAR" {
		t.Fatalf("should have written header, got: %v", rec.Header())
	}

	if res.state.wroteHeader != 200 {
		t.Fatalf("shoud have marked header was written, got: %v", res.state.wroteHeader)
	}
}

type out3 struct{ Bar string }

func (o *out3) Head(w http.ResponseWriter, r *http.Request) (err error) {
	return
}

func TestFullyValidResponseUsage(t *testing.T) {
	cfg := New().WithDecoding(epcoding.NewJSONDecoding()).WithEncoding(epcoding.NewJSONEncoding())

	var in in1

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", strings.NewReader(`{"Foo": "bar"}`))
	req = Negotiate(*cfg, req)
	res := NewResponse(rec, req, *cfg)
	if res.Bind(&in) {
		res.Render(&out3{strings.ToUpper(in.Foo)}, res.Validate(in))
	}

	if res.Error() != nil {
		t.Fatalf("unexpected, got: %v", res.Error())
	}

	if rec.Header().Get("Content-Type") != "application/json" {
		t.Fatalf("unexpected, got: %v", rec.Header())
	}

	if rec.Body.String() != `{"Bar":"BAR"}`+"\n" {
		t.Fatalf("unexpected, got: %v", rec.Body.String())
	}
}

type out4 struct{ FooBar string }

func (o out4) Template() string                                        { return "vt1" }
func (o out4) Head(w http.ResponseWriter, r *http.Request) (err error) { return }

func TestHTMLEncoding(t *testing.T) {
	view := template.New("root")
	view.New("vt1").Parse(`hello {{.FooBar}}`)
	view.New("error").Parse(`hello error: {{.ErrorMessage}}`)

	cfg := New().WithEncoding(epcoding.NewHTMLEncoding(view))

	t.Run("render error output with error template", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/", nil)
		req = Negotiate(*cfg, req)
		res := NewResponse(rec, req, *cfg)
		res.Render(out4{"world"}, errors.New("fail"))
		if rec.Body.String() != "hello error: Internal Server Error" {
			t.Fatalf("unexpected, got: %v", rec.Body.String())
		}
	})

	t.Run("render no error output with normal template", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/", nil)
		req = Negotiate(*cfg, req)
		res := NewResponse(rec, req, *cfg)
		res.Render(out4{"world"}, nil)
		if rec.Body.String() != "hello world" {
			t.Fatalf("unexpected, got: %v", rec.Body.String())
		}
	})
}

type out5 struct{ Bar string }

func (o *out5) Head(w http.ResponseWriter, r *http.Request) (err error) {
	w.WriteHeader(204)
	return SkipEncode
}

func TestResponseWithSkipEncode(t *testing.T) {
	cfg := New().WithDecoding(epcoding.NewJSONDecoding()).WithEncoding(epcoding.NewJSONEncoding())

	var in in1

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", strings.NewReader(`{"Foo": "bar"}`))
	req = Negotiate(*cfg, req)
	res := NewResponse(rec, req, *cfg)
	if res.Bind(&in) {
		res.Render(&out5{strings.ToUpper(in.Foo)}, res.Validate(in))
	}

	if res.Error() != nil {
		t.Fatalf("unexpected, got: %v", res.Error())
	}

	if rec.Header().Get("Content-Type") != "application/json" {
		t.Fatalf("unexpected, got: %v", rec.Header())
	}

	if rec.Body.String() != `` {
		t.Fatalf("unexpected, got: %v", rec.Body.String())
	}
}

// An input that reads itself from the request shouldn't trigger NothingBound
// validation error
func TestNonDecodingInput(t *testing.T) {
	cfg := New()
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	req = Negotiate(cfg, req)
	res := NewResponse(rec, req, cfg)

	var in in2
	if res.Bind(&in) {
		err := res.Validate(in)
		if err != nil {
			t.Fatalf("unexpected, got: %v", err)
		}

		if in.Foo != "barr" {
			t.Fatalf("unexpected, got: %v", err)
		}
	}
}

// An request without content type should use sniffing to still determine that
// it needs a JSON decoder to bind it
func TestSniffedJSONInput(t *testing.T) {
	cfg := New().WithDecoding(epcoding.NewXMLDecoding(), epcoding.NewJSONDecoding())

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", strings.NewReader(`{"Foo": "rab"}`))
	req = Negotiate(*cfg, req)
	res := NewResponse(rec, req, *cfg)

	var in in1
	ok := res.Bind(&in)
	if !ok {
		t.Fatalf("unexpected, got: %v", ok)
	}

	if in.Foo != "rab" {
		t.Fatalf("unexpected, got: %v", in.Foo)
	}
}

// TestStream of inputs to bind
func TestStreamingInput(t *testing.T) {
	cfg := New().WithDecoding(epcoding.NewXMLDecoding(), epcoding.NewJSONDecoding())

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", strings.NewReader(`{"Foo": "rab"}`+"\n"+`{"Foo": "oof"}`+"\n"+`{}`))
	req = Negotiate(*cfg, req)
	res := NewResponse(rec, req, *cfg)

	var n int
	for {
		var in in1
		if !res.Bind(&in) {
			break
		}

		n++
	}

	if n != 3 {
		t.Fatalf("unexpected, got: %v", n)
	}

	if rec.Code != 200 {
		t.Fatalf("unexpected, got: %v", rec.Code)
	}
}

type out204 struct{}

func (o out204) Head(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(204)
	return nil
}

func Test204ResponseWriting(t *testing.T) {
	cfg := New().WithEncoding(epcoding.NewJSONEncoding())
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/", nil)
	req = Negotiate(*cfg, req)
	res := NewResponse(rec, req, cfg)

	res.Render(out204{}, nil)

	if rec.Body.Len() != 0 {
		t.Fatalf("unexpected, got: %v", rec.Body.Len())
	}
}
