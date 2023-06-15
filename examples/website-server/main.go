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

type server struct {
	h             http.Handler
	conf          conf
	logger        logs.Logger
	emailer       email.Emailer
	contactForms  contact.Forms
	uploads       media.FileStorage
	authenticator *auth.OTPAuthenticator
}

type conf struct {
	Env               string                   `json:"env"`                 // ex: "prod" or "dev"
	Host              string                   `json:"host"`                // ex: "example.com"
	Port              int                      `json:"port"`                // ex: 8080
	Version           string                   `json:"version"`             // ex: "1.0.2"
	AdminEmailAddr    string                   `json:"admin_email_addr"`    // ex: "admin@example.com"
	SMTPEmailerConfig *email.SMTPEmailerConfig `json:"smtp_emailer_config"` // see pkg/email
	LogFilePath       string                   `json:"log_file_path"`       // ex: "/tmp/my_app.log"
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

const (
	contactRoute           = "/contact"
	UploadsRoute           = "/uploads"
	adminRoute             = "/admin"
	adminFileUploadRoute   = adminRoute + "/upload"
	adminLoginRoute        = adminRoute + "/login"
	adminConfirmLoginRoute = adminRoute + "/confirm-login"
)

func newServer() *server {
	s := &server{}

	// Load, decode and validate config
	config.MustLoad(&s.conf,
		config.TryLoadFile("config.dev.json", config.JSONDecoder), // load config.dev.json file
		config.TryLoadFile("config.json", config.JSONDecoder),     // or load config.json file
	)
	err := s.conf.validate()
	if err != nil {
		panic(err)
	}

	// Init logger
	if s.conf.Env == "prod" {
		panic("todo: set logger for prod")
	} else {
		s.logger = logs.NewTextLogger(os.Stderr, logs.MustOpenLogFile(s.conf.LogFilePath))
	}
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile) // also init std/logger

	// Init emailer
	if s.conf.Env == "prod" {
		s.emailer = email.NewSMTPEmailer(s.conf.SMTPEmailerConfig)
	} else {
		s.emailer = email.NewMockEmailer(os.Stderr, nil)
	}

	// Init authenticator
	if s.conf.Env == "prod" {
		panic("todo: set authenticator for prod")
	} else {
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

	// Init contact form DB
	if s.conf.Env == "prod" {
		panic("todo: set prod DB for contact form submissions")
	} else {
		s.contactForms = contact.MockDB{}
	}

	// Init media storage
	if s.conf.Env == "prod" {
		panic("todo: set prod storage for file uploads")
	} else {
		s.uploads = media.NewLocalDiskStorage("uploads")
	}

	// Init HTTP handler
	s.initHTTPHandler()

	return s
}

func (s *server) initHTTPHandler() {
	// Init HTTP endpoint h
	h := web.Routes{}
	h.Handle(serveHomePage(s), web.MatchPath("/"), web.MatchMethodGET)
	h.Handle(serveContactForm(s), web.MatchPath("/contact"), web.MatchMethodPOST)
	h.Handle(web.FileServer("uploads", UploadsRoute+"/"), web.MatchPathPrefix(UploadsRoute+"/"))

	h.Handle(authMiddleware(s)(serveAdminPage(s)), web.MatchPath(adminRoute), web.MatchMethodGET)
	h.Handle(authMiddleware(s)(serveAdminFileUpload(s)), web.MatchPath(adminFileUploadRoute), web.MatchMethodPOST)
	h.Handle(serveLoginForm(s), web.MatchPath(adminLoginRoute), web.MatchMethodPOST)
	h.Handle(serveConfirmLoginForm(s), web.MatchPath(adminConfirmLoginRoute), web.MatchMethodGET)

	h.Handle(web.ServeMonochromeFaviconPNG(nil), web.MatchPath("/favicon.ico"), web.MatchMethodGET)
	h.Handle(web.ServeSitemapXML("example.com", "/"), web.MatchPath("/sitemap.xml"), web.MatchMethodGET)
	h.Handle(serve404Page(s), web.CatchAll)
	s.h = h

	// Wrap global middleware
	s.h = web.AccessLoggingMiddleware(s.logger)(s.h)
	s.h = web.PanicRecoveryMiddleware(s.onPanic)(s.h)
}

// Log startup info
// and listen for incoming connections (with graceful shutdown)
func (s *server) run() error {
	s.logger.Log("starting HTTP server on port " + strconv.Itoa(s.conf.Port))
	return web.RunServer(web.NewServerWithDefaults(s.h, s.conf.Port))
}
