package ep

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/advanderveer/ep/v2/coding"
)

func BenchmarkHandlers(b *testing.B) {
	for i, c := range []struct {
		body  string
		left  http.Handler
		right http.Handler
	}{
		{ // baseline base-case overhead between ep handle and
			body:  ``,
			left:  New().Handle(func(ResponseWriter, *http.Request) {}),
			right: http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}),
		},

		{ // overhead between reflection and no-reflection
			body:  ``,
			left:  New().Handle(func() {}),
			right: New().Handle(func(ResponseWriter, *http.Request) {}),
		},

		{ // overhead when something needs to be encoded/decoded: ~0.00323 ms
			body: `{"Foo": "bar"}`,
			left: New(
				RequestDecoding(coding.JSON{}),
				ResponseEncoding(coding.JSON{}),
			).Handle(func(in struct{ Foo string }) (out struct{ Bar string }) {
				out.Bar = in.Foo
				return
			}),
			right: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var in struct{ Foo string }
				dec := json.NewDecoder(r.Body)
				err := dec.Decode(&in)
				if err != nil {
					http.Error(w, err.Error(), 500)
					return
				}

				out := struct{ Bar string }{in.Foo}
				enc := json.NewEncoder(w)
				err = enc.Encode(out)
				if err != nil {
					http.Error(w, err.Error(), 500)
					return
				}
			}),
		},
	} {
		b.Run(strconv.Itoa(i)+"-left", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				r := httptest.NewRequest("POST", "/", strings.NewReader(c.body))
				w := httptest.NewRecorder()
				c.left.ServeHTTP(w, r)
			}
		})

		b.Run(strconv.Itoa(i)+"-right", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				r := httptest.NewRequest("POST", "/", strings.NewReader(c.body))
				w := httptest.NewRecorder()
				c.right.ServeHTTP(w, r)
			}
		})
	}
}

func TestAppHandlePanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	New().Handle(1) // first arg must be a function
}

func TestAppHandleWithReflection(t *testing.T) {
	for i, c := range []struct {
		fn      interface{}
		body    string
		expCode int
		expBody string
	}{
		{func() {}, ``, 200, ``},
		{func(a string) string { return a }, `"foo"`, 200, `"foo"` + "\n"},
		{func(a *string) string { return *a }, `"rab"`, 200, `"rab"` + "\n"},
		{func(myCtx) string { return "bar" }, ``, 200, `"bar"` + "\n"},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/", strings.NewReader(c.body))

			New(
				RequestDecoding(coding.JSON{}),
				ResponseEncoding(coding.JSON{}),
			).Handle(c.fn).ServeHTTP(w, r)

			if w.Code != c.expCode {
				t.Fatalf("expected code: %d, got: %d", c.expCode, w.Code)
			}

			if w.Body.String() != c.expBody {
				t.Fatalf("expected body: '%s', got: '%s'", c.expBody, w.Body.String())
			}
		})
	}
}

func TestAppHandleWithoutReflection(t *testing.T) {
	errHook := func(err error) (out interface{}) {
		return struct {
			Message string `json:"message"`
		}{err.Error()}
	}

	failingReqHook := func(r *http.Request, in interface{}) error {
		return errors.New("failing request hook")
	}

	always404Hook := func(w http.ResponseWriter, r *http.Request, out interface{}) {
		w.WriteHeader(404)
	}

	manualWrite := func(w ResponseWriter, r *http.Request) {
		w.Write([]byte("foo"))
	}

	renderErr := func(w ResponseWriter, r *http.Request) {
		w.Render(nil, errors.New("foo"))
	}

	panicHandle := func(w ResponseWriter, r *http.Request) {
		panic("foo")
	}

	stream := func(w ResponseWriter, r *http.Request) {
		for {
			var in struct{ Foo string }
			if !w.Bind(&in) {
				break
			}

			w.Render(in, nil)
		}
	}

	justBind := func(w ResponseWriter, r *http.Request) {
		var in struct{}
		w.Bind(&in)
	}

	for i, c := range []struct {
		method  string
		body    string
		handle  func(ResponseWriter, *http.Request)
		opt     Option
		expBody string
		expCode int
	}{
		{ // manual writing in handle shout trigger response hooks
			opt:    Options(nil, ResponseHook(always404Hook)),
			method: "GET", handle: manualWrite,
			expBody: "foo", expCode: 404,
		},

		{ // render an error should also call the hooks, and shoud encode
			opt: Options(
				ResponseEncoding(coding.JSON{}),
				ResponseHook(always404Hook),
				ErrorHook(errHook)),

			method: "GET", handle: renderErr,
			expBody: `{"message":"foo"}` + "\n", expCode: 404,
		},

		{ // steaming input to output should work as expected, with hook
			opt: Options(
				RequestDecoding(coding.JSON{}),
				ResponseEncoding(coding.JSON{}),
				ResponseHook(always404Hook)),

			method: "POST", handle: stream,
			body:    `{"Foo": "bar"}` + "\n" + `{"Foo": "rab"}`,
			expBody: `{"Foo":"bar"}` + "\n" + `{"Foo":"rab"}` + "\n",
			expCode: 404,
		},

		{ // request hook error should render as error
			opt: Options(
				RequestHook(failingReqHook),
				ResponseEncoding(coding.JSON{}),
				ResponseHook(always404Hook),
				ErrorHook(errHook)),

			method: "GET", handle: justBind,
			expBody: `{"message":"response.bind: request hook failed: failing request hook"}` + "\n", expCode: 404,
		},

		{ // panic should also render with encoder
			opt: Options(
				ResponseEncoding(coding.JSON{}),
				ErrorHook(errHook)),

			method: "GET", handle: panicHandle,
			expBody: `{"message":"foo"}` + "\n", expCode: 200,
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			r := httptest.NewRequest(c.method, "/", strings.NewReader(c.body))
			w := httptest.NewRecorder()
			New(c.opt).Handle(c.handle).ServeHTTP(w, r)

			if w.Code != c.expCode {
				t.Fatalf("expected code: %d, got: %d", c.expCode, w.Code)
			}

			if w.Body.String() != c.expBody {
				t.Fatalf("expected body: '%s', got: '%s'", c.expBody, w.Body.String())
			}
		})
	}
}
