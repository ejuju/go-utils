# Simple Golang utilities

Idiomatic zero-dependency dead-simple Golang utility library.
Focused on utilities needed for building web services.

For protoypes, demos and learning purposes.

Project goal:
1. Have a set of Go packages covering various areas common to almost all web applications.
2. Keep the source code of these packages idiomatic, tiny and flexible.

Todo:
- [ ] Rate limiting middleware (`web.RateLimitingMiddleware`)
- [ ] Timeout middleware (`web.TimeoutMiddleware`)
- [ ] ??? Use a Go workspace with one module for packages and one module per example folder.
- [ ] Add admin space
- [ ] Add `cicd` package
- [ ] Add `kv` package (for implementing DB using a key-value store like BoltDB)

V2 (For later):
- [ ] Add `analytics` package
- [ ] Add `livechat` package (live support chat)
- [ ] Add `ab` package for A/B testing
- [ ] Add `shorturl` package for URL shortening and redirection with click tracking