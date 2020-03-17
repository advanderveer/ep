package epcoding

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

type urldec1 struct{}

func (d urldec1) Decode(v interface{}, vals url.Values) error {
	b, _ := json.Marshal(vals)
	_ = json.Unmarshal(b, v)
	return nil
}

type urldec2 struct{}

func (d urldec2) Decode(v interface{}, vals url.Values) error {
	return errors.New("fail")
}

type formIn1 struct{ Foo []string }

func TestFormDecoding(t *testing.T) {
	urld := urldec1{}
	formd := NewFormDecoding(urld)

	if !reflect.DeepEqual(formd.Accepts(), []string{"application/x-www-form-urlencoded", "multipart/form-data"}) {
		t.Fatalf("unexpected, got: %v", formd.Accepts())
	}

	formd.SetMaxMultipartMemory(1000)
	if formd.maxMultipartMem != 1000 {
		t.Fatalf("unexpected, got: %v", formd.maxMultipartMem)
	}

	t.Run("nil body", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/", nil)
		dec := formd.Decoder(req)

		var in formIn1
		err := dec.Decode(&in)
		if err != nil {
			t.Fatalf("failed to decode: %v", err)
		}

		err = dec.Decode(&in)
		if err != io.EOF {
			t.Fatalf("should have EOF now, got: %v", err)
		}
	})

	t.Run("non multi-part body", func(t *testing.T) {
		vals := url.Values{}
		vals.Set("Foo", "bar")

		req, _ := http.NewRequest("POST", "/", strings.NewReader(vals.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		dec := formd.Decoder(req)

		var in formIn1
		err := dec.Decode(&in)
		if err != nil {
			t.Fatalf("failed to decode: %v", err)
		}

		if len(in.Foo) != 1 || in.Foo[0] != "bar" {
			t.Fatalf("unexpected, got: %v", in.Foo)
		}
	})

	t.Run("invalid non multi-part body", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/", strings.NewReader(`Foo=b%?r`))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		dec := formd.Decoder(req)

		var in formIn1
		err := dec.Decode(&in)
		if err == nil {
			t.Fatalf("unexpected: %v", err)
		}
	})

	t.Run("invalid multi-part body", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)
		mpw := multipart.NewWriter(buf)
		mpw.WriteField("Foo", "bar")
		err := mpw.Close()
		if err != nil {
			t.Fatalf("failed to close: %v", err)
		}

		req, _ := http.NewRequest("POST", "/", buf)
		req.Header.Set("Content-Type", "multipart/form-data")

		dec := formd.Decoder(req)

		var in formIn1
		err = dec.Decode(&in)
		if err == nil {
			t.Fatalf("unexpected: %v", err)
		}
	})

	t.Run("valid multi-part body", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)
		mpw := multipart.NewWriter(buf)
		mpw.WriteField("Foo", "bar")
		err := mpw.Close()
		if err != nil {
			t.Fatalf("failed to close: %v", err)
		}

		req, _ := http.NewRequest("POST", "/", buf)
		req.Header.Set("Content-Type", "multipart/form-data; boundary="+mpw.Boundary())

		dec := formd.Decoder(req)

		var in formIn1
		err = dec.Decode(&in)
		if err != nil {
			t.Fatalf("unexpected: %v", err)
		}

		if len(in.Foo) != 1 || in.Foo[0] != "bar" {
			t.Fatalf("unexpected, got: %v", in.Foo)
		}
	})

	t.Run("invalid non multi-part body", func(t *testing.T) {
		formd := NewFormDecoding(urldec2{})
		req, _ := http.NewRequest("POST", "/", strings.NewReader(`Foo=bar`))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		dec := formd.Decoder(req)

		var in formIn1
		err := dec.Decode(&in)
		if err == nil {
			t.Fatalf("unexpected: %v", err)
		}
	})
}
