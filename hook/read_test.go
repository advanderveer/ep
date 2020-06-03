package hook

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

type input1 struct{ Foo string }

func (i *input1) Read(r *http.Request) error {
	i.Foo = r.Method
	return nil
}

type input2 struct{ Bar string }

func (_ input2) Read(r *http.Request) error {
	return errors.New("foo")
}

type input3 struct{ Foo string }

func (i *input3) Read(r *http.Request) {
	i.Foo = r.Method
	return
}

func TestReadHook(t *testing.T) {
	t.Run("nil should be skipped", func(t *testing.T) {
		err := Read(nil, nil) //no input should work fine
		if err != nil {
			t.Fatalf("failed to read, got: %v", err)
		}
	})

	t.Run("should have read method", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/", nil)
		in := &input1{}
		err := Read(r, in)
		if err != nil {
			t.Fatalf("failed to read, got: %v", err)
		}

		if in.Foo != "GET" {
			t.Fatalf("unexpected, got: %v", in.Foo)
		}
	})

	t.Run("should have read method", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/", nil)
		in := &input3{}
		err := Read(r, in)
		if err != nil {
			t.Fatalf("failed to read, got: %v", err)
		}

		if in.Foo != "GET" {
			t.Fatalf("unexpected, got: %v", in.Foo)
		}
	})

	t.Run("should pass back error", func(t *testing.T) {
		err := Read(nil, input2{})
		if err == nil || err.Error() != "foo" {
			t.Fatalf("failed to read, got: %v", err)
		}
	})
}
