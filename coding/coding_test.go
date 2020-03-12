package epcoding

import (
	"net/http"
	"testing"
)

func TestEncodingNegotiation(t *testing.T) {
	t.Run("empty header empty encodings", func(t *testing.T) {
		hdr := http.Header{}
		enc := NegotiateEncoding(hdr, []Encoding{})
		if enc != nil {
			t.Fatalf("unexpected, got: %v", enc)
		}
	})

	t.Run("simple header empty encodings", func(t *testing.T) {
		hdr := http.Header{}
		hdr.Set("Accept", "application/json")

		enc := NegotiateEncoding(hdr, []Encoding{})
		if enc != nil {
			t.Fatalf("unexpected, got: %v", enc)
		}
	})

	t.Run("simple header one encoding", func(t *testing.T) {
		hdr := http.Header{}
		hdr.Set("Accept", "application/json")

		jsone := NewJSONEncoding()
		enc := NegotiateEncoding(hdr, []Encoding{jsone})
		if enc != jsone {
			t.Fatalf("unexpected, got: %v", enc)
		}
	})

	t.Run("duplicate encoding produce", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Fatalf("should panic")
			}
		}()

		hdr := http.Header{}
		hdr.Set("Accept", "application/json")

		jsone := NewJSONEncoding()
		jsonb := NewJSONEncoding()
		_ = NegotiateEncoding(hdr, []Encoding{jsone, jsonb})

	})
}

func TestDecodingNegotiation(t *testing.T) {
	t.Run("without any decoders, should be nil", func(t *testing.T) {
		hdr := http.Header{}
		dec := NegotiateDecoding(hdr, []Decoding{})
		if dec != nil {
			t.Fatalf("unexpected, got: %v", dec)
		}
	})

	t.Run("without any header, should be default", func(t *testing.T) {
		hdr := http.Header{}

		jsond := NewJSONDecoding()
		xmld := NewXMLDecoding()
		dec := NegotiateDecoding(hdr, []Decoding{xmld, jsond})
		if dec != xmld {
			t.Fatalf("unexpected, got: %v", dec)
		}
	})

	t.Run("should select json", func(t *testing.T) {
		hdr := http.Header{}
		hdr.Set("Content-Type", "application/json; charset=UTF-8")

		jsond := NewJSONDecoding()
		xmld := NewXMLDecoding()
		dec := NegotiateDecoding(hdr, []Decoding{xmld, jsond})
		if dec != jsond {
			t.Fatalf("unexpected, got: %v", dec)
		}
	})
}
