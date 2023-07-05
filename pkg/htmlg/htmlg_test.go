package htmlg

import (
	"testing"

	"github.com/ejuju/go-utils/pkg/cssg"
)

func TestHTMLGeneration(t *testing.T) {
	t.Run("can generate valid HTML strings", func(t *testing.T) {
		tests := []struct {
			description    string
			input          Stringer
			expectedOutput string
		}{
			{
				description:    "can generate a valid h1 tag (no attributes)",
				input:          Create("h1", WithString("Hello world")),
				expectedOutput: "<h1>Hello world</h1>",
			},
			{
				description:    "can generate a valid h1 tag (with attributes)",
				input:          Create("h1", SetAttr([2]string{"id", "hero"}), WithString("Hello world")),
				expectedOutput: `<h1 id="hero">Hello world</h1>`,
			},
			{
				description:    "can generate a valid void tag (no closing tag and inner HTML)",
				input:          Create("br"),
				expectedOutput: "<br>",
			},
			{
				description:    "can generate a valid element with one child",
				input:          Create("h1", Wrap(Create("span", WithString("hello")))),
				expectedOutput: "<h1><span>hello</span></h1>",
			},
			{
				description:    "can generate a valid element with many children",
				input:          Create("h1", Wrap(Create("span"), Create("div"))),
				expectedOutput: "<h1><span></span><div></div></h1>",
			},
			{
				description:    "can generate a valid element with many children (incl. nil children)",
				input:          Create("h1", Wrap(nil, Create("span"), nil, Create("div"), nil)),
				expectedOutput: "<h1><span></span><div></div></h1>",
			},
			{
				description:    "can generate a valid element with inline styles",
				input:          Create("h1", Style(cssg.DeclarationGroup{{"color", "black"}, {"margin", "16px"}}...)),
				expectedOutput: `<h1 style="color: black; margin: 16px"></h1>`,
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
