package main

import (
	"encoding/json"
	"image/color"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/ejuju/go-utils/pkg/auth"
	"github.com/ejuju/go-utils/pkg/conf"
	"github.com/ejuju/go-utils/pkg/contact"
	"github.com/ejuju/go-utils/pkg/email"
	"github.com/ejuju/go-utils/pkg/logs"
	"github.com/ejuju/go-utils/pkg/media"
	"github.com/ejuju/go-utils/pkg/web"
)

func main() {
	err := newServer().run()
	if err != nil {
		panic(err)
	}
}

type config struct {
	Env               string                   `json:"env"`                 // ex: "prod" or "dev"
	Host              string                   `json:"host"`                // ex: "example.com"
	Port              int                      `json:"port"`                // ex: 8080
	Version           string                   `json:"version"`             // ex: "1.0.2"
	AdminEmailAddr    string                   `json:"admin_email_addr"`    // ex: "admin@example.com"
	SMTPEmailerConfig *email.SMTPEmailerConfig `json:"smtp_emailer_config"` // see pkg/email
}

var brandColor = color.RGBA{R: 0, G: 250, B: 255, A: 255}

type server struct {
	h             http.Handler
	conf          config
	logger        logs.Logger
	emailer       email.Emailer
	contactForms  contact.Forms
	uploads       media.Storage
	authenticator *auth.OTPAuthenticator
}

const (
	contactRoute           = "/contact"
	adminRoute             = "/admin"
	adminFileUploadRoute   = adminRoute + "/upload"
	adminLoginRoute        = adminRoute + "/login"
	adminConfirmLoginRoute = adminRoute + "/confirm-login"
)

func newServer() *server {
	s := &server{}

	// Load and decode config file
	conf.MustLoad(&s.conf, conf.LoadFile("config.json", json.Unmarshal))

	// Init logger
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	if s.conf.Env == "prod" {
		s.logger = logs.NewTextLogger(os.Stderr)
	} else {
		s.logger = logs.NewTextLogger(os.Stderr, logs.MustOpenLogFile("/tmp/app_logs.txt"))
	}

	// Init emailer
	if s.conf.Env == "prod" {
		s.emailer = email.NewSMTPEmailer(s.conf.SMTPEmailerConfig)
	} else {
		s.emailer = email.NewMockEmailer(os.Stderr, nil)
	}

	// Init authenticator
	if s.conf.Env == "prod" {
		panic("todo")
	} else {
		s.authenticator = auth.NewOTPAuthenticator(&auth.OTPAuthenticatorConfig{
			Host:                s.conf.Host,
			Users:               auth.NewStaticUsers(s.conf.AdminEmailAddr),
			Emailer:             s.emailer,
			ConfirmLoginRoute:   adminConfirmLoginRoute,
			SuccessfulLoginPath: adminRoute,
			CookieName:          "auth",
		})
	}

	// Init contact form DB
	if s.conf.Env == "prod" {
		panic("todo")
	} else {
		s.contactForms = contact.MockDB{}
	}

	// Init media storage
	if s.conf.Env == "prod" {
		panic("todo")
	} else {
		s.uploads = media.NewLocalDiskStorage("uploads")
	}

	// Init HTTP endpoint h
	h := web.Routes{}
	h.With(handleHomePage(s), web.MatchPath("/"), web.MatchMethodGET)
	h.With(handleContactForm(s), web.MatchPath("/contact"), web.MatchMethodPOST)

	h.With(authMiddleware(s)(handleAdminPage(s)), web.MatchPath(adminRoute), web.MatchMethodGET)
	h.With(authMiddleware(s)(handleAdminFileUpload(s)), web.MatchPath(adminFileUploadRoute), web.MatchMethodPOST)
	h.With(handleLoginForm(s), web.MatchPath(adminLoginRoute), web.MatchMethodPOST)
	h.With(handleConfirmLoginForm(s), web.MatchPath(adminConfirmLoginRoute), web.MatchMethodGET)
	h.With(web.ServeSimpleFavicon(brandColor), web.MatchPath("/favicon.ico"), web.MatchMethodGET)
	h.With(web.ServeSitemapXML("example.com", "/"), web.MatchPath("/sitemap.xml"), web.MatchMethodGET)
	h.With(web.FileServer("uploads", adminFileUploadRoute+"/"), web.MatchPathPrefix(adminFileUploadRoute+"/"))
	h.With(handleNotFound(s), web.CatchAll)
	s.h = h

	// Wrap global middleware
	s.h = web.AccessLoggingMiddleware(s.logger)(s.h)
	s.h = web.PanicRecoveryMiddleware(s.onPanic)(s.h)

	return s
}

// Log startup info
// and listen for incoming connections (with graceful shutdown)
func (s *server) run() error {
	s.logger.Log("starting HTTP server on port " + strconv.Itoa(s.conf.Port))
	return web.RunServer(web.NewServerWithDefaults(s.h, s.conf.Port))
}
