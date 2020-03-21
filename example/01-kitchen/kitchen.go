package main

import (
	"context"
	"net/http"
	"strings"

	"github.com/advanderveer/ep"
)

type (
	KitchenInput struct {
		Foo string
	}

	KitchenOutput struct {
		ep.StatusCreated
		Bar string
	}
)

func HandleKitchen(res *ep.Response, req *http.Request) {
	var in KitchenInput
	if res.Bind(&in) {
		res.Render(KitchenAction(req.Context(), in, res.Validate(in)))
	}
}

func KitchenAction(ctx context.Context, in KitchenInput, verr error) (out *KitchenOutput, err error) {
	if verr != nil {
		return nil, ep.Error(422, verr) // error is wrapped and will be used as message
	}

	if in.Foo == "bogus" {
		return nil, ep.Errorf(404, "couldn't find 'Foo': %s", in.Foo) // new error
	}

	if in.Foo == "conflict" {
		return nil, ep.Error(409) //use default http.Status message
	}

	out = &KitchenOutput{}
	out.Bar = strings.ToUpper(in.Foo)
	return
}
