package coding

import (
	"errors"
	"html/template"
	"net/http/httptest"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

type output1 struct{ Foo string }

func (o output1) Template() string { return "root" }

type output2 struct{ Foo string }

func (o output2) Template() *template.Template {
	return template.Must(template.New("root").Parse(`hello2 {{.Foo}}!`))
}

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
		{NewHTML(tmpl1), output1{"bar"}, nil, "text/html", `hello bar!`},
		{NewHTML(nil), output2{"bar"}, nil, "text/html", `hello2 bar!`},
		{NewHTML(tmpl1), struct{}{}, NoTemplateSpecified, "text/html", ``},
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
	type Input2 struct{ Foo []string }

	for i, c := range []struct {
		dec        Decoding
		in         interface{}
		body       string
		ct         string
		expErr     error
		expAccepts string
		expIn      interface{}
	}{
		{
			JSON{}, &struct{ Foo string }{}, `{"Foo": "bar"}`, "",
			nil, "application/json, application/vnd.api+json",
			&struct{ Foo string }{"bar"},
		},

		{
			XML{}, &Input{}, `<Output><Foo>bar</Foo></Output>`, "",
			nil, "application/xml, text/xml",
			&Input{"bar"},
		},

		{
			NewForm(urldec1{}), &Input2{}, `Foo=bar`, "",
			nil, "application/x-www-form-urlencoded, multipart/form-data",
			&Input2{}, // failed because content-type was not set to urlencode
		},
		{
			NewForm(urldec1{}), &Input2{}, `Foo=bar`, "application/x-www-form-urlencoded",
			nil, "application/x-www-form-urlencoded, multipart/form-data",
			&Input2{[]string{"bar"}},
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			if c.dec.Accepts() != c.expAccepts {
				t.Fatalf("expected decoder to accept '%s', got: '%s'", c.expAccepts, c.dec.Accepts())
			}

			r := httptest.NewRequest("POST", "/?Foo=bar", strings.NewReader(c.body))
			if c.ct != "" {
				r.Header.Set("Content-Type", c.ct)
			}

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
