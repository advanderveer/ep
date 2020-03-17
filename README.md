# ep
A miniature framework to reduce code duplication in writing HTTP endpoints

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
- [ ] MUST   fully test coding package
- [ ] MUST   be ergonomic to have translated templates as a response
- [x] SHOULD allow configuring defaults for endpoint config
- [x] SHOULD make the Config method more ergonomic to use
- [ ] SHOULD make endpoints extremely easy to test
- [ ] SHOULD allow endpoints to be fuzzed
- [ ] SHOULD come with nice logging support
- [x] SHOULD remove progress keeping from reader
- [ ] SHOULD handle panics in the handle
- [ ] Could  add Conf constructors for different types of endpoints: Rest, Form
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
- [ ] COULD  support response buffering for errors that occur halway writing
- [ ] WONT   do content-encoding negotiation, complex: https://github.com/nytimes/gziphandler
- [x] WONT   add a H/HF method for endpoints that are just the handle/exec func