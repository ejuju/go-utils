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

var commonTmpls = []string{"_layout.gohtml"}

const layoutTmpl = "layout"

var (
	errorPageTmpl = web.MustParseHTMLTemplate(commonTmpls, "_error.gohtml")
	homePageTmpl  = web.MustParseHTMLTemplate(commonTmpls, "_home.gohtml")
	adminPageTmpl = web.MustParseHTMLTemplate(commonTmpls, "_admin.gohtml")
	loginPageTmpl = web.MustParseHTMLTemplate(commonTmpls, "_login.gohtml")
)

func (s *server) onPanic(w http.ResponseWriter, r *http.Request, err any) {
	// TODO: notify admin

	// Log error
	stackstr := strings.ReplaceAll(string(debug.Stack()), "\n", " ")
	stackstr = strings.ReplaceAll(string(stackstr), "\t", " ")
	s.logger.Log(fmt.Sprintf("%s %v %s", web.VisitorHash(r), err, stackstr))

	// Respond to client
	renderErrorPage(w, r, http.StatusInternalServerError, fmt.Errorf("%v", err))
}

func renderPage(w http.ResponseWriter, r *http.Request, statusCode int, t *template.Template, data map[string]any) {
	web.RenderHTMLTmpl(w, statusCode, t, layoutTmpl, map[string]any{
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

func authMiddleware(s *server) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session := s.authenticator.Authenticate(w, r)
			if session == nil {
				renderPage(w, r, http.StatusForbidden, loginPageTmpl, nil)
				return
			}
			h.ServeHTTP(w, r)
		})
	}
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

func serveAdminPage(s *server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fnames, err := s.uploads.List()
		if err != nil {
			renderPage(w, r, http.StatusInternalServerError, adminPageTmpl, map[string]any{"ListError": err})
			return
		}
		renderPage(w, r, http.StatusOK, adminPageTmpl, map[string]any{"Uploads": fnames})
	}
}

func serveAdminFileUpload(s *server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(8000); err != nil {
			renderPage(w, r, http.StatusBadRequest, adminPageTmpl, map[string]any{"UploadError": err})
			return
		}
		f, fh, err := r.FormFile("file")
		if err != nil {
			renderPage(w, r, http.StatusBadRequest, adminPageTmpl, map[string]any{"UploadError": err})
			return
		}
		if err = s.uploads.Store(fh.Filename, f); err != nil {
			renderPage(w, r, http.StatusInternalServerError, adminPageTmpl, map[string]any{"UploadError": err})
			return
		}
		http.Redirect(w, r, adminRoute, http.StatusSeeOther)
	}
}

func serveLoginForm(s *server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		emailAddr := r.FormValue("email-address")
		err := s.authenticator.SendLoginLinkByEmail(emailAddr)
		if err != nil {
			renderErrorPage(w, r, http.StatusBadRequest, err)
			return
		}
		renderPage(w, r, http.StatusOK, loginPageTmpl, map[string]any{"FormSuccess": true})
	}
}

func serveConfirmLoginForm(s *server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.authenticator.LoginWithLink(w, r,
			func(err error) { renderErrorPage(w, r, http.StatusInternalServerError, err) },
			func(err error) { renderErrorPage(w, r, http.StatusBadRequest, err) },
		)
	}
}
