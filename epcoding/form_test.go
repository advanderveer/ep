package epcoding

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

type urldec1 struct{}

func (_ urldec1) Decode(v interface{}, vals url.Values) error {
	b, _ := json.Marshal(vals)
	_ = json.Unmarshal(b, v)
	return nil
}

func TestFormDecodeEOF(t *testing.T) {
	r := httptest.NewRequest("POST", "/", nil)

	d := NewForm(urldec1{}).Decoder(r)
	err := d.Decode(nil)
	if err != nil {
		t.Fatalf("unexpected, got: %v", err)
	}

	err = d.Decode(nil)
	if err != io.EOF {
		t.Fatalf("unexpected, got: %v", err)
	}
}

func TestMultipartFormDecoding(t *testing.T) {
	buf := bytes.NewBuffer(nil)

	mpw := multipart.NewWriter(buf)
	mpw.WriteField("Foo", "bar")
	err := mpw.Close()
	if err != nil {
		t.Fatalf("failed to close: %v", err)
	}

	r := httptest.NewRequest("POST", "/", buf)
	r.Header.Set("Content-Type", "multipart/form-data; boundary="+mpw.Boundary())

	var v struct{ Foo []string }
	d := NewForm(urldec1{}).Decoder(r)
	err = d.Decode(&v)
	if err != nil {
		t.Fatalf("unexpected, got: %v", err)
	}

	if v.Foo[0] != "bar" {
		t.Fatalf("unexpected, got: %v", err)
	}
}

func TestFormDecodingError(t *testing.T) {
	r := httptest.NewRequest("POST", "/", strings.NewReader(`Foo=b%?r`))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	d := NewForm(urldec1{}).Decoder(r)
	err := d.Decode(nil)
	if err == nil || err.Error() != `invalid URL escape "%?r"` {
		t.Fatalf("unexpected, got: %v", err)
	}
}
