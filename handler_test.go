package ep

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
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

func action1(in handle1Input, verr error) (out handle1Output, err error) {
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

func panicHandle(res *Response, req *http.Request) {
	panic("bar")
}

func TestPanicHandling(t *testing.T) {
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/", strings.NewReader(`{"Foo": "rab"}`))
	req.Header.Set("Accept", "application/json")

	New().
		WithHook(ServerErrHook).
		WithLanguage("nl", "en-US").
		WithEncoding(epcoding.NewXMLEncoding(), epcoding.NewJSONEncoding()).
		WithDecoding(epcoding.NewXMLDecoding(), epcoding.NewJSONDecoding()).
		HandlerFunc(panicHandle).ServeHTTP(rec, req)

	if rec.Code != 500 {
		t.Fatalf("unexpected, got: %v", rec.Code)
	}

	if rec.Body.String() != `{"ErrorMessage":"Internal Server Error"}`+"\n" {
		t.Fatalf("unexpected, got: %v", rec.Body.String())
	}
}

func panic2Handle(res *Response, req *http.Request) {

	// if desired we can panic an invalid input message all the way up
	// and allow the framework to present it
	panic(Error(422, errors.New("foo")))
}

func TestPanicInvalidInputHandling(t *testing.T) {
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/", strings.NewReader(`{"Foo": "rab"}`))
	req.Header.Set("Accept", "application/json")

	New().
		WithHook(AppErrHook).
		WithLanguage("nl", "en-US").
		WithEncoding(epcoding.NewXMLEncoding(), epcoding.NewJSONEncoding()).
		WithDecoding(epcoding.NewXMLDecoding(), epcoding.NewJSONDecoding()).
		HandlerFunc(panic2Handle).ServeHTTP(rec, req)

	if rec.Code != 422 {
		t.Fatalf("unexpected, got: %v", rec.Code)
	}

	if rec.Body.String() != `{"ErrorMessage":"foo"}`+"\n" {
		t.Fatalf("unexpected, got: %v", rec.Body.String())
	}
}

func TestServerHandling(t *testing.T) {
	s := httptest.NewServer(New().
		WithLanguage("nl", "en-US").
		WithEncoding(epcoding.NewXMLEncoding(), epcoding.NewJSONEncoding()).
		WithDecoding(epcoding.NewXMLDecoding(), epcoding.NewJSONDecoding()).
		HandlerFunc(handle1))

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			resp, err := s.Client().Post(s.URL, "application/json", strings.NewReader(`{"Foo": "rab"}`))
			if err != nil {
				t.Fatalf("unexpected, got: %v", err)
			}

			body, _ := ioutil.ReadAll(resp.Body)
			if string(body) != `<handle1Output><Bar>RAB</Bar></handle1Output>` {
				t.Fatalf("unexpected, got: %v", string(body))
			}
		}()
	}

	wg.Wait()
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
