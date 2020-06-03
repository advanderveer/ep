package hook

import (
	"bytes"
	"errors"
	"log"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/advanderveer/ep"
	"github.com/advanderveer/ep/coding"
)

func TestPrivateErrorLogs(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	logs := log.New(buf, "", 0)

	NewStandardError(logs)(errors.New("foo"))

	if buf.String() != "foo\n" {
		t.Fatalf("should have logged, got: %v", buf.String())
	}
}

func TestPrivateErrorWithResponseHookAndEncoding(t *testing.T) {
	for i, c := range []struct {
		enc     coding.Encoding
		err     error
		expCode int
		expBody string
	}{
		// should not create outputs for non-ep errors and nil
		{coding.JSON{}, nil, 200, "null\n"},
		{coding.JSON{}, errors.New("foo"), 200, "null\n"},

		{coding.JSON{}, ep.Err("foo"), 500, `{"message":"Internal Server Error"}` + "\n"},
		{coding.JSON{}, ep.Err(ep.DecoderError), 400, `{"message":"Bad Request"}` + "\n"},
		{coding.JSON{}, ep.Err(ep.UnsupportedError), 415, `{"message":"Unsupported Media Type"}` + "\n"},
		{coding.JSON{}, ep.Err(ep.UnacceptableError), 406, `{"message":"Not Acceptable"}` + "\n"},
		{coding.XML{}, ep.Err(ep.UnacceptableError), 406, `<Error><Message>Not Acceptable</Message></Error>`},
		{coding.NewHTML(nil), ep.Err(ep.UnacceptableError), 406, `Not Acceptable`},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			w := httptest.NewRecorder()

			out := NewStandardError(nil)(c.err)
			Status(w, nil, out)

			err := c.enc.Encoder(w).Encode(out)
			if err != nil {
				t.Fatalf("unexpected, got: %v", err)
			}

			if w.Code != c.expCode {
				t.Fatalf("expected %#v, got: %#v", c.expCode, w.Code)
			}

			if w.Body.String() != c.expBody {
				t.Fatalf("expected %#v, got: %#v", c.expBody, w.Body.String())
			}
		})
	}
}
