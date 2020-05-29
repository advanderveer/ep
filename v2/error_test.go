package ep

import (
	"errors"
	"strconv"
	"testing"
)

func TestErrorBuilding(t *testing.T) {
	e1 := Err("my message")
	if e1.Error() != "my message" {
		t.Fatalf("unexpected, got: %v", e1.Error())
	}

	e2 := Err("my other message", e1)
	if e2.Error() != "my other message: my message" {
		t.Fatalf("unexpected, got: %v", e2.Error())
	}

	e3 := Err("my other message", e1, Op("my op"))
	if e3.Error() != "my op: my other message: my message" {
		t.Fatalf("unexpected, got: %v", e3.Error())
	}

	e4 := Err("my other message", e1, Op("my op"), OtherError)
	if e4.Error() != "my op: my other message: my message" {
		t.Fatalf("unexpected, got: %v", e4.Error())
	}
}

func TestErrorUnwrapping(t *testing.T) {
	e1 := errors.New("foo")
	e2 := Err(e1, "foo")
	if !errors.Is(e2, e1) {
		t.Fatalf("e2 should be a e1")
	}

	e3 := Err(Op("my.op"))
	e4 := Err(e3, "foo")
	if !errors.Is(e4, Err(Op("my.op"))) {
		t.Fatalf("e4 should be a e3")
	}

	e5 := Err(Op("my.op"))
	e6 := Err(e5, "foo")
	if errors.Is(e6, Err(Op("my.other.op"))) {
		t.Fatalf("e6 should not be a e5")
	}
}

func TestErrorBuildingPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	Err(1)
}

func TestErrorMatching(t *testing.T) {
	for i, c := range []struct {
		err    error
		target error
		expect bool
	}{
		{Err("foo"), Err("bar"), false},
		{Err("foo"), Err("foo"), true},
		{Err("foo"), Err(Op("op"), "foo"), false},
		{Err("foo", Op("op")), Err(Op("op"), "foo"), true},
		{Err("foo", Op("op2")), Err(Op("op"), "foo"), false},
		{Err(Err("foo")), Err(Err("foo")), true},
		{Err(Err("foo")), Err(Err("bar")), false},
		{Err("foo"), errors.New("foo"), false},
		{Err(ErrorKind(1)), Err(ErrorKind(1)), true},
		{Err(ErrorKind(155)), Err(ErrorKind(2)), false},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			actual := errors.Is(c.err, c.target)
			if actual != c.expect {
				t.Errorf("errors.Is(%s, %s)=%v, got: %v", c.err, c.target, c.expect, actual)
			}
		})
	}

}
