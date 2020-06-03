package ephook

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"
)

type output5 struct{}

func (_ output5) Status() int { return 301 }

func (_ output5) Redirect() string { return "/" }

type output6 struct{}

func (_ output6) Redirect() string { return "/foo" }

type output7 struct{}

func (_ output7) Redirect() string { return "" }

func TestRedirectHook(t *testing.T) {
	for i, c := range []struct {
		out       interface{}
		expCode   int
		expBody   string
		expHeader http.Header
	}{
		{nil, 200, ``, nil},
		{output7{}, 200, ``, nil},
		{output5{}, 301, "<a href=\"/\">Moved Permanently</a>.\n\n", http.Header{
			"Content-Type": {"text/html; charset=utf-8"},
			"Location":     {"/"},
		}},
		{output6{}, 303, "<a href=\"/foo\">See Other</a>.\n\n", http.Header{
			"Content-Type": {"text/html; charset=utf-8"},
			"Location":     {"/foo"},
		}},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/", nil)
			Redirect(w, r, c.out)

			if w.Code != c.expCode {
				t.Fatalf("expected %#v, got: %#v", c.expCode, w.Code)
			}

			if w.Body.String() != c.expBody {
				t.Fatalf("expected %#v, got: %#v", c.expBody, w.Body.String())
			}

			if c.expHeader != nil && !reflect.DeepEqual(w.Header(), c.expHeader) {
				t.Fatalf("expected %#v, got: %#v", c.expHeader, w.Header())
			}
		})
	}
}
