# ep [![](https://godoc.org/github.com/advanderveer/ep?status.svg)](https://pkg.go.dev/github.com/advanderveer/ep?tab=doc)
A miniature Go(lang) framework for rapid development of [http.Handlers](https://pkg.go.dev/net/http?tab=doc#Handler) 
while reducing code duplication and increasing readability. Designed to build 
both APIs and regular web applications more efficiently while keeping the 
flexibility that is expected in the Go ecosystem.

__Features:__

- Works with __any HTTP router__ that accepts the http.Handler interface
- Supports any library for __input validation__, both system-wide or customized per endpoint
- Provides customizable __error handling__ for system errors and errors specific to your application
- Automatically encodes and decodes HTTP payloads using __content negotation__ 
- Uses __language negotiation__ os your code can use best supported language for translations
- __Well tested__, benchmarked and depends only on the standard library
- Supports __streaming__ requests and responses

## Making error handling less unwieldy
- Don't need a logger to see them
- Don't need to create custom output
- Don't need to setup a custom template with a specific name

The general premise is ok, the response will store any errors that it encounters
during the response lifecycle and renders an output of it instead of the actual
output. But having the different error types is confusing.

- Maybe just one config that turns errors into outputs. Also handle panics?
- Suitable for errors that are always handler wide: unexpected server errors and
  client errors.
- AppErrors are confusing, just for custom status code? Basically an utility to 
  render an error output with a certain status code and message.
- Question is, may the error returned from exec methods be used for rendering
  "expected" errors, like 404 or 412? If so, an handler basically has two outputs
  the "appError" output and the output type. That is confusing, it's probably
  better if the main Output type would have a "error" embedded type and an 
  "valid" embedded type, with optional rendering of it. The status would be
  a function of that with hooks.
- Question: what will be the default behaviour of the error handler func. Like
  net/html it might just log to stderr by default. But what json/html would it
  render?
- What happens if the output doesn't have a template method? maybe configure the
  template encoding with a default template so it doesn't crash with weird 
  errors?
- Or maybe just skip encoding by default?

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
- [x] MUST   find an alternative for comparing error interface values in Render: not needed, users can just retur nil
- [x] MUST   have a better way to debug unexpected error responses for development: add client and server error logging
- [x] MUST   handle panics in the handle, with the server error message rendering, should also be easy to debug
- [x] MUST   re-think usecase of rest endpoint returning error
- [x] MUST   don't write body if response is 204 or other status without a body
- [x] MUST   allow html template to accept any kind of template (interface), rename to template encoding
- [x] MUST   not server 500 status code if skipEncode is provided as an error to render
- [x] MUST   set default error template to "error.html" it is corresponds to an actual file in the most common case
- [x] SHOULD implement error hooks for handling error outputs
- [x] SHOULD implement hooks for common status responses
- [x] SHOULD implement contextual output as a hook
- [ ] SHOULD make handling errors less unwieldy, need to add a logger to see them, need to create custom outputs, needs to setup html template with correct name
- [x] SHOULD when both query decoder and body decoder is configured, should be easier to protect against CSRF posts with all query params
- [ ] SHOULD make it clear in docs or with an error that the order of hooks is important, if one calls "writeHeader" the others won't be able to change the header
- [ ] SHOULD have a clearer error when here is no html template defined for "error"
- [ ] SHOULD add more logging methods to the logger to track
- [ ] SHOULD in general, make it easier to render some response with just a status code and a simple body (no encoding)
- [ ] SHOULD also call head hooks when not using render (but just resp.Write())
- [x] SHOULD allow outputs to embed a type that will be populated with the request context
- [x] SHOULD allow configuring defaults for endpoint config
- [x] SHOULD make the Config method more ergonomic to use
- [x] SHOULD come with build-in logging support to debug client and server errors
- [x] SHOULD remove progress keeping from reader
- [x] SHOULD be able to return all kinds of app errors with status code from exec
- [x] SHOULD also make it more ergonomic to just render a 204, 404, Conflict and other common exec status responses for REST endpoints
- [x] SHOULD make it ergonomic to render output with common 2xx status codes: 201, 204
- [x] SHOULD make it ergonomic to redirect the client
- [x] SHOULD make AppError fields public
- [x] SHOULD rename "Check" on input to "Validate", way more obvious and less suprising
- [x] SHOULD SkipEncode should also work when returned directly to the render
- [x] COULD  include a more composable ways for behaviour to be added to an output: what if it redirects and sets a cookie
- [x] COULD  allow middleware to install a hook that is called just before the first byte is written to the response body 
             for middleware that needs to write a header
- [ ] COULD  check for the user that the output in a pointer value if setContext would be called
- [x] COULD  make withEncoding and withHooks consistent in naming (one with s, other without)
- [x] COULD  have package level doc summary for coding package
- [ ] COULD  not get nil pointer if status created is embedded on a nil output struct. Instead, embedding 
			 should trigger behaviour differently
- [ ] COULD  use the configuration pattern as described here: https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis
- [ ] COULD  turn most of the coding tests into table tests
- [ ] COULD  provide tooling to make endpoints extremely easy to test
- [ ] COULD  provide tooling to fuzz endpoint
- [ ] COULD  add Conf constructors for different types of endpoints: Rest, Form
- [x] COULD  make config method on endpoint optional
- [x] COULD  move per endpoint config to where Handler is called instead
- [ ] COULD  come with a nice error page for development
- [ ] COULD  rename 'epcoding' to just 'coding'
- [ ] COULD  rename coding to something else entirely, cofusing with HTTP encoding header name
- [ ] COULD  create http request interface for easier testing
- [x] COULD  remove reqProgress counter
- [ ] COULD  allow input.Read to return special error that prevents decoding
- [x] COULD  allow output.Head to return special error that prevents encoding
- [ ] COULD  better test language negotiation
- [ ] COULD  support response buffering for errors that occur halway writing the response
- [ ] COULD  allow JSON encoder configuration, i.e: indentation
- [ ] COULD  be more flexible with what content get's accepted for decoding: (i.e application/vnd.api+json should match json)
- [x] COULD  allow configuration what content-type will be written for a encoder: i.e: application/vnd.api+json
- [ ] COULD  also handle panics in the negotiation code
- [ ] COULD  assert status codes send to Error, Errorf to be in range of 400-600
- [ ] COULD  support something like this: https://github.com/mozillazg/go-httpheader on output structs
- [ ] COULD  encode response status also from output struct tags: maybe use AWS SDK approach of tagging with 'location:"header/uri/body"'
- [x] WONT   do content-encoding negotiation, complex: https://github.com/nytimes/gziphandler, deserves dedicated package
- [x] WONT   add a H/HF method for endpoints that are just the handle/exec func
- [x] WONT  return an error from handle as well, since that might be a common usecase. We want to motivate to move into exec function

