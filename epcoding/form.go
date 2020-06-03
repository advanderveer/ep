package epcoding

import (
	"io"
	"net/http"
	"net/url"
)

// URLValuesDecoder can be implemented to automate the decoding of url.Values
type URLValuesDecoder interface {
	Decode(v interface{}, d url.Values) error
}

// NewForm initializes a decoder that parses form request bodies
func NewForm(uvd URLValuesDecoder) Decoding {
	return &formDecoding{uvd}
}

type (
	formDecoding struct {
		uvd URLValuesDecoder
	}

	formDecoder struct {
		r *http.Request
		d *formDecoding
	}
)

func (d *formDecoding) Accepts() string {
	return "application/x-www-form-urlencoded, multipart/form-data"
}

func (d *formDecoding) Decoder(r *http.Request) Decoder {
	return &formDecoder{r, d}
}

func (d *formDecoder) Decode(v interface{}) (err error) {
	if d.r == nil {
		return io.EOF
	}

	defer func() { d.r = nil }() // flag as done

	// try to parse as multipart, if that fails, attempt a form parse
	err = d.r.ParseMultipartForm(32 << 20) // 32 MB
	if err != nil && err == http.ErrNotMultipart {

		// NOTE: technically speaking the ParseMultiPartForm already called
		// ParseForm for us but we don't rely on that implementation detail
		err = d.r.ParseForm()
	}

	if err != nil {
		return
	}

	// This decoder only cares about the body url values
	return d.d.uvd.Decode(v, d.r.PostForm)
}
