package ephook

import "net/http"

// Read is a request hook that checks if the input has a Read method with
// full control over reading the request for each input that is decoded.
func Read(r *http.Request, in interface{}) error {
	switch outt := in.(type) {
	case interface{ Read(*http.Request) error }:
		return outt.Read(r)
	case interface{ Read(*http.Request) }:
		outt.Read(r)
		return nil
	}

	return nil
}
