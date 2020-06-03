package ephook

import (
	"net/http/httptest"
	"testing"
)

type output1 struct{}

func (o *output1) Status() int { return 300 }

func TestStatusHookWithNil(t *testing.T) {
	w := httptest.NewRecorder()

	var v1 *output1
	Status(w, nil, v1)

	if w.Code != 300 {
		t.Fatalf("unexpected, got: %v", w.Code)
	}
}
