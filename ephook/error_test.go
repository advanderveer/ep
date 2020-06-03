package ephook

import (
	"bytes"
	"errors"
	"log"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/advanderveer/ep"
	"github.com/advanderveer/ep/epcoding"
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
		enc     epcoding.Encoding
		err     error
		expCode int
		expBody string
	}{
		// should not create outputs for non-ep errors and nil
		{epcoding.JSON{}, nil, 200, "null\n"},
		{epcoding.JSON{}, errors.New("foo"), 200, "null\n"},

		{epcoding.JSON{}, ep.Err("foo"), 500, `{"message":"Internal Server Error"}` + "\n"},
		{epcoding.JSON{}, ep.Err(ep.DecoderError), 400, `{"message":"Bad Request"}` + "\n"},
		{epcoding.JSON{}, ep.Err(ep.UnsupportedError), 415, `{"message":"Unsupported Media Type"}` + "\n"},
		{epcoding.JSON{}, ep.Err(ep.UnacceptableError), 406, `{"message":"Not Acceptable"}` + "\n"},
		{epcoding.XML{}, ep.Err(ep.UnacceptableError), 406, `<Error><Message>Not Acceptable</Message></Error>`},
		{epcoding.NewHTML(nil), ep.Err(ep.UnacceptableError), 406, `<!doctype html><html lang="en"><head><title>Not Acceptable</title></head><body>Not Acceptable</body></html>`},
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
