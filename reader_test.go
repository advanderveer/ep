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
	var progress int

	req, _ := http.NewRequest("GET", "/", strings.NewReader(`<?xml version="1.0" ?>`))
	rd := NewReader(req.Body, &progress)

	t.Run("test sniff", func(t *testing.T) {
		mt := rd.Sniff()
		if mt != "text/xml; charset=utf-8" {
			t.Fatalf("unexpected, got: %v", mt)
		}

		if progress != 0 {
			t.Fatalf("unexpected, got: %v", progress)
		}
	})

	t.Run("reading", func(t *testing.T) {
		buf := make([]byte, 10)
		rd.Read(buf)
		if progress != 10 {
			t.Fatalf("unexpected, got: %v", progress)
		}

		buf = make([]byte, 512)
		rd.Read(buf)
		if progress != 22 {
			t.Fatalf("unexpected, got: %v", progress)
		}
	})
}

func TestJSONSniff(t *testing.T) {
	for i, c := range []struct{ s string }{
		{`{}`}, {`[]`}, {`""`},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var progress int
			req, _ := http.NewRequest("GET", "/", strings.NewReader(c.s))
			r := NewReader(req.Body, &progress)
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
			var progress int
			req, _ := http.NewRequest("GET", "/", strings.NewReader(c.s))
			r := NewReader(req.Body, &progress)
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
