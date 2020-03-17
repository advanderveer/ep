package epcoding

import (
	"io"
	"net/http"
	"net/url"
	"strings"
)

// DefaultMaxMultipartMem is the amount of memory that is allowed to be used
// by multipart form reading
var DefaultMaxMultipartMem = int64(1024 * 1024)

// URLValuesDecoder can be implemented to automate the decoding of url.Values
type URLValuesDecoder interface {
	Decode(v interface{}, d url.Values) error
}

// FormDecoding can be used to decode multipart/form-data or urlencoded
// request bodies.
type FormDecoding struct {
	d               URLValuesDecoder
	maxMultipartMem int64
}

// NewFormDecoding inits a new decoding
func NewFormDecoding(d URLValuesDecoder) *FormDecoding { return &FormDecoding{d, 0} }

// SetMaxMultipartMemory changes the max amount of allowed memory for multipart
// decoding
func (d *FormDecoding) SetMaxMultipartMemory(nbytes int64) *FormDecoding {
	d.maxMultipartMem = nbytes
	return d
}

// Accepts returns what content-types this decoder can handle
func (d FormDecoding) Accepts() []string {
	return []string{
		"application/x-www-form-urlencoded",
		"multipart/form-data",
	}
}

// Decoder creates an actual decoder for the request
func (d FormDecoding) Decoder(r *http.Request) Decoder {
	return &FormDecoder{d.maxMultipartMem, r, d.d}
}

// FormDecoder implements decoder for forms
type FormDecoder struct {
	maxMultipartMem int64
	r               *http.Request
	d               URLValuesDecoder
}

// Decode will attempt to decode form data from the request body into
// the provided value.
func (d *FormDecoder) Decode(v interface{}) (err error) {
	if d.r == nil {
		return io.EOF
	}

	defer func() { d.r = nil }() // only decode once, no streaming

	// without a body, nothing left to do
	if d.r.ContentLength == 0 || d.r.Body == nil {
		return
	}

	// if the Request didn't provide application/x-www-form-urlencoded as the
	// content type. The req.PostForm will now be non-nil
	err = d.r.ParseForm()
	if err != nil {
		return err
	}

	mem := d.maxMultipartMem
	if mem == 0 {
		mem = DefaultMaxMultipartMem
	}

	// if the content-type indiceas it is multipart, try to decode it as such
	if strings.Contains(strings.ToLower(d.r.Header.Get("Content-Type")), "multipart/form-data") {
		err = d.r.ParseMultipartForm(mem)
		if err != nil {
			return err
		}
	}

	err = d.d.Decode(v, d.r.PostForm)
	if err != nil {
		return err
	}

	return
}
