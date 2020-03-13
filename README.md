# ep
A miniature framework to reduce code duplication in writing HTTP endpoints

## Backlog
- [ ] MUST   also add HTTP language negotiation
- [ ] MUST 	 clean up the config 
- [ ] MUST   allow exec to return a InvalidAfterValidation error
- [ ] MUST   benchmark worst case sniffing and negotiation overhead
- [ ] MUST   run with race checker to check for race conditions
- [ ] MUST   fully test coding package
- [ ] MUST   test form decoder with file upload
- [ ] SHOULD be ergonomic to have translated templates as a response
- [ ] SHOULD make the Config method more ergonomic to use
- [ ] SHOULD make endpoints extremely easy to test
- [ ] SHOULD allow endpoints to be fuzzed
- [ ] COULD  come with a nice error page for development
- [ ] COULD  rename 'epcoding' to just 'coding'
- [ ] COULD  create http request interface for easier testing
- [ ] COULD  add HTTP encoding negotiation
- [ ] COULD  remove reqProgress counter
- [ ] COULD  allow input.Read to return special error that prevents decoding
- [ ] COULD  allow output.Head to return special error that prevents encoding
- [ ] COULD  better test language negotiation