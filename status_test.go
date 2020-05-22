package ep

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStatusCreated(t *testing.T) {
	s := StatusCreated{}
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)

	err := StatusCreatedHook(s, rec, req)
	if err != nil {
		t.Fatalf("unexpected, got: %v", err)
	}

	if rec.Code != 201 {
		t.Fatalf("unexpected, got: %v", rec.Code)
	}

	if rec.Header().Get("Location") != "" {
		t.Fatalf("unexpected, got: %v", rec.Header())
	}

	t.Run("with location", func(t *testing.T) {
		s := StatusCreated{}
		s.SetLocation("http://google.com")

		rec := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/", nil)

		err := StatusCreatedHook(s, rec, req)
		if err != nil {
			t.Fatalf("unexpected, got: %v", err)
		}

		if rec.Header().Get("Location") != "http://google.com" {
			t.Fatalf("unexpected, got: %v", rec.Header())
		}
	})
}

func TestStatusNoContent(t *testing.T) {
	s := StatusNoContent{}
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)

	err := StatusNoContentHook(s, rec, req)
	if err != nil {
		t.Fatalf("unexpected, got: %v", err)
	}

	if rec.Code != 204 {
		t.Fatalf("unexpected, got: %v", rec.Code)
	}
}

func TestStatusRedirect(t *testing.T) {
	s := StatusRedirect{}
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)

	err := StatusRedirectHook(s, rec, req)
	if err != nil {
		t.Fatalf("unexpected, got: %v", err)
	}

	if rec.Code != 200 {
		t.Fatalf("unexpected, got: %v", rec.Code)
	}

	if rec.Header().Get("Location") != "" {
		t.Fatalf("unexpected, got: %v", rec.Header())
	}

	t.Run("set just location", func(t *testing.T) {
		s.SetRedirect("/")

		rec = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", "/", nil)
		err := StatusRedirectHook(s, rec, req)
		if err != nil {
			t.Fatalf("unexpected, got: %v", err)
		}

		if rec.Code != 303 {
			t.Fatalf("unexpected, got: %v", rec.Code)
		}

		if rec.Header().Get("Location") != "/" {
			t.Fatalf("unexpected, got: %v", rec.Header())
		}
	})

	t.Run("set with code", func(t *testing.T) {
		s.SetRedirect("/a", 307)

		rec = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", "/", nil)
		err := StatusRedirectHook(s, rec, req)
		if err != nil {
			t.Fatalf("unexpected, got: %v", err)
		}

		if rec.Code != 307 {
			t.Fatalf("unexpected, got: %v", rec.Code)
		}

		if rec.Header().Get("Location") != "/a" {
			t.Fatalf("unexpected, got: %v", rec.Header())
		}
	})

	t.Run("should panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("The code did not panic")
			}
		}()

		s.SetRedirect("/a", 307, 301)
	})
}
