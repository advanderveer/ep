package ep

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/advanderveer/ep/v2/coding"
)

func TestApp(t *testing.T) {
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
