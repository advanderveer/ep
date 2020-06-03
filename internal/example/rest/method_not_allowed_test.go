package rest

import (
	"net/http/httptest"
	"testing"
)

func TestMethodNotAllowed(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("PATCH", "/idea", nil)
	h := New()
	h.ServeHTTP(w, r)

	if w.Code != 405 {
		t.Fatalf("unexpected, got: %d", w.Code)
	}

	if w.Body.String() != `{"message":"Method Not Allowed"}`+"\n" {
		t.Fatalf("unexpected, got: %s", w.Body.String())
	}
}
