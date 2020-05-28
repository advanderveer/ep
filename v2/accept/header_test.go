package accept

import (
	"reflect"
	"testing"
)

var parseAcceptTests = []struct {
	s        string
	expected []AcceptSpec
}{
	{"text/html", []AcceptSpec{{"text/html", 1, 0}}},
	{"text/html; q=0", []AcceptSpec{{"text/html", 0, 0}}},
	{"text/html; q=0.0", []AcceptSpec{{"text/html", 0, 0}}},
	{"text/html; q=1", []AcceptSpec{{"text/html", 1, 0}}},
	{"text/html; q=1.0", []AcceptSpec{{"text/html", 1, 0}}},
	{"text/html; q=0.1", []AcceptSpec{{"text/html", 0.1, 0}}},
	{"text/html;q=0.1", []AcceptSpec{{"text/html", 0.1, 0}}},
	{"text/html, text/plain", []AcceptSpec{{"text/html", 1, 0}, {"text/plain", 1, 0}}},
	{"text/html; q=0.1, text/plain", []AcceptSpec{{"text/html", 0.1, 0}, {"text/plain", 1, 0}}},
	{"iso-8859-5, unicode-1-1;q=0.8,iso-8859-1", []AcceptSpec{{"iso-8859-5", 1, 0}, {"unicode-1-1", 0.8, 0}, {"iso-8859-1", 1, 0}}},
	{"iso-8859-1", []AcceptSpec{{"iso-8859-1", 1, 0}}},
	{"*", []AcceptSpec{{"*", 1, 0}}},
	{"da, en-gb;q=0.8, en;q=0.7", []AcceptSpec{{"da", 1, 0}, {"en-gb", 0.8, 0}, {"en", 0.7, 0}}},
	{"da, q, en-gb;q=0.8", []AcceptSpec{{"da", 1, 0}, {"q", 1, 0}, {"en-gb", 0.8, 0}}},
	{"image/png, image/*;q=0.5", []AcceptSpec{{"image/png", 1, 0}, {"image/*", 0.5, 0}}},

	// bad cases
	{"value1; q=0.1.2", []AcceptSpec{{"value1", 0.1, 0}}},
	{"da, en-gb;q=foo", []AcceptSpec{{"da", 1, 0}}},
}

func TestParseAccept(t *testing.T) {
	for _, tt := range parseAcceptTests {
		actual := Parse([]string{tt.s})
		if !reflect.DeepEqual(actual, tt.expected) {
			t.Errorf("Parse(h, %q)=%v, want %v", tt.s, actual, tt.expected)
		}
	}
}

func TestParseAcceptPartIndex(t *testing.T) {
	actual := Parse([]string{"text/html; q=0.1, text/plain", "val2"})
	if len(actual) != 3 {
		t.Fatalf("unexpected, got: %d", len(actual))
	}

	if actual[0].Index != 0 || actual[1].Index != 0 || actual[2].Index != 1 {
		t.Fatalf("part indexes not correct, got: %d %d %d",
			actual[0].Index, actual[1].Index, actual[2].Index)
	}
}

var parseValueAndParamsTests = []struct {
	s      string
	value  string
	params map[string]string
}{
	{`text/html`, "text/html", map[string]string{}},
	{`text/html  `, "text/html", map[string]string{}},
	{`text/html ; `, "text/html", map[string]string{}},
	{`tExt/htMl`, "text/html", map[string]string{}},
	{`tExt/htMl; fOO=";"; hellO=world`, "text/html", map[string]string{
		"hello": "world",
		"foo":   `;`,
	}},
	{`text/html; foo=bar, hello=world`, "text/html", map[string]string{"foo": "bar"}},
	{`text/html ; foo=bar `, "text/html", map[string]string{"foo": "bar"}},
	{`text/html ;foo=bar `, "text/html", map[string]string{"foo": "bar"}},
	{`text/html; foo="b\ar"`, "text/html", map[string]string{"foo": "bar"}},
	{`text/html; foo="bar\"baz\"qux"`, "text/html", map[string]string{"foo": `bar"baz"qux`}},
	{`text/html; foo="b,ar"`, "text/html", map[string]string{"foo": "b,ar"}},
	{`text/html; foo="b;ar"`, "text/html", map[string]string{"foo": "b;ar"}},
	{`text/html; FOO="bar"`, "text/html", map[string]string{"foo": "bar"}},
	{`form-data; filename="file.txt"; name=file`, "form-data", map[string]string{"filename": "file.txt", "name": "file"}},
}

func TestParseValueAndParams(t *testing.T) {
	for _, tt := range parseValueAndParamsTests {
		value, params := ParseValueAndParams(tt.s)
		if value != tt.value {
			t.Errorf("%q, value=%q, want %q", tt.s, value, tt.value)
		}
		if !reflect.DeepEqual(params, tt.params) {
			t.Errorf("%q, param=%#v, want %#v", tt.s, params, tt.params)
		}
	}
}
