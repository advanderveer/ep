package web

import (
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func TestRegister(t *testing.T) {
	for i, c := range []struct {
		method  string
		body    string
		ct      string
		expCode int
		expBody string
	}{
		{"GET", ``, "", 200, `form`},
		{"POST", `email=foo`, "application/x-www-form-urlencoded", 422, `forminvalid`},
		{"POST", `email=foo&password=bar`, "application/x-www-form-urlencoded", 301, ``},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(c.method, "/register", strings.NewReader(c.body))
			if c.ct != "" {
				r.Header.Set("Content-Type", c.ct)
			}

			New().ServeHTTP(w, r)

			if w.Code != c.expCode {
				t.Fatalf("expected %d, got: %d", c.expCode, w.Code)
			}

			if w.Body.String() != c.expBody {
				t.Fatalf("expected %s, got: %s", c.expBody, w.Body.String())
			}
		})
	}
}
