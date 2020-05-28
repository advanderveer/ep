package accept

import (
	"testing"
)

var negotiateAcceptTest = []struct {
	asks             []string
	offers           []string
	expect           int
	expectMatchedAsk int
}{
	{[]string{"text/html, */*;q=0"}, []string{"x/y"}, -1, -1},
	{[]string{"text/html, */*"}, []string{"x/y"}, 0, 0},
	{[]string{"text/html, image/png"}, []string{"text/html", "image/png"}, 0, 0},
	{[]string{"text/html, image/png"}, []string{"image/png", "text/html"}, 0, 0},
	{[]string{"text/html, image/png; q=0.5"}, []string{"image/png"}, 0, 0},
	{[]string{"text/html, image/png; q=0.5"}, []string{"text/html"}, 0, 0},
	{[]string{"text/html, image/png; q=0.5"}, []string{"foo/bar"}, -1, -1},
	{[]string{"text/html, image/png; q=0.5"}, []string{"image/png", "text/html"}, 1, 0},
	{[]string{"text/html, image/png; q=0.5"}, []string{"text/html", "image/png"}, 0, 0},
	{[]string{"text/html;q=0.5, image/png"}, []string{"image/png"}, 0, 0},
	{[]string{"text/html;q=0.5, image/png"}, []string{"text/html"}, 0, 0},
	{[]string{"text/html;q=0.5, image/png"}, []string{"image/png", "text/html"}, 0, 0},
	{[]string{"text/html;q=0.5, image/png"}, []string{"text/html", "image/png"}, 1, 0},
	{[]string{"image/png, image/*;q=0.5"}, []string{"image/jpg", "image/png"}, 1, 0},
	{[]string{"image/png, image/*;q=0.5"}, []string{"image/jpg"}, 0, 0},
	{[]string{"image/png, image/*;q=0.5"}, []string{"image/jpg", "image/gif"}, 0, 0},
	{[]string{"image/png, image/*"}, []string{"image/jpg", "image/gif"}, 0, 0},
	{[]string{"image/png, image/*"}, []string{"image/gif", "image/jpg"}, 0, 0},
	{[]string{"image/png, image/*"}, []string{"image/gif", "image/png"}, 1, 0},
	{[]string{"image/png, image/*"}, []string{"image/png", "image/gif"}, 0, 0},
	{[]string{"image/*", "image/png"}, []string{"image/png", "image/gif"}, 0, 1},
	{[]string{"text/html;q=0.5", "image/png"}, []string{"text/html", "image/png"}, 1, 1},
}

func TestAcceptNegotiate(t *testing.T) {
	for _, tt := range negotiateAcceptTest {
		actual, matchedAsk := Negotiate(tt.asks, tt.offers)
		if actual != tt.expect {
			t.Errorf("NegotiateAccept(%v, %#v)=%q, want %q", tt.asks, tt.offers, actual, tt.expect)
		}

		if matchedAsk != tt.expectMatchedAsk {
			t.Errorf("Expected matchedAsk to be %d, got: %d", tt.expectMatchedAsk, matchedAsk)
		}

	}
}

var negotiateLanguageTests = []struct {
	s      string
	offers []string
	expect int
}{
	{"en-GB,en;q=0.9,en-US;q=0.8,nl;q=0.7,it;q=0.6", []string{"xy-YX"}, -1},
	{"en-GB,en;q=0.9,en-US;q=0.8,nl;q=0.7,it;q=0.6", []string{"en-GB"}, 0},
	{"en-GB,en;q=0.9,en-US;q=0.8,nl;q=0.7,it;q=0.6", []string{"nl"}, 0},
	{"en-GB,en;q=0.9,en-US;q=0.8,nl;q=0.7,it;q=0.6", []string{"xy"}, -1},
	{"en-GB,en;q=0.9,en-US;q=0.8,nl;q=0.7,it;q=0.6", []string{"en-US", "en-GB"}, 1},
	{"en-GB,en;q=0.9,en-US;q=0.8,nl;q=0.7,it;q=0.6", []string{"it", "nl", "en-US"}, 2},
}

func TestLanguageNegotiate(t *testing.T) {
	for _, tt := range negotiateLanguageTests {
		actual, _ := Negotiate([]string{tt.s}, tt.offers)
		if actual != tt.expect {
			t.Errorf("NegotiateLanguage(%q, %#v)=%q, want %q", tt.s, tt.offers, actual, tt.expect)
		}
	}
}

var negotiateEncodingTests = []struct {
	s      string
	offers []string
	expect int
}{
	{"br;q=1.0, gzip;q=0.8, *;q=0.1", []string{"xy-YX"}, -1},
	{"br;q=1.0, gzip;q=0.8, *;q=0.1", []string{"bogus", "gzip"}, 1},
	{"br;q=1.0, gzip;q=0.8, *;q=0.1", []string{"br", "gzip"}, 0},
}

func TestEncodingNegotiate(t *testing.T) {
	for _, tt := range negotiateEncodingTests {
		actual, _ := Negotiate([]string{tt.s}, tt.offers)
		if actual != tt.expect {
			t.Errorf("NegotiateLanguage(%q, %#v)=%q, want %q", tt.s, tt.offers, actual, tt.expect)
		}
	}
}
