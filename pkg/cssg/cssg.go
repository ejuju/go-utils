package cssg

import "strings"

// Stringer can generate a valid CSS string.
type Stringer interface{ CSSString() string }

// Decl represents a CSS declaration containing a property and value.
type Decl [2]string

// CSSString generates a valid CSS string representation of the declaration.
func (d Decl) CSSString() string { return d[0] + ": " + d[1] }

// DeclarationGroup groups multiple declarations together.
type DeclarationGroup []Decl

func (dg DeclarationGroup) CSSString() string {
	out := ""
	for i, d := range dg {
		if i != 0 {
			out += "; "
		}
		out += d.CSSString()
	}
	return out
}

// Rule represents a CSS rule (linking declarations to selectors).
type Rule struct {
	selectors    []string
	declarations DeclarationGroup
}

// CSSString generates a valid CSS string representation of the rule.
func (r *Rule) CSSString() string {
	return strings.Join(r.selectors, ", ") + " { " + r.declarations.CSSString() + " }"
}

// NewRule allocates a new CSS rule with the given selectors.
func NewRule(selectors ...string) *Rule { return (&Rule{}).Select(selectors...) }

// Select appends the given selectors to the rules' list of selectors.
func (r *Rule) Select(selectors ...string) *Rule {
	r.selectors = append(r.selectors, selectors...)
	return r
}

// With appends the given declarations to the rules' list of declarations.
func (r *Rule) With(declarations DeclarationGroup) *Rule {
	r.declarations = append(r.declarations, declarations...)
	return r
}

// RuleGroup represent a list of CSS rules.
type RuleGroup []*Rule

// CSSString generates a valid CSS string representation of the rules.
func (rs RuleGroup) CSSString() string {
	out := ""
	for i, r := range rs {
		if i != 0 {
			out += "  "
		}
		out += r.CSSString()
	}
	return out
}
