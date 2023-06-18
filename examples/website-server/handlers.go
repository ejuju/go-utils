package main

import (
	"fmt"
	"html/template"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/ejuju/go-utils/pkg/contact"
	"github.com/ejuju/go-utils/pkg/web"
)

const (
	contactRoute           = "/contact"
	PublicFilesRoute       = "/uploads"
	adminRoute             = "/admin"
	adminFileUploadRoute   = adminRoute + "/upload"
	adminLoginRoute        = adminRoute + "/login"
	adminConfirmLoginRoute = adminRoute + "/confirm-login"
)

var commonTmpls = []string{"_layout.gohtml"}

const layoutTmpl = "layout"

var (
	errorPageTmpl = web.MustParseHTMLTemplate(commonTmpls, "_error.gohtml")
	homePageTmpl  = web.MustParseHTMLTemplate(commonTmpls, "_home.gohtml")
)

func (s *server) onPanic(err any, w http.ResponseWriter, r *http.Request) {
	// TODO: notify admin

	// Log error
	stackstr := strings.ReplaceAll(string(debug.Stack()), "\n", " ")
	stackstr = strings.ReplaceAll(string(stackstr), "\t", " ")
	s.logger.Log(fmt.Sprintf("%s %v %s", web.VisitorHash(r, s.conf.UsesXForwardedFor), err, stackstr))

	// Respond to client
	renderErrorPage(w, r, http.StatusInternalServerError, fmt.Errorf("%v", err))
}

func renderPage(w http.ResponseWriter, r *http.Request, statusCode int, t *template.Template, data map[string]any) {
	web.RenderHTMLTemplate(w, statusCode, t, layoutTmpl, map[string]any{
		"Request": r,
		"Local":   data,
	})
}

func renderErrorPage(w http.ResponseWriter, r *http.Request, statusCode int, err error) {
	renderPage(w, r, statusCode, errorPageTmpl, map[string]any{
		"Error":      err.Error(),
		"StatusCode": statusCode,
		"StatusText": http.StatusText(statusCode),
	})
}

func serve404Page(s *server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		renderErrorPage(w, r, http.StatusNotFound, fmt.Errorf("%s not found", r.URL))
	}
}

func serveHomePage(s *server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) { renderPage(w, r, http.StatusOK, homePageTmpl, nil) }
}

func serveContactForm(s *server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse form data
		form, err := contact.ParseAndValidateForm(r, "email-address", "message")
		if err != nil {
			renderPage(w, r, http.StatusBadRequest, homePageTmpl, map[string]any{"FormError": err})
			return
		}

		// TODO: Send notification email to admin
		// TODO: Send confirmation email to user
		fmt.Printf("New contact form submission received:\n%#v\n", form)

		// Store message
		err = s.contactForms.SaveNew(form)
		if err != nil {
			renderPage(w, r, http.StatusInternalServerError, homePageTmpl, map[string]any{"FormError": err})
			return
		}

		renderPage(w, r, http.StatusOK, homePageTmpl, map[string]any{"FormSuccess": true})
	}
}
