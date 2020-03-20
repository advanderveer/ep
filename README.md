# ep
A miniature framework to reduce code duplication in writing HTTP endpoints

## Debugging errors more easily
Specifically bad request and internal server errors
Would be nice if there are nice error pages that can be rendered
Logging could also help with that

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
- [ ] MUST   have a better way to debug unexpected error responses for development: add factories for verbose errors
- [x] MUST   handle panics in the handle, with the server error message rendering, should also be easy to debug
- [x] MUST   re-think usecase of rest endpoint returning error
- [x] MUST   don't write body if response is 204 or other status without a body
- [x] SHOULD allow configuring defaults for endpoint config
- [x] SHOULD make the Config method more ergonomic to use
- [ ] SHOULD come with build-in logging support to debug client and server errors
- [x] SHOULD remove progress keeping from reader
- [ ] SHOULD turn most of the coding tests into table tests
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
- [ ] COULD  allow output.Head to return special error that prevents encoding
- [ ] COULD  better test language negotiation
- [ ] COULD  support response buffering for errors that occur halway writing the response
- [ ] COULD  allow JSON encoder configuration, i.e: indentation
- [ ] COULD  be more flexible with what content get's accepted for decoding: (i.e application/vnd.api+json should match json)
- [x] COULD  allow configuration what content-type will be written for a encoder: i.e: application/vnd.api+json
- [x] WONT   do content-encoding negotiation, complex: https://github.com/nytimes/gziphandler, deserves dedicated package
- [x] WONT   add a H/HF method for endpoints that are just the handle/exec func
- [x] WONT  return an error from handle as well, since that might be a common usecase. We want to motivate to move into exec function