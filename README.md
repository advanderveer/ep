# ep
A miniature framework to reduce code duplication in writing HTTP endpoints

## Backlog
- [x] MUST   get kitchen example back to work
- [x] MUST   also add HTTP language negotiation
- [x] MUST   output.Head and input.Check are now optional
- [ ] MUST 	 clean up the config and make config ergonomic 
- [x] MUST   allow exec to return a InvalidInput error
- [ ] MUST   benchmark worst case sniffing and negotiation overhead
- [ ] MUST   run with race checker to check for race conditions
- [ ] MUST   fully test coding package
- [ ] MUST   test file upload usecase
- [ ] SHOULD be ergonomic to have translated templates as a response
- [ ] SHOULD make the Config method more ergonomic to use
- [ ] SHOULD make endpoints extremely easy to test
- [ ] SHOULD allow endpoints to be fuzzed
- [x] SHOULD remove progress keeping from reader
- [x] COULD  make config method on endpoint optional
- [ ] COULD  come with a nice error page for development
- [ ] COULD  rename 'epcoding' to just 'coding'
- [ ] COULD  create http request interface for easier testing
- [x] COULD  remove reqProgress counter
- [ ] COULD  allow input.Read to return special error that prevents decoding
- [ ] COULD  allow output.Head to return special error that prevents encoding
- [ ] COULD  better test language negotiation
- [ ] WONT   do content-encoding negotiation, complex: https://github.com/nytimes/gziphandler