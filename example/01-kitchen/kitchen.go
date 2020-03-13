package main

import (
	"context"
	"net/http"
	"strings"

	"github.com/advanderveer/ep"
	"github.com/advanderveer/ep/coding"
)

type (
	KitchenInput  struct{ Foo string }
	KitchenOutput struct{ Bar string }
)

type Kitchen struct{}

func (e Kitchen) Config(c *ep.Config) {
	c.SetDecodings(NewFormDecoding())
	c.SetEncodings(epcoding.NewJSONEncoding())
}

func (e Kitchen) Handle(res *ep.Response, req *http.Request) {
	var in KitchenInput
	if res.Bind(&in) {
		res.Render(e.Exec(req.Context(), in, res.Validate(in)))
	}
}

func (e Kitchen) Exec(ctx context.Context, in KitchenInput, verr error) (err error, out KitchenOutput) {
	out.Bar = strings.ToUpper(in.Foo)
	return
}
