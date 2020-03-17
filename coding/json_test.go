package epcoding

import (
	"bytes"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

func TestDecoding(t *testing.T) {
	type Input struct{ Foo string }

	d := NewJSONDecoding()

	d.SetAccepts([]string{"foo"})
	if !reflect.DeepEqual(d.accepts, []string{"foo"}) {
		t.Fatalf("unexpected, got: %v", d.accepts)
	}

	req, _ := http.NewRequest("POST", "/", strings.NewReader(`{"Foo": "bbar"}`))
	dec := d.Decoder(req)

	var v Input
	err := dec.Decode(&v)
	if err != nil {
		t.Fatalf("unexpected, got: %v", err)
	}
}

func TestEncoding(t *testing.T) {
	type Output struct{ Foo string }

	e := NewJSONEncoding()
	e.SetProduces("foo")

	if e.produces != "foo" {
		t.Fatalf("unexpected, got: %v", e.produces)
	}

	buf := bytes.NewBuffer(nil)
	enc := e.Encoder(buf)
	err := enc.Encode(Output{"Bar"})
	if err != nil {
		t.Fatalf("failed to encode: %v", err)
	}

	if buf.String() != `{"Foo":"Bar"}`+"\n" {
		t.Fatalf("unexpected, got: %v", buf.String())
	}
}
