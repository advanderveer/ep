package epcoding

import (
	"bytes"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

func TestXMLDecoding(t *testing.T) {
	type Input struct{ Foo string }

	d := NewXMLDecoding()
	d.SetAccepts([]string{"foo/bar"})

	if !reflect.DeepEqual(d.Accepts(), []string{"foo/bar"}) {
		t.Fatalf("unepected, got: %v", d.Accepts())
	}

	req, _ := http.NewRequest("POST", "/", strings.NewReader(`<Input><Foo>bar</Foo></Input>`))

	var v Input
	dec := d.Decoder(req)
	err := dec.Decode(&v)
	if err != nil {
		t.Fatalf("unexpected, got: %v", err)
	}

	if v.Foo != "bar" {
		t.Fatalf("unexpected, got: %v", v.Foo)
	}
}

func TestXMLEncoding(t *testing.T) {
	type Output struct{ Foo string }

	e := NewXMLEncoding()
	e.SetProduces("foo")

	if e.Produces() != "foo" {
		t.Fatalf("unexpected, got: %v", e.Produces())
	}

	buf := bytes.NewBuffer(nil)
	enc := e.Encoder(buf)
	err := enc.Encode(Output{"Bar"})
	if err != nil {
		t.Fatalf("failed to encode: %v", err)
	}

	if buf.String() != `<Output><Foo>Bar</Foo></Output>` {
		t.Fatalf("unexpected, got: %v", buf.String())
	}
}
