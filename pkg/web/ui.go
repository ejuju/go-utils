package web

import (
	"html/template"
	"net/http"
)

func MustParseHTMLTemplate(commonTmpls []string, path string) *template.Template {
	return template.Must(template.ParseFiles(append(commonTmpls, path)...))
}

// Will panic when an error occurs during rendering, make sure you handle panic recovery in a middleware.
func RenderHTMLTmpl(
	w http.ResponseWriter,
	statusCode int,
	t *template.Template,
	tname string,
	data map[string]any,
) {
	w.WriteHeader(statusCode)
	err := t.ExecuteTemplate(w, tname, data)
	if err != nil {
		panic(err)
	}
}
