package ep

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/advanderveer/ep/coding"
)

func TestNegotiate(t *testing.T) {
	jsone := epcoding.NewJSONEncoding()
	jsond := epcoding.NewJSONDecoding()

	cfg := New().
		WithLanguage("it", "en-GB").
		WithEncoding(epcoding.NewXMLEncoding(), jsone).
		WithDecoding(epcoding.NewXMLDecoding(), jsond)

	req, _ := http.NewRequest("GET", "/", strings.NewReader(`{}`))
	req.Header.Set("Accept-Language", "en-GB,en;q=0.9,en-US;q=0.8,nl;q=0.7,it;q=0.6")
	req.Header.Set("Accept", "application/json")

	req = Negotiate(*cfg, req)

	if Language(req.Context()) != "en-GB" {
		t.Fatalf("unexpected, got: %v", Language(req.Context()))
	}

	if Encoding(req.Context()) != jsone {
		t.Fatalf("unexpected, got: %v", Encoding(req.Context()))
	}

	if Decoding(req.Context()) != jsond {
		t.Fatalf("unexpected, got: %v", Decoding(req.Context()))
	}
}

type handle1Input struct{ Foo string }
type handle1Output struct{ Bar string }

func handle1(res *Response, req *http.Request) {
	var in handle1Input
	if res.Bind(&in) {
		res.Render(action1(in, res.Validate(in)))
	}
}

func action1(in handle1Input, verr error) (err error, out handle1Output) {
	out.Bar = strings.ToUpper(in.Foo)
	return
}

func TestBasicHandler(t *testing.T) {
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/", strings.NewReader(`{"Foo": "rab"}`))
	req.Header.Set("Accept", "application/json")

	New().
		WithLanguage("nl", "en-US").
		WithEncoding(epcoding.NewXMLEncoding(), epcoding.NewJSONEncoding()).
		WithDecoding(epcoding.NewXMLDecoding(), epcoding.NewJSONDecoding()).
		HandlerFunc(handle1).ServeHTTP(rec, req)

	if rec.Body.String() != `{"Bar":"RAB"}`+"\n" {
		t.Fatalf("unexpected, got: %v", rec.Body.String())
	}
}

func BenchmarkBaseOverhead(b *testing.B) {
	h := New().HandlerFunc(func(res *Response, req *http.Request) {})
	for i := 0; i < b.N; i++ {
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/", nil)
		h.ServeHTTP(rec, req)
	}
}

func BenchmarkJSONNegotiateOverhead(b *testing.B) {
	h := New().
		WithLanguage("nl", "en-US").
		WithEncoding(epcoding.NewXMLEncoding(), epcoding.NewJSONEncoding()).
		WithDecoding(epcoding.NewXMLDecoding(), epcoding.NewJSONDecoding()).
		HandlerFunc(handle1)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {

		rec := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/", strings.NewReader(`{"Foo": "rab"}`))
		req.Header.Set("Accept", "application/json")
		h.ServeHTTP(rec, req)

		if rec.Body.String() != `{"Bar":"RAB"}`+"\n" {
			b.Fatalf("unexpected, got: %v", rec.Body.String())
		}
	}
}
