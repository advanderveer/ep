package epcoding

import (
	"errors"
	"io"
	"net/http"
	"net/url"
)

type URLValuesDecoder interface {
	Decode(v interface{}, d url.Values) error
}

type FormDecoding struct{ d URLValuesDecoder }

func NewFormDecoding(d URLValuesDecoder) *FormDecoding { return &FormDecoding{d} }

func (d FormDecoding) Accepts() []string {
	return []string{
		"application/x-www-form-urlencoded",
		"multipart/form-data",
	}
}

func (d FormDecoding) Decoder(r *http.Request) Decoder {
	return &FormDecoder{r, d.d}
}

type FormDecoder struct {
	r *http.Request
	d URLValuesDecoder
}

func (d FormDecoder) Decode(v interface{}) (err error) {
	if d.r == nil {
		return io.EOF
	}

	defer func() { d.r = nil }() // only decode once, no streaming

	// without a body, nothing left to do
	if d.r.ContentLength == 0 || d.r.Body == nil {
		return
	}

	// try to just parse multipart, internally it calls ParseForm first, so when
	// the encoding is not multi-part decoding will still work
	_ = d.r.ParseMultipartForm(1024 * 10)
	if d.r.PostForm == nil {
		return errors.New("invalid multipart or urlencoded post data")
	}

	err = d.d.Decode(v, d.r.PostForm)
	if err != nil {
		return err
	}

	return
}
