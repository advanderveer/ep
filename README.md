# ep
A miniature Go(lang) framework for rapid development of [http.Handlers](https://pkg.go.dev/net/http?tab=doc#Handler) 
while reducing code duplication and increasing readability. Designed to build 
both APIs and regular web applications more efficiently while keeping the 
flexibility that is expected in the Go ecosystem.

__Features:__

- Works with __any HTTP router__ that accepts the http.Handler interface
- Systematic  __error handling__ for HTTP handlers with full customization options
- Automatically encodes and decodes HTTP payloads using __content negotation__ 
- __Well tested__, benchmarked and depends only on the standard library
- Supports __streaming__ requests and responses as a first-class citizen

## Backlog V0.1
- [x] SHOULD be able to use error hooks to assert if errors are server, client or
      		 something more specific and return relevant output
- [x] SHOULD allow request hooks to return an error that is both server, or client
- [x] SHOULD handle panics in handlers as server errors
- [x] SHOULD add and test xml encoding/decoding
- [x] SHOULD add and test template encoding (text/html) and figure out how to
      		 pass the template name to the encoder
- [x] SHOULD add and test form decoding
- [x] SHOULD allow outputs to specify template method that returns the template
             type directly, not just the name
- [x] SHOULD make sure that the PrivateError hook also sets the "nosniff"
			 header like the std library
- [x] SHOULD make it clear that the build-in error hook only creates outputs for
             ep.Error errors
- [x] SHOULD should only show negotation errors when bind is actually called,
             handlers might not even want to decode. I.e a static endpoint that
             just returns a struct as output. AND:
- [x] SHOULD move negotiate to new response creation then
- [x] SHOULD not set error to nil when hooks fail to deliver an output, users
             should be able to just return an output themselves that implements
             the error interface
- [x] SHOULD return 500 when a response hook panics with nil pointer
- [x] SHOULD should be possible to prevent the body from being rendered, passing
             nil is not the correct mechanism. But the actual output should still
             be taken into account. Would be nice if it doesn't require a magic
             variable from the ep package so it doesn't need to be included
             everywhere.
             - Option 1: use reflection to check for nil value in interface
             	CON: Performance: +/-4ns
             	CON: doesn't allow empty body when value is NOT nil
             	PRO: usefull in preventing hooks calling on nil values
             - Option 2: check for a method outside of response hooks
- [x] SHOULD be able to bind empty body to allow the implementation to handle it
- [x] SHOULD also allow Read() method on input that doesn't return error
- [x] SHOULD write a basic rest example to test and apply v2
- [x] SHOULD be able to use bind with a Read implementation that reads the body
             and don't error with no-decoders
- [x] SHOULD come up with a metter name for the Empty() method
- [x] SHOULD make sure that the redirect hook behaves identical to the std lib
             redirect method. Redirect hook checks for a method to determine the
             url to redirect to, and also asserts the status method on itself
             or else takes a sensible redirect default
             - But what if two hooks trigger writeHeader? The first one takes
             precedense, so redirect hook should be put in front of 
- [x] SHOULD test the redirect hook in the rest example        
- [x] SHOULD test concurrent use of the callable     
- [x] SHOULD benchmark callable compared to non-reflection use
- [x] SHOULD properly do errors in callable logic
- [x] SHOULD rename coding and hook to `epcoding` and `ephook` 
- [ ] SHOULD rename app, should be open to users of the lib. Overloaded term
- [ ] SHOULD make sure that standard error hook returns valid html (as std lib)
       
- [ ] COULD  research if it is possible to reduce the nr of allocations when 
             decoding and encoding. It adds about 14 more allocations compared
             to the std lib variant. 
- [ ] COULD  limit the header lenght used during negotiation so it doesn't 
             allow for DDOS attacks
- [ ] COULD  make a response hook that sets cookies
- [ ] COULD  allow xml/json/form/template encoder/decoder configuration with the
             option pattern or outputs implementing a certain interface. The 
             latter is more flexible
- [ ] COULD  use a default encoding when the client specifies an accept header
             and none of the encoders match (the first configfured encoding is 
             always the default)
- [ ] COULD  lift up the error kind when nesting errors
- [ ] COULD  move the error to a separate package if it can fully replace the
             stdlib errors package
- [ ] COULD  detect if decoding should happen for an input based on whether the
             hooks have read data from the request body instead of checking a
             magic method on the input.
- [x] COULD  add some (optional) reflect sparkles for creating the handle func
             since the reflecting can be done out of the hot path. Maybe take
             inspiration from the std lib rpc package
- [x] COULD  make the response.Render() method take variadic nr of interface{}
             arguments such that exec methods can return any nr of outputs.
             response.Bind() might also be able to bind more then one input.
             Might be usefull if the endpoint has a two distict outputs with
             different templates and logic? Errors are already different outputs

- [x] WONT   return text/plain if template encoder is specified a text template. 
             each encoder should only return one type of content, we only support
             html for now
- [x] WONT   add a specific output that renders as nil, instead skipping encoding 
             instead have magic methods that are asserted.
- [x] WONT   prevent hooks from calling interface methods if the value is nil and
             causing a panic, requires reflect so maybe disable with a flag. It's
             up to the implementation to check if the value is nil

## Backlog V0.0
- [x] MUST   get kitchen example back to work
- [x] MUST   also add HTTP language negotiation
- [x] MUST   output.Head and input.Check are now optional
- [x] MUST   clean up the config and make config ergonomic 
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
- [x] SHOULD make handling errors less unwieldy, need to add a logger to see them, need to create custom outputs, 
             needs to setup html template with correct name
- [x] SHOULD when both query decoder and body decoder is configured, should be easier to protect against CSRF posts with all query params
- [x] SHOULD have type for OnErrorRender function signature
- [ ] SHOULD not overwrite content-type header if it is already set by the implementation explicitely, for example if it
             writes pdf.
- [ ] SHOULD make it clear in docs or with an error that the order of hooks is important, if one calls "writeHeader" the 
             others won't be able to change the header
- [ ] SHOULD have a clearer error when here is no html template defined for "error"
- [ ] SHOULD in general, make it easier to render some response with just a status code and a simple body (no encoding)
- [ ] SHOULD also call head hooks when not using render (but just resp.Write())
- [ ] SHOULD have a default OnErrorRender that logs it to stderr
- [ ] SHOULD use a single error value in the response instead of client and server, but instead: use an error type that 
             describes it as client or server and should change the OnErrorRender signature to allow customization of that
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
- [ ] COULD  have an hook that sets the status code based on whether the response is in a client or server error mode
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
- [x] WONT   return an error from handle as well, since that might be a common usecase. We want to motivate to move into exec function
- [x] WONT   add more logging methods to the logger to track, logging was not really used at all