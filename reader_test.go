package ep

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"strconv"
	"strings"
	"testing"
)

func TestReadProgress(t *testing.T) {

	req, _ := http.NewRequest("GET", "/", strings.NewReader(`<?xml version="1.0" ?>`))
	rd := NewReader(req.Body)

	t.Run("test sniff", func(t *testing.T) {
		mt := rd.Sniff()
		if mt != "text/xml; charset=utf-8" {
			t.Fatalf("unexpected, got: %v", mt)
		}
	})
}

func TestJSONSniff(t *testing.T) {
	for i, c := range []struct{ s string }{
		{`{}`}, {`[]`}, {`""`},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/", strings.NewReader(c.s))
			r := NewReader(req.Body)
			ct := r.Sniff()

			// if json package can decode it should detect as JSON
			var v interface{}
			err := json.Unmarshal([]byte(c.s), &v)
			if err == nil {
				if ct != "application/json; charset=utf-8" {
					t.Fatalf("unexpected, got: %v", ct)
				}
			} else {
				t.Fatalf("couldn't decode: %v", err)
			}
		})
	}
}

func TestXMLSniff(t *testing.T) {
	for i, c := range []struct{ s string }{
		{`<foo/>`},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/", strings.NewReader(c.s))
			r := NewReader(req.Body)
			ct := r.Sniff()

			// if json package can decode it should detect as JSON
			var v interface{}
			err := xml.Unmarshal([]byte(c.s), &v)
			if err == nil {
				if ct != "text/xml; charset=utf-8" {
					t.Fatalf("unexpected, got: %v", ct)
				}
			} else {
				t.Fatalf("couldn't decode: %v", err)
			}
		})
	}
}
