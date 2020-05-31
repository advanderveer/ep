package ep

import "errors"

type Op string

type ErrorKind uint8

const (
	OtherError        ErrorKind = iota
	ServerError                 // unexpected server condition
	EmptyRequestError           // Empty request encountered
	UnacceptableError           // no encoder supports what the client accepts
	UnsupportedError            // no decoder supports the content type sent by the client
	RequestHookError            // request hook failed to run
	DecoderError                // decoder failed while decoding
	EncoderError                // encoder failed while encoding
)

type Error struct {
	msg  string
	err  error
	op   Op
	kind ErrorKind
}

func (e *Error) Unwrap() error {
	return e.err
}

func (e *Error) Is(target error) bool {
	terr, ok := target.(*Error)
	if !ok {
		return false
	}

	if terr.op != "" && terr.op != e.op {
		return false
	}

	if terr.kind != OtherError && terr.kind != e.kind {
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
		case ErrorKind:
			e.kind = at
		default:
			panic("ep: unsupported argument for building error")
		}
	}

	return e
}
