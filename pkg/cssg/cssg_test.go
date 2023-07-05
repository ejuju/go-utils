package cssg

import "testing"

func TestCSSGeneration(t *testing.T) {
	t.Run("can generate valid CSS strings", func(t *testing.T) {
		tests := []struct {
			description    string
			input          Stringer
			expectedOutput string
		}{
			{
				description:    "can generate valid rule",
				input:          NewRule("body").With(DeclarationGroup{{"margin", "0"}, {"color", "blue"}}),
				expectedOutput: "body { margin: 0; color: blue }",
			},
			{
				description: "can generate valid rule group",
				input: RuleGroup{
					NewRule(":root").With(DeclarationGroup{{"--myvar", "myvalue"}}),
					NewRule("*").With(DeclarationGroup{{"font", "inherit"}}),
				},
				expectedOutput: ":root { --myvar: myvalue }  * { font: inherit }",
			},
		}

		for _, test := range tests {
			t.Run(test.description, func(t *testing.T) {
				got := test.input.CSSString()
				want := test.expectedOutput
				if got != want {
					t.Fatalf("want %q but got %q", want, got)
				}
			})
		}
	})
}
