package web

import (
	"net/http/httptest"
	"testing"
)

func TestNotFound(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/bogus", nil)
	h := New()
	h.ServeHTTP(w, r)

	if w.Code != 404 {
		t.Fatalf("unexpected, got: %d", w.Code)
	}

	if w.Body.String() != `Not Found` {
		t.Fatalf("unexpected, got: %s", w.Body.String())
	}
}
