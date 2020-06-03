package rest

import (
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNotFound(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/bogus", strings.NewReader(`{}`))
	h := New()
	h.ServeHTTP(w, r)

	if w.Code != 404 {
		t.Fatalf("unexpected, got: %d", w.Code)
	}

	if w.Body.String() != `{"message":"Not Found"}`+"\n" {
		t.Fatalf("unexpected, got: %s", w.Body.String())
	}
}
