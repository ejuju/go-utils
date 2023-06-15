package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/ejuju/go-utils/pkg/auth"
	"github.com/ejuju/go-utils/pkg/config"
	"github.com/ejuju/go-utils/pkg/contact"
	"github.com/ejuju/go-utils/pkg/email"
	"github.com/ejuju/go-utils/pkg/logs"
	"github.com/ejuju/go-utils/pkg/media"
	"github.com/ejuju/go-utils/pkg/validation"
	"github.com/ejuju/go-utils/pkg/web"
)

func main() {
	err := newServer().run()
	if err != nil {
		panic(err)
	}
}

func init() {
	// init std/logger for external libraries using it
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

type server struct {
	h             http.Handler
	conf          conf
	logger        logs.Logger
	emailer       email.Emailer
	contactForms  contact.Forms
	uploads       media.FileStorage
	authenticator *auth.OTPAuthenticator
}

func newServer() *server {
	s := &server{}
	config.MustLoad(&s.conf,
		config.TryLoadFile("config.dev.json", config.JSONDecoder), // load config.dev.json file
		config.TryLoadFile("config.json", config.JSONDecoder),     // or load config.json file
	)
	err := s.conf.validate()
	if err != nil {
		panic(err)
	}

	switch s.conf.Env {
	case "prod":
		panic("not implemented yet")
		// s.emailer = email.NewSMTPEmailer(s.conf.SMTPEmailerConfig)
		// todo: set logger for prod
		// todo: set authenticator for prod
		// todo: set prod DB for contact form submissions
		// todo: set prod storage for file uploads
	case "dev":
		s.logger = logs.NewTextLogger(os.Stderr, logs.MustOpenLogFile(s.conf.LogFilePath))
		s.emailer = email.NewMockEmailer(os.Stderr, nil)
		s.uploads = media.NewLocalDiskStorage("uploads")
		s.contactForms = contact.MockDB{}
		s.authenticator = auth.NewOTPAuthenticator(&auth.OTPAuthenticatorConfig{
			Host:                s.conf.Host,
			ConfirmLoginRoute:   adminConfirmLoginRoute,
			SuccessfulLoginPath: adminRoute,
			CookieName:          "auth",
			Emailer:             s.emailer,
			Users:               auth.NewMockUsers("admin@local"),
			Sessions:            auth.MockSessions{},
			OTPs:                auth.MockOTPs{},
		})
	}

	// Init HTTP endpoint h
	h := web.Routes{}
	h.Handle(serveHomePage(s), web.MatchPath("/"), web.MatchMethodGet)
	h.Handle(serveContactForm(s), web.MatchPath("/contact"), web.MatchMethodPost)
	h.Handle(web.FileServer("uploads", UploadsRoute+"/"), web.MatchPathPrefix(UploadsRoute+"/"))

	h.Handle(authMiddleware(s)(serveAdminPage(s)), web.MatchPath(adminRoute), web.MatchMethodGet)
	h.Handle(authMiddleware(s)(serveAdminFileUpload(s)), web.MatchPath(adminFileUploadRoute), web.MatchMethodPost)
	h.Handle(serveLoginForm(s), web.MatchPath(adminLoginRoute), web.MatchMethodPost)
	h.Handle(serveConfirmLoginForm(s), web.MatchPath(adminConfirmLoginRoute), web.MatchMethodGet)

	h.Handle(web.ServeMonochromeFaviconPNG(nil), web.MatchPath("/favicon.ico"), web.MatchMethodGet)
	h.Handle(web.ServeSitemapXML("example.com", "/"), web.MatchPath("/sitemap.xml"), web.MatchMethodGet)
	h.Handle(serve404Page(s), web.CatchAll)
	s.h = h

	// Wrap global middleware
	s.h = web.Wrap(s.h, web.AccessLoggingMiddleware(s.logger))
	s.h = web.Wrap(s.h, web.PanicRecoveryMiddleware(s.onPanic))

	return s
}

func (s *server) run() error {
	// Log startup info
	// and listen for incoming connections (with graceful shutdown)
	s.logger.Log("starting HTTP server on port " + strconv.Itoa(s.conf.Port))
	return web.RunServer(web.NewServerWithDefaults(s.h, s.conf.Port))
}

type conf struct {
	Env               string                   `json:"env"`                  // ex: "prod" or "dev"
	Host              string                   `json:"host"`                 // ex: "example.com"
	Port              int                      `json:"port"`                 // ex: 8080
	Version           string                   `json:"version"`              // ex: "1.0.2"
	AdminEmailAddr    string                   `json:"admin_email_addr"`     // ex: "admin@example.com"
	SMTPEmailerConfig *email.SMTPEmailerConfig `json:"smtp_emailer_config"`  // see pkg/email
	LogFilePath       string                   `json:"log_file_path"`        // ex: "/tmp/my_app.log"
	UsesXForwardedFor bool                     `json:"uses_x_forwarded_for"` // use if the server is behind a reverse proxy
}

func (c *conf) validate() error {
	return validation.Validate(
		validation.CheckStringIsEither(c.Env, "prod", "dev"),
		validation.CheckNetworkPort(c.Port),
		validation.CheckEmailAddress(c.AdminEmailAddr),
		validation.CheckWhen(c.Env == "prod", validation.CheckMultiple(
			validation.CheckNotNil(c.SMTPEmailerConfig),
			func() error { return c.SMTPEmailerConfig.Validate() },
		)),
	)
}
