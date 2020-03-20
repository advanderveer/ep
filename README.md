# ep
A miniature framework to reduce code duplication in writing HTTP endpoints

## Render error comparison
In general the goals is to allow users to just return the validation error
to indicate that the output should be rendered normally, but return earlier

## Backlog
- [x] MUST   get kitchen example back to work
- [x] MUST   also add HTTP language negotiation
- [x] MUST   output.Head and input.Check are now optional
- [x] MUST 	 clean up the config and make config ergonomic 
- [x] MUST   allow exec to return a InvalidInput error
- [x] MUST   allow default configuration to be configured
- [x] MUST   test file upload usecase
- [x] MUST   allow (url) query decoding when request body is JSON (first decoder to implement an interface is used?)
- [x] MUST   have an form decoding that just takes an interface to do the actual decoding
- [x] MUST   benchmark worst case sniffing, negotiation and base overhead
- [x] MUST   run with race checker to check for race conditions
- [x] MUST   allow outputs to overwrite the template to use
- [x] MUST   be able to cache output templates
- [x] MUST   be ergonomic to have translated templates as a response, or other (error) customizations
- [x] MUST   fully test coding package
- [x] MUST   find an alternative for comparing error interface values in Render: not actually needed
- [ ] MUST   have a better way to debug unexpected error responses for development: add factories for verbose errors
- [ ] MUST   re-think usecase of rest endpoint returning error
- [ ] MUST   don't write body if response is 204 or other status without a body
- [x] SHOULD allow configuring defaults for endpoint config
- [x] SHOULD make the Config method more ergonomic to use
- [ ] SHOULD come with build-in logging support to debug client and server errors
- [x] SHOULD remove progress keeping from reader
- [ ] SHOULD handle panics in the handle, with the same error message rendering
- [ ] SHOULD turn most of the coding tests into table tests
- [ ] COULD  provide tooling to make endpoints extremely easy to test
- [ ] COULD  provide tooling to fuzz endpoint
- [ ] COULD  add Conf constructors for different types of endpoints: Rest, Form
- [x] COULD  make config method on endpoint optional
- [x] COULD  move per endpoint config to where Handler is called instead
- [ ] COULD  return an error from handle as well, since that might be a common usecase
- [ ] COULD  come with a nice error page for development
- [ ] COULD  rename 'epcoding' to just 'coding'
- [ ] COULD  rename coding to something else entirely, cofusing with HTTP encoding header name
- [ ] COULD  create http request interface for easier testing
- [x] COULD  remove reqProgress counter
- [ ] COULD  allow input.Read to return special error that prevents decoding
- [ ] COULD  allow output.Head to return special error that prevents encoding
- [ ] COULD  better test language negotiation
- [ ] COULD  support response buffering for errors that occur halway writing the response
- [ ] COULD  allow JSON encoder configuration, i.e: indentation
- [ ] COULD  be more flexible with what content get's accepted for decoding: (i.e application/vnd.api+json should match json)
- [ ] COULD  allow configuration what content-type will be written for a encoder: i.e: application/vnd.api+json
- [ ] WONT   do content-encoding negotiation, complex: https://github.com/nytimes/gziphandler
- [x] WONT   add a H/HF method for endpoints that are just the handle/exec func

## Bugs

```
server.go:3059: http: panic serving [::1]:56146: runtime error: comparing uncomparable type validator.ValidationErrors
goroutine 15 [running]:
net/http.(*conn).serve.func1(0xc000444000)
	/usr/local/go/src/net/http/server.go:1772 +0x139
panic(0x15e9ca0, 0xc000429520)
	/usr/local/go/src/runtime/panic.go:973 +0x396
github.com/advanderveer/ep.(*Response).Render(0xc0003db480, 0x161f300, 0xc000434ba0, 0x176d160, 0xc000434b80)
	/Users/adam/Projects/go/pkg/mod/github.com/advanderveer/ep@v0.0.2/response.go:155 +0x139
github.com/advanderveer/arc-assignment/endpoint.UserSignup.Handle(0xc0003db480, 0xc000281400)
	/Users/adam/Projects/go/src/github.com/advanderveer/arc-assignment/endpoint/user_signup.go:29 +0x15c
github.com/advanderveer/ep.Handler.ServeHTTP(0xc0003dac80, 0x176be80, 0x1ba4610, 0x1776920, 0xc00042a2a0, 0xc000281400)
	/Users/adam/Projects/go/pkg/mod/github.com/advanderveer/ep@v0.0.2/handler.go:20 +0xe1
github.com/gorilla/mux.(*Router).ServeHTTP(0xc00030c0c0, 0x1776920, 0xc00042a2a0, 0xc000280f00)
	/Users/adam/Projects/go/pkg/mod/github.com/gorilla/mux@v1.7.4/mux.go:210 +0xe2
net/http.serverHandler.ServeHTTP(0xc00042a1c0, 0x1776920, 0xc00042a2a0, 0xc000280f00)
	/usr/local/go/src/net/http/server.go:2807 +0xa3
net/http.(*conn).serve(0xc000444000, 0x1777da0, 0xc00041aec0)
	/usr/local/go/src/net/http/server.go:1895 +0x86c
created by net/http.(*Server).Serve
	/usr/local/go/src/net/http/server.go:2933 +0x35c
```

```
server.go:3059: http: panic serving [::1]:59469: ep/response: failed to render: http: request method or response status code does not allow body
goroutine 18 [running]:
net/http.(*conn).serve.func1(0xc0000b2000)
	/usr/local/go/src/net/http/server.go:1772 +0x139
panic(0x15bbda0, 0xc000094c70)
	/usr/local/go/src/runtime/panic.go:973 +0x396
github.com/advanderveer/ep.(*Response).Render(0xc0000c26c0, 0x16100a0, 0x1bad7c0, 0x0, 0x0)
	/Users/adam/Projects/go/pkg/mod/github.com/advanderveer/ep@v0.0.4/response.go:172 +0x1cd
github.com/advanderveer/arc-assignment/endpoint.DeleteIdea.Handle(0xc0000c26c0, 0xc0000d1700)
	/Users/adam/Projects/go/src/github.com/advanderveer/arc-assignment/endpoint/delete_idea.go:36 +0xe9
github.com/advanderveer/ep.Handler.ServeHTTP(0xc0003fab00, 0x1771100, 0x1bac610, 0x177be20, 0xc0000bc620, 0xc0000d1700)
	/Users/adam/Projects/go/pkg/mod/github.com/advanderveer/ep@v0.0.4/handler.go:20 +0xe1
github.com/gorilla/mux.(*Router).ServeHTTP(0xc00032c0c0, 0x177be20, 0xc0000bc620, 0xc0000d1200)
	/Users/adam/Projects/go/pkg/mod/github.com/gorilla/mux@v1.7.4/mux.go:210 +0xe2
net/http.serverHandler.ServeHTTP(0xc0004421c0, 0x177be20, 0xc0000bc620, 0xc0000d1200)
	/usr/local/go/src/net/http/server.go:2807 +0xa3
net/http.(*conn).serve(0xc0000b2000, 0x177d2a0, 0xc00009c1c0)
	/usr/local/go/src/net/http/server.go:1895 +0x86c
created by net/http.(*Server).Serve
	/usr/local/go/src/net/http/server.go:2933 +0x35c
```