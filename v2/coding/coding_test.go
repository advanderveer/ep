package coding

import (
	"errors"
	"net/http/httptest"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"text/template"
)

type output1 struct{ Foo string }

func (o output1) Template() string { return "root" }

func TestEncodings(t *testing.T) {
	tmpl1 := template.Must(template.New("root").Parse(`hello {{ .Foo }}!`))

	for i, c := range []struct {
		enc         Encoding
		out         interface{}
		expErr      error
		expProduces string
		expBody     string
	}{
		{JSON{}, struct{}{}, nil, "application/json", `{}` + "\n"},
		{XML{}, output1{"bar"}, nil, "application/xml", `<output1><Foo>bar</Foo></output1>`},
		{NewTemplate(tmpl1), output1{"bar"}, nil, "text/html", `hello bar!`},
		{NewTemplate(tmpl1), struct{}{}, NoTemplateSpecified, "text/html", ``},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			if c.enc.Produces() != c.expProduces {
				t.Fatalf("expected encoder to produce '%s', got: '%s'", c.expProduces, c.enc.Produces())
			}

			w := httptest.NewRecorder()
			e := c.enc.Encoder(w)
			err := e.Encode(c.out)
			if !errors.Is(err, c.expErr) {
				t.Fatalf("expected error: %#v, got: %#v", c.expErr, err)
			}

			if w.Body.String() != c.expBody {
				t.Fatalf("expected encoding body '%s', got: '%s'", c.expBody, w.Body.String())
			}
		})
	}
}

func TestDecodings(t *testing.T) {
	type Input struct{ Foo string }

	for i, c := range []struct {
		dec        Decoding
		in         interface{}
		body       string
		expErr     error
		expAccepts string
		expIn      interface{}
	}{
		{
			JSON{}, &struct{ Foo string }{}, `{"Foo": "bar"}`,
			nil, "application/json, application/vnd.api+json",
			&struct{ Foo string }{"bar"},
		},

		{
			XML{}, &Input{}, `<Output><Foo>bar</Foo></Output>`,
			nil, "application/xml, text/xml",
			&Input{"bar"},
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			if c.dec.Accepts() != c.expAccepts {
				t.Fatalf("expected decoder to accept '%s', got: '%s'", c.expAccepts, c.dec.Accepts())
			}

			r := httptest.NewRequest("POST", "/", strings.NewReader(c.body))
			d := c.dec.Decoder(r)
			err := d.Decode(c.in)
			if !errors.Is(err, c.expErr) {
				t.Fatalf("expected error: %#v, got: %#v", c.expErr, err)
			}

			if !reflect.DeepEqual(c.in, c.expIn) {
				t.Fatalf("expected in to be %#v, got: %#v", c.expIn, c.in)
			}
		})
	}
}
