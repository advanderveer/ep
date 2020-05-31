package ep

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http/httptest"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/advanderveer/ep/v2/coding"
)

func TestNegotiateResponseEncoder(t *testing.T) {
	for i, c := range []struct {
		accept string
		encs   []coding.Encoding
		expEnc coding.Encoder
		expErr error
		expCT  string
	}{
		{
			expErr: Err(Op("negotiateEncoder"), ServerError),
		},
		{
			encs:   []coding.Encoding{coding.JSON{}},
			expEnc: new(json.Encoder),
			expCT:  "application/json",
		},
		{
			accept: "foo/bar",
			encs:   []coding.Encoding{coding.JSON{}},
			expErr: Err(Op("negotiateEncoder"), UnacceptableError),
		},
		{
			accept: "application/json",
			encs:   []coding.Encoding{coding.JSON{}},
			expEnc: new(json.Encoder),
			expCT:  "application/json",
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			r := httptest.NewRequest("GET", "/", nil)
			r.Header.Set("Accept", c.accept)

			w := httptest.NewRecorder()

			enc, ct, err := negotiateEncoder(r, w, c.encs)
			if !errors.Is(err, c.expErr) {
				t.Fatalf("expected error %#v, got: %#v", c.expErr, err)
			}

			if ct != c.expCT {
				t.Fatalf("expected encoding ct to be %s, got: %s", c.expCT, ct)
			}

			if c.expEnc == nil && enc != c.expEnc {
				t.Fatalf("expected nil decoder, got: %#v", enc)
			} else if reflect.TypeOf(enc) != reflect.TypeOf(c.expEnc) {
				t.Fatalf("expected decoder type %T, got: %T", c.expEnc, enc)
			}
		})
	}
}

func TestNegotiateRequestDecoder(t *testing.T) {
	for i, c := range []struct {
		body   string
		ct     string
		decs   []coding.Decoding
		expDec coding.Decoder
		expErr error
	}{
		{
			"", "application/json", nil, nil,
			Err(Op("negotiateDecoder"), EmptyRequestError),
		},
		{
			"{}", "application/json ; charset=UTF-8", nil, nil,
			Err(Op("negotiateDecoder"), UnsupportedError),
		},
		{
			"{}", "foo/bar ; charset=UTF-8",
			[]coding.Decoding{coding.JSON{}}, nil,
			Err(Op("negotiateDecoder"), UnsupportedError),
		},
		{
			"{}", "application/json ; charset=UTF-8",
			[]coding.Decoding{coding.JSON{}}, &json.Decoder{}, nil,
		},
		{
			" {", "",
			[]coding.Decoding{coding.JSON{}}, &json.Decoder{}, nil,
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			r := httptest.NewRequest("POST", "/", strings.NewReader(c.body))
			if c.ct != "" {
				r.Header.Set("Content-Type", c.ct)
			}

			bdec, err := negotiateDecoder(r, c.decs)
			if !errors.Is(err, c.expErr) {
				t.Fatalf("expected error %#v, got: %#v", c.expErr, err)
			}

			if c.expDec == nil && bdec != c.expDec {
				t.Fatalf("expected nil decoder, got: %#v", bdec)
			} else if reflect.TypeOf(bdec) != reflect.TypeOf(c.expDec) {
				t.Fatalf("expected decoder type %T, got: %T", c.expDec, bdec)
			}

			all, _ := ioutil.ReadAll(r.Body)
			if string(all) != c.body {
				t.Fatalf("body not intact after negotiate, got: %v", string(all))
			}
		})
	}
}

func TestDetectContentType(t *testing.T) {
	for i, c := range []struct {
		content string
		expCT   string
	}{
		{
			" ",
			"text/plain; charset=utf-8",
		},
		{
			" {",
			"application/json; charset=utf-8",
		},
		{
			"[",
			"application/json; charset=utf-8",
		},
		{
			`"`,
			"application/json; charset=utf-8",
		},
		{
			`<?xml version="1.0" ?>`,
			"text/xml; charset=utf-8",
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			ct := detectContentType([]byte(c.content))
			if ct != c.expCT {
				t.Fatalf("expected ct to be: %v, got: %v", c.expCT, ct)
			}
		})
	}
}
