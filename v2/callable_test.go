package ep

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http/httptest"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

func TestNewCallable(t *testing.T) {
	for i, c := range []struct {
		fn     interface{}
		expErr error
		expTyp reflect.Type
	}{
		{nil, Err("not a func"), nil},
		{func() {}, nil, nil},
		{func(string) {}, nil, reflect.TypeOf("")},
		{func(context.Context) {}, nil, nil},
		{func(u, v string) {}, Err("first must be ctx"), nil},
		{func(myCtx, string) {}, nil, reflect.TypeOf("")},
		{func(context.Context, *string) {}, nil, reflect.TypeOf((*string)(nil))},
		{func(u, v, w string) {}, Err("at most 2 args"), nil},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			call, err := newCallable(c.fn)
			if !errors.Is(err, c.expErr) {
				t.Fatalf("expected '%v', got: '%v'", c.expErr, err)
			}

			if call == nil {
				return
			}

			if call.inpt != c.expTyp {
				t.Fatalf("expected '%v', got: '%v'", c.expTyp, call.inpt.Kind())
			}
		})
	}
}

func TestCallableInput(t *testing.T) {
	foo := "foo"

	for i, c := range []struct {
		fn      interface{}
		expZero bool
		buf     string
		expIn   interface{}
	}{
		{func() {}, true, ``, nil},
		{func(*string) {}, false, `"foo"`, &foo},
		{func(string) {}, false, `"foo"`, &foo},
		{
			func(struct{ Foo string }) {}, false, `{"Foo":"bar"}`,
			&struct{ Foo string }{"bar"},
		},
		{
			func(*struct{ Bar string }) {}, false, `{"Bar":"foo"}`,
			&struct{ Bar string }{"foo"},
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			call, err := newCallable(c.fn)
			if err != nil {
				t.Fatalf("unexpected, got: %v", err)
			}

			inp := call.Input()
			if inp == reflect.ValueOf(nil) && c.expZero {
				return
			}

			in := inp.Interface()

			dec := json.NewDecoder(strings.NewReader(c.buf))
			err = dec.Decode(in)
			if err != nil {
				t.Fatalf("unexpected, got: %v", err)
			}

			if !reflect.DeepEqual(in, c.expIn) {
				t.Fatalf("expected %v, got: %v", c.expIn, in)
			}
		})
	}
}

func TestCallableArgs(t *testing.T) {
	foo := ""

	for i, c := range []struct {
		fn      interface{}
		expArgs []reflect.Value
	}{
		{func() {}, nil},
		{func(string) {}, []reflect.Value{reflect.ValueOf("")}},
		{func(*string) {}, []reflect.Value{reflect.ValueOf(&foo)}},
		{func(myCtx) {}, []reflect.Value{reflect.ValueOf(context.Background())}},
		{func(myCtx, string) {}, []reflect.Value{reflect.ValueOf(context.Background()), reflect.ValueOf("")}},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			call, err := newCallable(c.fn)
			if err != nil {
				t.Fatalf("unexpected, got: %v", err)
			}

			inp := call.Input()
			req := httptest.NewRequest("GET", "/", nil)
			args := call.Args(req, inp)

			if fmt.Sprint(c.expArgs) != fmt.Sprint(args) {
				t.Fatalf("expected %v, got: %v", c.expArgs, args)
			}
		})
	}
}

func TestCallableCall(t *testing.T) {
	for i, c := range []struct {
		fn      interface{}
		buf     string
		expOuts []interface{}
	}{
		{func() {}, ``, []interface{}{}},
		{func(a string) string { return a }, `"foo"`, []interface{}{"foo"}},
		{func(a *string) string { return *a }, `"bar"`, []interface{}{"bar"}},
		{func(ctx context.Context, a string) string {
			if ctx == nil {
				t.Fatal("should not be nil")
			}

			return a
		}, `"rab"`, []interface{}{"rab"}},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			call, err := newCallable(c.fn)
			if err != nil {
				t.Fatalf("unexpected, got: %v", err)
			}

			inp := call.Input()

			if c.buf != "" {
				dec := json.NewDecoder(strings.NewReader(c.buf))
				err = dec.Decode(inp.Interface())
				if err != nil {
					t.Fatalf("unexpected, got: %v", err)
				}
			}

			req := httptest.NewRequest("GET", "/", nil)
			args := call.Args(req, inp)

			outs := call.Call(args)
			if !reflect.DeepEqual(outs, c.expOuts) {
				t.Fatalf("expected: %v, got: %v", c.expOuts, outs)
			}
		})
	}
}

type myCtx interface{ context.Context }

func TestCanBeAssignedContext(t *testing.T) {
	for i, c := range []struct {
		val interface{}
		exp bool
	}{
		{func(string) {}, false},
		{func(context.Context) {}, true},
		{func(myCtx) {}, true},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			act := canBeAssignedContext(reflect.TypeOf(c.val).In(0))
			if act != c.exp {
				t.Fatalf("expected: %v, got: %v", c.exp, act)
			}
		})
	}
}
