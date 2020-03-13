package ep

import (
	"net/http"
	"strings"
	"testing"

	"github.com/advanderveer/ep/coding"
)

func TestNegotiate(t *testing.T) {
	jsone := epcoding.NewJSONEncoding()
	jsond := epcoding.NewJSONDecoding()

	cfg := &Config{}
	cfg.SetLanguages("it", "en-GB")
	cfg.SetEncodings(epcoding.NewXMLEncoding(), jsone)
	cfg.SetDecodings(epcoding.NewXMLDecoding(), jsond)

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
