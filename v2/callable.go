package ep

import (
	"context"
	"net/http"
	"reflect"
)

// callable represents a handler that just specifies inputs and outputs
type callable struct {
	fnt  reflect.Type
	fnv  reflect.Value
	inpt reflect.Type
}

// newCallable reflects on the provided function 'f' to create the callable
func newCallable(f interface{}) (c *callable, err error) {
	c = &callable{fnv: reflect.ValueOf(f)}
	c.fnt = reflect.TypeOf(f)
	if c.fnt == nil || c.fnt.Kind() != reflect.Func {
		return nil, Err("not a func")
	}

	switch c.fnt.NumIn() {
	case 0:
	case 1:
		if canBeAssignedContext(c.fnt.In(0)) {
			break
		}

		c.inpt = c.fnt.In(0)
	case 2:
		if !canBeAssignedContext(c.fnt.In(0)) {
			return nil, Err("first must be ctx")
		}

		c.inpt = c.fnt.In(1)
	default:
		return nil, Err("at most 2 args")
	}

	return
}

// Input returns the input argument value such that it can be used to bind
func (c *callable) Input() reflect.Value {
	if c.inpt == nil {
		return reflect.ValueOf(nil)
	}

	if c.inpt.Kind() == reflect.Ptr {
		return reflect.New(c.inpt.Elem())
	}

	return reflect.New(c.inpt)
}

// inArg returns the input argument based on its type
func (c *callable) inArg(in reflect.Value) reflect.Value {
	if c.inpt.Kind() == reflect.Ptr {
		return in
	}

	return in.Elem()
}

// Setup the arguments for calling the callable
func (c *callable) Args(r *http.Request, in reflect.Value) []reflect.Value {
	args := make([]reflect.Value, c.fnt.NumIn())
	switch len(args) {
	case 0:
	case 1:
		if in == reflect.ValueOf(nil) {
			args[0] = reflect.ValueOf(r.Context())
			break
		}

		args[0] = c.inArg(in)
	case 2:
		args[0] = reflect.ValueOf(r.Context())
		args[1] = c.inArg(in)
	}

	return args
}

// Call the callable with arguments and return its outputs
func (c *callable) Call(args []reflect.Value) []interface{} {
	outs := c.fnv.Call(args)

	result := make([]interface{}, 0, len(outs))
	for _, out := range outs {
		result = append(result, out.Interface())
	}

	return result
}

// keep the ctx type, this is the idiom to get it: https://godoc.org/reflect#example-TypeOf
var ctxTyp = reflect.TypeOf((*context.Context)(nil)).Elem()

// returns whether the type implements the standard lib context interface
func canBeAssignedContext(typ reflect.Type) bool {
	if typ.Kind() != reflect.Interface {
		return false
	}

	return typ.Implements(ctxTyp)
}
