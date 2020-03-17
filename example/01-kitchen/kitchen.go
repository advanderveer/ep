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

func HandleKitchen(res *ep.Response, req *http.Request) {
	var in KitchenInput
	if res.Bind(&in) {
		res.Render(KitchenAction(req.Context(), in, res.Validate(in)))
	}
}

func KitchenAction(ctx context.Context, in KitchenInput, verr error) (out KitchenOutput, err error) {
	out.Bar = strings.ToUpper(in.Foo)
	return
}
