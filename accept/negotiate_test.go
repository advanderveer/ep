// Copyright 2013 The Go Authors. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file or at
// https://developers.google.com/open-source/licenses/bsd.

package accept

import (
	"net/http"
	"testing"
)

var negotiateAcceptTest = []struct {
	s            string
	offers       []string
	defaultOffer string
	expect       string
}{
	{"text/html, */*;q=0", []string{"x/y"}, "", ""},
	{"text/html, */*", []string{"x/y"}, "", "x/y"},
	{"text/html, image/png", []string{"text/html", "image/png"}, "", "text/html"},
	{"text/html, image/png", []string{"image/png", "text/html"}, "", "image/png"},
	{"text/html, image/png; q=0.5", []string{"image/png"}, "", "image/png"},
	{"text/html, image/png; q=0.5", []string{"text/html"}, "", "text/html"},
	{"text/html, image/png; q=0.5", []string{"foo/bar"}, "", ""},
	{"text/html, image/png; q=0.5", []string{"image/png", "text/html"}, "", "text/html"},
	{"text/html, image/png; q=0.5", []string{"text/html", "image/png"}, "", "text/html"},
	{"text/html;q=0.5, image/png", []string{"image/png"}, "", "image/png"},
	{"text/html;q=0.5, image/png", []string{"text/html"}, "", "text/html"},
	{"text/html;q=0.5, image/png", []string{"image/png", "text/html"}, "", "image/png"},
	{"text/html;q=0.5, image/png", []string{"text/html", "image/png"}, "", "image/png"},
	{"image/png, image/*;q=0.5", []string{"image/jpg", "image/png"}, "", "image/png"},
	{"image/png, image/*;q=0.5", []string{"image/jpg"}, "", "image/jpg"},
	{"image/png, image/*;q=0.5", []string{"image/jpg", "image/gif"}, "", "image/jpg"},
	{"image/png, image/*", []string{"image/jpg", "image/gif"}, "", "image/jpg"},
	{"image/png, image/*", []string{"image/gif", "image/jpg"}, "", "image/gif"},
	{"image/png, image/*", []string{"image/gif", "image/png"}, "", "image/png"},
	{"image/png, image/*", []string{"image/png", "image/gif"}, "", "image/png"},
}

func TestAcceptNegotiate(t *testing.T) {
	for _, tt := range negotiateAcceptTest {
		r := &http.Request{Header: http.Header{"Accept": {tt.s}}}
		actual := Negotiate("Accept", r.Header, tt.offers, tt.defaultOffer)
		if actual != tt.expect {
			t.Errorf("NegotiateAccept(%q, %#v, %q)=%q, want %q", tt.s, tt.offers, tt.defaultOffer, actual, tt.expect)
		}
	}
}

var negotiateLanguageTests = []struct {
	s            string
	offers       []string
	defaultOffer string
	expect       string
}{
	{"en-GB,en;q=0.9,en-US;q=0.8,nl;q=0.7,it;q=0.6", []string{"xy-YX"}, "", ""},
	{"en-GB,en;q=0.9,en-US;q=0.8,nl;q=0.7,it;q=0.6", []string{"en-GB"}, "", "en-GB"},
	{"en-GB,en;q=0.9,en-US;q=0.8,nl;q=0.7,it;q=0.6", []string{"nl"}, "", "nl"},
	{"en-GB,en;q=0.9,en-US;q=0.8,nl;q=0.7,it;q=0.6", []string{"xy"}, "default", "default"},
	{"en-GB,en;q=0.9,en-US;q=0.8,nl;q=0.7,it;q=0.6", []string{"en-US", "en-GB"}, "", "en-GB"},
	{"en-GB,en;q=0.9,en-US;q=0.8,nl;q=0.7,it;q=0.6", []string{"it", "nl", "en-US"}, "", "en-US"},
}

func TestLanguageNegotiate(t *testing.T) {
	for _, tt := range negotiateLanguageTests {
		r := &http.Request{Header: http.Header{"Accept-Language": {tt.s}}}
		actual := Negotiate("Accept-Language", r.Header, tt.offers, tt.defaultOffer)
		if actual != tt.expect {
			t.Errorf("NegotiateLanguage(%q, %#v, %q)=%q, want %q", tt.s, tt.offers, tt.defaultOffer, actual, tt.expect)
		}
	}
}

var negotiateEncodingTests = []struct {
	s            string
	offers       []string
	defaultOffer string
	expect       string
}{
	{"br;q=1.0, gzip;q=0.8, *;q=0.1", []string{"xy-YX"}, "", ""},
	{"br;q=1.0, gzip;q=0.8, *;q=0.1", []string{"bogus", "gzip"}, "", "gzip"},
	{"br;q=1.0, gzip;q=0.8, *;q=0.1", []string{"br", "gzip"}, "", "br"},
}

func TestEncodingNegotiate(t *testing.T) {
	for _, tt := range negotiateEncodingTests {
		r := &http.Request{Header: http.Header{"Accept-Encoding": {tt.s}}}
		actual := Negotiate("Accept-Encoding", r.Header, tt.offers, tt.defaultOffer)
		if actual != tt.expect {
			t.Errorf("NegotiateLanguage(%q, %#v, %q)=%q, want %q", tt.s, tt.offers, tt.defaultOffer, actual, tt.expect)
		}
	}
}
