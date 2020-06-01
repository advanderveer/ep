package rest

import (
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func TestCreateIdea(t *testing.T) {
	for i, c := range []struct {
		body    string
		expCode int
		expBody string
		expLoc  string
	}{
		{"", 400, `{"message":"Bad Request"}` + "\n", ""},
		{"{}", 422, `{"message":"Name is empty"}` + "\n", ""},
		{`{"name": "foo"}`, 201, ``, "/ideas"},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/idea", strings.NewReader(c.body))
			New().ServeHTTP(w, r)

			if w.Code != c.expCode {
				t.Fatalf("expected %d, got: %d", c.expCode, w.Code)
			}

			if w.Body.String() != c.expBody {
				t.Fatalf("expected %s, got: %s", c.expBody, w.Body.String())
			}

			if w.Header().Get("Location") != c.expLoc {
				t.Fatalf("expected %s, got: %s", c.expLoc, w.Header().Get("Location"))
			}
		})
	}
}
