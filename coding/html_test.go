package epcoding

import (
	"bytes"
	"html/template"
	"testing"
)

type tmplInput struct{ Name string }

func (in tmplInput) Template() string { return "t1" }

func TestHTMLEncoding(t *testing.T) {
	type input struct{ Name string }

	v := template.Must(template.New("").Parse(`Hello {{.Name}}`))
	template.Must(v.New("t1").Parse(`Bye, {{.Name}}`))
	e := NewHTMLEncoding(v)

	e = e.SetProduces("text/t1")
	if e.Produces() != "text/t1" {
		t.Fatalf("unexpected, got: %v", e.Produces())
	}

	t.Run("encode non templated", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)
		enc := e.Encoder(buf)

		err := enc.Encode(input{"Foo"})
		if err != nil {
			t.Fatalf("failed to encode: %v", err)
		}

		if buf.String() != "Hello Foo" {
			t.Fatalf("unexpected, got: %v", buf.String())
		}
	})

	t.Run("encode templated", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)
		enc := e.Encoder(buf)

		err := enc.Encode(tmplInput{"Foo"})
		if err != nil {
			t.Fatalf("failed to encode: %v", err)
		}

		if buf.String() != "Bye, Foo" {
			t.Fatalf("unexpected, got: %v", buf.String())
		}
	})

}
