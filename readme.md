# Simple Golang utilities

Idiomatic zero-dependency dead-simple Golang utility library.
Focused on utilities needed for building web services.

For protoypes, demos and learning purposes.

Features and packages:
- `auth` Authentication utilities
	- [x] Passwordless authenticator (with email magic link) (`auth.OTPAuthenticator`)
- `config` Configuration file loading utilities
	- [x] Load configuration (`config.Load`) from at least one source (can try several files)
	- [x] Support loading configuration from file (`config.TryLoadFile`) or string `config.TryLoadString`
	- [x] Support "JSON" format (`config.JSONDecoder`)
	- [ ] Support ".env" format (`config.DotEnvDecoder`)
- `contact` Contact form utilities
	- [x] Contact form DB interface (`contact.Forms`)
	- [x] Contact form struct definition (`contact.FormData`)
	- [x] Contact form validation (`contact.ParseAndValidateForm`)
	- [x] Contact form mock DB implementation (`contact.ParseAndValidateForm`)
- `email` Emailing utilities
	- [x] Email sender interface (`email.Emailer`)
	- [x] Email struct definition (`email.Email`)
	- [x] Convert email to SMTP message string (`email.SMTPMsg`)
	- [x] Email sender mock implementation (`email.MockEmailer`)
	- [x] Email sender implementation using SMTP client (`email.SMTPEmailer`)
- `logs` Logging utilities
	- [x] Logger interface definition (`logs.Logger`)
	- [x] Logger implementation for text logs (`logs.TextLogger`)
- `media` Media file storage
	- [x] Media file storage interface (`media.Storage`)
	- [x] Media file storage implementation in local FS directory (`media.LocalFileStorage`)
- `uid` Unique ID generation
	- [x] Generate ID of a certain byte length (using crypto/rand) (`uid.NewID`)
- `validation` Input validation utilities
	- [x] Validator interface (`validation.Check`)
	- [x] Check several validators sequentially (`validation.Validate`)
	- [x] Email address validator (`validation.CheckEmailAddress`)
	- [x] String match (`validation.CheckStringIs`, `validation.CheckStringIsEither`)
	- [x] String length (`validation.CheckUTF8StringMinLength`, `validation.CheckUTF8StringMaxLength`)
	- [x] Network port number (`validation.CheckNetworkPort`)
- `web` HTTP and web server utilities
	- [x] HTTP routing (`web.Routes` http.Handler implementation)
	- [x] Run HTTP server with graceful shutdown (via `web.RunServer`)
	- [x] Access logging middleware (`web.AccessLoggingMiddleware`)
	- [x] Panic recovery middleware (`web.PanicRecoveryMiddleware`)
	- [x] Permanent redirect (as http.HandlerFunc) (`web.PermanentRedirectHandler`)
	- [x] XML Sitemap generation (`web.SitemapXML`)
	- [x] Robots.txt generation with disallowed routes (`web.RobotsTXTDisallowedRoutes`)
	- [x] HTML template parsing helper (`web.MustParseHTMLTemplate`)
	- [x] HTML template rendering helper (`web.RenderHTMLTemplate`)
	- [x] Get hash from request IP address and user-agent (`web.VisitorHash`)
	- [x] Read and serve file from memory (as http.HandlerFunc) from memory (`web.ServeRaw`)
	- [x] Generate favicon PNG and serve it from memory (as http.HandlerFunc) (`web.ServeMonochromeFaviconPNG`)
	- [ ] Rate limiting middleware (`web.RateLimitingMiddleware`)
	- [ ] Timeout middleware (`web.TimeoutMiddleware`)

Todo:
- [ ] Use a Go workspace with one module for packages and one module per example folder.
- [ ] Add admin space
- [ ] Add `cms` package
- [ ] Add `cicd` package
- [ ] Add `kv` package (for implementing DB using a key-value store like BoltDB)

V2 (For later):
- [ ] Add `analytics` package
- [ ] Add `livechat` package (live support chat)
- [ ] Add `ab` package for A/B testing