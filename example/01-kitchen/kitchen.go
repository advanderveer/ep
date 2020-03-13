package main

import (
	"context"
	"net/http"
	"strings"

	"github.com/advanderveer/ep"
)

type (
	KitchenInput  struct{ Foo string }
	KitchenOutput struct{ Bar string }
)

type Kitchen struct{}

func (e Kitchen) Handle(res *ep.Response, req *http.Request) {
	var in *KitchenInput
	if res.Bind(in) {
		res.Render(e.Exec(req.Context(), in, res.Validate(in)))
	}
}

func (e Kitchen) Exec(ctx context.Context, in *KitchenInput, verr error) (err error, out KitchenOutput) {
	if in == nil {
		return //nothing to do
	}

	out.Bar = strings.ToUpper(in.Foo)
	return
}
