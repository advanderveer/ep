# ep
A miniature Go(lang) framework for rapid development of [http.Handlers](https://pkg.go.dev/net/http?tab=doc#Handler) 
while reducing code duplication and increasing readability. Designed to build 
both APIs and regular web applications more efficiently while keeping the 
flexibility that is expected in the Go ecosystem.

__Features:__

- Works with __any HTTP router__ that accepts the http.Handler interface
- Provides customizable __error handling__ for system errors and errors specific to your application
- Automatically encodes and decodes HTTP payloads using __content negotation__ 
- __Well tested__, benchmarked and depends only on the standard library
- Supports __streaming__ requests and responses

## Build-in Hooks
Would be nice if ep comes with a hook that writes the status code based on a
method that the output can implement. Support both `Head(h Header) int`
or `Status() int`

With that enabled, it is possible to create an error hook that comes with
sensible defaults and logs to a provided logger. But how to support such a
build-in error with the template encoder? Maybe allow returning a template
itself. 

How about allowing handlers to return outputs that are errors (implement 
the Error() string interface). Does that work?

The PrivateError hook only handles ep errors (and panics, which are always
safe to consider to be server errors i guess). 

## Request Encoding negotiation

- Bind might be called on a non-nil input that just wants to extract query
variables in a GET request
- Bind might be called with no decoders configured but decoding of the body done 
with read methods on input structs
- Take the approach of the standard lib, ParseForm only actually touches the
body if a content-type is set. Then, what about sniffing?
- Possibly allow input also to skip decoding with an Empty() method interface

What is het meaning of calling bind to the developer. Is it just about decoding
or about making sure the hooks are called. If the content-type is set, we attempt
to find a decoder. But if we don't succeed, do we error?

Should the framework be eager to decode, or not

## Backlog
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
- [ ] SHOULD make sure that the redirect hook behaves identical to the std lib
             redirect method. Redirect hook checks for a method to determine the
             url to redirect to, and also asserts the status method on itself
             or else takes a sensible redirect default
             - But what if two hooks trigger writeHeader? The first one takes
             precedense, so redirect hook should be put in front of 
- [ ] SHOULD test the redirect hook in the rest example             

- [ ] COULD  limit the header lenght used during negotiation so it doesn't 
             allow for DDOS attacks
- [ ] COULD  make a response hook that sets cookies
- [ ] COULD  add some (optional) reflect sparkles for creating the handle func
             since the reflecting can be done out of the hot path. Maybe take
             inspiration from the std lib rpc package
- [ ] COULD  make the response.Render() method take variadic nr of interface{}
             arguments such that exec methods can return any nr of outputs.
             response.Bind() might also be able to bind more then one input.
             Might be usefull if the endpoint has a two distict outputs with
             different templates and logic? Errors are already different outputs
- [ ] COULD  add a specific output that renders as nil, instead skipping encoding 
- [ ] COULD  allow xml/json/form/template encoder/decoder configuration with the
             option pattern or outputs implementing a certain interface. The 
             latter is more flexible
- [ ] COULD  prevent hooks from calling interface methods if the value is nil and
      		 causing a panic, requires reflect so maybe disable with a flag
- [ ] COULD  use a default encoding when the client specifies an accept header
      		 and none of the encoders match (the first configfured encoding
      		 is always the default)
- [ ] COULD  lift up the error kind when nesting errors
- [ ] COULD  move the error to a separate package if it can fully replace the
             stdlib errors package
- [ ] COULD  detect if decoding should happen for an input based on whether the
             hooks have read data from the request body instead of checking a
             magic method on the input.

- [x] WONT   return text/plain if template encoder is specified a text template. 
             each encoder should only return one type of content, we only support
             html for now

## Hook Usecases
- [x] SHOULD have easy to use hook that allows output to set statuscode
- [ ] SHOULD have a hook that can read mux params from the request
- [ ] SHOULD have a hook that makes redirects easy
- [ ] SHOULD support hook that supports query decoding of the url
- [ ] SHOULD support a hook that make csrf tokens available to outputs (instead of through context)
- [ ] SHOULD support a hook that adds a localizer to each output (instead of through context)
- [ ] SHOULD support a hook that rewrites the session cookie on each response
- [x] SHOULD have error hooks that does sensible defaults for error rendering
		- but what to show for template render
		- but only makes sense with status response hook
