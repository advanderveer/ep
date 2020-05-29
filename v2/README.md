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

## Backlog
- [x] SHOULD be able to use error hooks to assert if errors are server, client or
      		 something more specific and return relevant output
- [x] SHOULD allow request hooks to return an error that is both server, or client
- [x] SHOULD handle panics in handlers as server errors
- [x] SHOULD add and test xml encoding/decoding
- [x] SHOULD add and test template encoding (text/html) and figure out how to
      		 pass the template name to the encoder
- [ ] SHOULD add and test form decoding
- [ ] COULD  add a specific output that renders as nil, instead skipping encoding 
- [ ] COULD  allow JSON encoder/decoder configuration, i.e: indentation
- [ ] COULD  prevent hooks from calling interface methods if the value is nil and
      		 causing a panic, requires reflect so maybe disable with a flag
- [ ] COULD  use a default encoding when the client specifies an accept header
      		 and none of the encoders match (the first configfured encoding
      		 is always the default)
- [ ] COULD  lift up the error kind when nesting errors

## Hook Usecases
- [ ] SHOULD have easy to use hook that allows output to set statuscode
- [ ] SHOULD have a hook that can read mux params from the request
- [ ] SHOULD have a hook that makes redirects easy
- [ ] SHOULD support a hook that make csrf tokens available to outputs
- [ ] SHOULD support a hook that adds a localizer to each output
- [ ] SHOULD support a hook that rewrites the session cookie on each response
- [ ] SHOULD have error hooks that does sensible defaults for error rendering
		- but what to show for template render
		- but only makes sense with status response hook