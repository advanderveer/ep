package ep

import "errors"

type Op string

type Error struct {
	msg string
	err error
	op  Op
}

func (e *Error) Is(target error) bool {
	terr, ok := target.(*Error)
	if !ok {
		return false
	}

	if terr.op != "" && terr.op != e.op {
		return false
	}

	if terr.msg != "" && terr.msg != e.msg {
		return false
	}

	if terr.err != nil {
		return errors.Is(e.err, terr.err)
	}

	return true
}

func (e *Error) Error() string {
	if e.err == nil {
		return e.msg
	}

	if e.op == "" {
		return e.msg + ": " + e.err.Error()
	}

	return string(e.op) + ": " + e.msg + ": " + e.err.Error()
}

func Err(args ...interface{}) *Error {
	e := &Error{}
	for _, arg := range args {
		switch at := arg.(type) {
		case Op:
			e.op = at
		case error:
			e.err = at
		case string:
			e.msg = at
		default:
			panic("ep: unsupported argument for building error")
		}
	}

	return e
}
