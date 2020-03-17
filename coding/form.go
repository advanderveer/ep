package epcoding

import (
	"errors"
	"io"
	"net/http"
	"net/url"
)

type UrlValuesDecoder interface {
	Decode(v interface{}, d url.Values) error
}

type FormDecoding struct{ d UrlValuesDecoder }

func NewFormDecoding(d UrlValuesDecoder) *FormDecoding { return &FormDecoding{d} }

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
	d UrlValuesDecoder
}

func (d FormDecoder) Decode(v interface{}) (err error) {
	if d.r == nil {
		return io.EOF
	}

	defer func() { d.r = nil }() // only decode once, no streaming

	// @TODO but if what if json is first decoder and we still want to
	// decode queries: querie decoding needs to be first-class configuration

	// always parse query with this decoder
	q, err := url.ParseQuery(d.r.URL.RawQuery)
	if err != nil {
		return err
	}

	err = d.d.Decode(v, q)
	if err != nil {
		return err
	}

	// without a body, nothing left to do
	if d.r.ContentLength == 0 {
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
