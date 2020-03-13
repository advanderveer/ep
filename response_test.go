package ep

import (
	"errors"
	"html/template"
	"net/http"
	"net/http/httptest"
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

func TestResponseBinding(t *testing.T) {
	cfg := &Config{}

	t.Run("bind without any input", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/", nil)
		req = Negotiate(*cfg, req)
		res := NewResponse(nil, req, *cfg)
		ok := res.Bind(nil)
		if !ok {
			t.Fatalf("unexpected, got: %v", ok)
		}
	})

	cfg = &Config{}
	cfg.Decoders(epcoding.NewJSONDecoding())

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

	cfg = &Config{}
	cfg.Decoders(epcoding.NewJSONDecoding())
	cfg.Encoders(epcoding.NewJSONEncoding())

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
	cfg := &Config{}

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
	})
}

type val1 struct{}

var val1err = errors.New("invalid")

func (v val1) Validate(interface{}) error { return val1err }

func TestResponseValidation(t *testing.T) {
	cfg := &Config{}

	t.Run("validate without a request body", func(t *testing.T) {
		var v in1

		req, _ := http.NewRequest("GET", "/", nil)
		req = Negotiate(*cfg, req)
		res := NewResponse(nil, req, *cfg)
		err := res.Validate(v)
		if err != errIn1 {
			t.Fatalf("unexpected, got: %v", err)
		}

		if res.Error() != errIn1 {
			t.Fatalf("unexpected, got: %v", res.Error())
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

		if res.Error() != errIn1 {
			t.Fatalf("unexpected, got: %v", res.Error())
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
	cfg := &Config{}

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

	t.Run("rendering an non-validation error", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/", nil)
		req = Negotiate(*cfg, req)
		res := NewResponse(rec, req, *cfg)
		res.Render(errors.New("foo"), nil)

		if rec.Code != http.StatusInternalServerError {
			t.Fatalf("unexpected, got: %v", rec.Code)
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

		res.Render(nil, out1{})

		if rec.Code != http.StatusInternalServerError {
			t.Fatalf("unexpected, got: %v", rec.Code)
		}

		if res.Error() != out1Err {
			t.Fatalf("should have registered as error")
		}
	})

	cfg = &Config{}
	cfg.Encoders(epcoding.NewJSONEncoding())

	t.Run("rendering output that cannot be encoded", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/", nil)
		req = Negotiate(*cfg, req)
		res := NewResponse(rec, req, *cfg)
		res.Render(nil, out2{})

		if rec.Code != http.StatusInternalServerError {
			t.Fatalf("unexpected, got: %v", rec.Code)
		}
	})

	t.Run("rendering InvalidInputError", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/", nil)
		req = Negotiate(*cfg, req)
		res := NewResponse(rec, req, *cfg)
		res.Render(InvalidInput, nil)

		if rec.Code != http.StatusOK {
			t.Fatalf("unexpected, got: %v", rec.Code)
		}
	})
}

var _ http.ResponseWriter = &Response{}

func TestResponseWriting(t *testing.T) {
	cfg := &Config{}
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	req = Negotiate(*cfg, req)
	res := NewResponse(rec, req, Config{})
	res.Header().Set("Foo", "BAR")
	res.Write([]byte("foo"))

	if rec.Body.String() != "foo" {
		t.Fatalf("unexpected, got: %v", rec.Body.String())
	}

	if rec.Header().Get("Foo") != "BAR" {
		t.Fatalf("should have written header, got: %v", rec.Header())
	}

	if res.state.wroteHeader != true {
		t.Fatalf("shoud have marked header was written, got: %v", res.state.wroteHeader)
	}
}

type out3 struct{ Bar string }

func (o *out3) Head(w http.ResponseWriter, r *http.Request) (err error) {
	return
}

func TestFullyValidResponseUsage(t *testing.T) {
	cfg := &Config{}
	cfg.Decoders(epcoding.NewJSONDecoding())
	cfg.Encoders(epcoding.NewJSONEncoding())

	var in in1

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", strings.NewReader(`{"Foo": "bar"}`))
	req = Negotiate(*cfg, req)
	res := NewResponse(rec, req, *cfg)
	if res.Bind(&in) {
		res.Render(res.Validate(in), &out3{strings.ToUpper(in.Foo)})
	}

	if res.Error() != nil {
		t.Fatalf("unexpected, got: %v", res.Error())
	}

	if rec.Body.String() != `{"Bar":"BAR"}`+"\n" {
		t.Fatalf("unexpected, got: %v", rec.Body.String())
	}
}

type out4 struct{ FooBar string }

var vt1 = template.Must(template.New("").Parse(`hello {{.FooBar}}`))
var et1 = template.Must(template.New("").Parse(`hello error: {{.ErrorMessage}}`))

func (o out4) Head(w http.ResponseWriter, r *http.Request) (err error) { return }

func TestHTMLEncoding(t *testing.T) {
	cfg := &Config{}
	cfg.Encoders(epcoding.NewHTMLEncoding(vt1, et1))

	t.Run("render error output with error template", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/", nil)
		req = Negotiate(*cfg, req)
		res := NewResponse(rec, req, *cfg)
		res.Render(errors.New("fail"), out4{"world"})
		if rec.Body.String() != "hello error: Internal Server Error" {
			t.Fatalf("unexpected, got: %v", rec.Body.String())
		}
	})

	t.Run("render no error output with normal template", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/", nil)
		req = Negotiate(*cfg, req)
		res := NewResponse(rec, req, *cfg)
		res.Render(nil, out4{"world"})
		if rec.Body.String() != "hello world" {
			t.Fatalf("unexpected, got: %v", rec.Body.String())
		}
	})
}

// An input that reads itself from the request shouldn't trigger NothingBound
// validation error
func TestNonDecodingInput(t *testing.T) {
	cfg := Config{}
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
	cfg := &Config{}
	cfg.Decoders(epcoding.NewXMLDecoding(), epcoding.NewJSONDecoding())

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
	cfg := &Config{}
	cfg.Decoders(epcoding.NewXMLDecoding(), epcoding.NewJSONDecoding())

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
}
