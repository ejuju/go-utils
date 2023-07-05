package wui

import (
	"net/http"
	"testing"

	"github.com/ejuju/go-utils/pkg/htmlg"
)

func TestHTMLGeneration(t *testing.T) {
	t.Run("can generate valid HTML elements", func(t *testing.T) {
		tests := []struct {
			description    string
			input          htmlg.Stringer
			expectedOutput string
		}{
			{
				description:    "can generate a HTML page with DOCTYPE declaration",
				input:          NewPage(nil, Head(), Body()),
				expectedOutput: `<!DOCTYPE html><html><head></head><body></body></html>`,
			},
			{
				description:    "can generate a valid favicon link",
				input:          NewFavicon("/favicon.ico"),
				expectedOutput: `<link rel="icon" href="/favicon.ico">`,
			},
			{
				description:    "can generate a valid title",
				input:          NewTitle("MyTitle"),
				expectedOutput: `<title>MyTitle</title>`,
			},
			{
				description:    "can generate a valid meta-description",
				input:          NewMetaDescription("MyDescription"),
				expectedOutput: `<meta name="description" content="MyDescription">`,
			},
			{
				description:    "can generate a valid meta-description",
				input:          NewForm("/endpoint", http.MethodPost),
				expectedOutput: `<form action="/endpoint" method="POST"></form>`,
			},
		}

		for _, test := range tests {
			t.Run(test.description, func(t *testing.T) {
				got := test.input.HTMLString()
				want := test.expectedOutput
				if got != want {
					t.Fatalf("want %q but got %q", want, got)
				}
			})
		}
	})
}
