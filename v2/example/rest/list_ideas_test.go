package rest

import (
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestListIdeas(t *testing.T) {
	for i, c := range []struct {
		query   string
		expCode int
		expBody string
	}{
		{"", 200, `[{"name":"existing"}]` + "\n"},
		{"?name=foo", 200, ``},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/idea"+c.query, nil)
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
