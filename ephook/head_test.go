package ephook

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"
)

type output2 struct{}

func (o output2) Head(h http.Header) {
	h.Set("X-Foo", "bar")
}

type output3 struct{}

func (o output3) Head(w http.ResponseWriter) {
	w.Header().Set("Bar", "foo")
}

type output4 struct{}

func (o output4) Head(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Foobar", "rab")
}

func TestHeadHook(t *testing.T) {
	for i, c := range []struct {
		out       interface{}
		expHeader http.Header
	}{
		{},
		{output2{}, http.Header{"X-Foo": []string{"bar"}}},
		{output3{}, http.Header{"Bar": []string{"foo"}}},
		{output4{}, http.Header{"Foobar": []string{"rab"}}},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			w := httptest.NewRecorder()
			Head(w, nil, c.out)

			if c.expHeader != nil && !reflect.DeepEqual(w.Header(), c.expHeader) {
				t.Fatalf("expected %#v, got: %#v", c.expHeader, w.Header())
			}
		})
	}
}
