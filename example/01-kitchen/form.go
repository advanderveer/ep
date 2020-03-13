package main

import (
	"io"
	"net/http"

	"github.com/advanderveer/ep/coding"
	"github.com/go-playground/form/v4"
)

type FormDecoding struct{ d *form.Decoder }

func NewFormDecoding() *FormDecoding { return &FormDecoding{form.NewDecoder()} }

func (d FormDecoding) Accepts() []string {
	return []string{"application/x-www-form-urlencoded"}
}

func (d FormDecoding) Decoder(r *http.Request) epcoding.Decoder {
	return &FormDecoder{r, d.d}
}

type FormDecoder struct {
	r *http.Request
	d *form.Decoder
}

func (d FormDecoder) Decode(v interface{}) (err error) {
	if d.r == nil {
		return io.EOF
	}

	defer func() { d.r = nil }()

	err = d.r.ParseForm()
	if err != nil {
		return err
	}

	err = d.d.Decode(v, d.r.Form)
	if err != nil {
		return err
	}

	return
}
