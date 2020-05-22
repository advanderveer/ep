package ep

import (
	"net/url"
	"reflect"
	"testing"

	"github.com/advanderveer/ep/coding"
)

type val2 struct{}

func (v val2) Validate(interface{}) error { return nil }

type qdec struct{}

func (d qdec) Decode(interface{}, url.Values) error { return nil }

func TestConf(t *testing.T) {

	jsone1 := epcoding.NewJSONEncoding()
	jsond1 := epcoding.NewJSONDecoding()
	xmle1 := epcoding.NewXMLEncoding()
	xmld1 := epcoding.NewXMLDecoding()
	v1 := val2{}
	// cef1 := func(error) Output { return nil }
	// aef1 := func(*AppError) Output { return nil }
	qdec1 := qdec{}

	t.Run("encoding methods", func(t *testing.T) {
		c1 := New()
		c2 := c1.WithEncoding(jsone1)
		c2 = c2.WithEncoding(xmle1)
		if len(c2.Encodings()) != 2 || c2.Encodings()[1] != xmle1 {
			t.Fatalf("unexpected, got: %v", c2.Encodings())
		}

		c2 = c2.SetEncodings(xmle1)
		if len(c2.Encodings()) != 1 || c2.Encodings()[0] != xmle1 {
			t.Fatalf("unexpected, got: %v", c2.Encodings())
		}

		if !reflect.DeepEqual(c1.Encodings(), c2.Encodings()) {
			t.Fatalf("should have edited original")
		}
	})

	t.Run("decoding methods", func(t *testing.T) {
		c1 := New()
		c2 := c1.WithDecoding(jsond1)
		c2 = c2.WithDecoding(xmld1)
		if len(c2.Decodings()) != 2 || c2.Decodings()[1] != xmld1 {
			t.Fatalf("unexpected, got: %v", c2.Decodings())
		}

		c2 = c2.SetDecodings(xmld1)
		if len(c2.Decodings()) != 1 || c2.Decodings()[0] != xmld1 {
			t.Fatalf("unexpected, got: %v", c2.Decodings())
		}

		if !reflect.DeepEqual(c1.Decodings(), c2.Decodings()) {
			t.Fatalf("should have edited original")
		}
	})

	t.Run("language methods", func(t *testing.T) {
		c1 := New()
		c2 := c1.WithLanguage("en-GB")
		c2 = c2.WithLanguage("en-US")
		if len(c2.Languages()) != 2 || c2.Languages()[1] != "en-US" {
			t.Fatalf("unexpected, got: %v", c2.Languages())
		}

		c2 = c2.SetLanguages("en-US")
		if len(c2.Languages()) != 1 || c2.Languages()[0] != "en-US" {
			t.Fatalf("unexpected, got: %v", c2.Languages())
		}

		if !reflect.DeepEqual(c1.Languages(), c2.Languages()) {
			t.Fatalf("should have edited original")
		}
	})

	t.Run("should copy", func(t *testing.T) {
		c1 := (New()).WithEncoding(jsone1)
		c2 := c1.Copy().WithEncoding(xmle1)

		if len(c2.Encodings()) != 2 || len(c1.Encodings()) != 1 {
			t.Fatalf("should be distinct confs")
		}
	})

	t.Run("set validator", func(t *testing.T) {
		c1 := (New()).SetValidator(v1)

		if c1.Validator() != v1 {
			t.Fatalf("unexpected, got: %v", c1.Validator())
		}
	})

	t.Run("set validator", func(t *testing.T) {
		c1 := (New()).SetQueryDecoder(qdec1)

		if c1.QueryDecoder() != qdec1 {
			t.Fatalf("unexpected, got: %v", c1.QueryDecoder())
		}
	})
}
