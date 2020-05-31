package rest

import (
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func TestCreateIdea(t *testing.T) {
	for i, c := range []struct {
		body string
	}{
		{}, // @TODO sending an empty body to an endpoint that binds to an input should be an error
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/idea", strings.NewReader(c.body))
			New().ServeHTTP(w, r)

			println(w.Body.String())
		})
	}
}
