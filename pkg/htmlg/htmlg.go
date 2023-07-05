package htmlg

import (
	"io"
	"net/http"
	"strings"

	"github.com/ejuju/go-utils/pkg/cssg"
)

// Stringer can generate a valid HTML string.
type Stringer interface {
	HTMLString() string
}

// String represents a literal HTML string.
type String string

// HTMLString returns the underlying string without modification.
func (s String) HTMLString() string { return string(s) }

// Text returns an Stringer that returns given strings with <br> tags in between them.
func Text(s ...string) String { return String(strings.Join(s, "<br>")) }

// Element represents a HTML element with a tag and optional attribute and inner HTML.
// Use Create to instanciate an element.
type Element struct {
	tag        string
	attributes [][2]string
	innerHTML  Stringer
}

// HTMLString returns a valid HTML representation of this element.
func (e *Element) HTMLString() string {
	s := "<" + e.tag // Add opening tag start and name
	for _, attr := range e.attributes {
		// Add attributes
		s += " " + attr[0] + `="` + attr[1] + `"`
	}
	s += ">" // Add opening tag end
	if _, ok := voidElementsTags[e.tag]; ok {
		// Return if void element (= element with no closing tag and inner HTML)
		return s
	}
	if e.innerHTML != nil {
		// Add inner HTML if provided
		s += e.innerHTML.HTMLString()
	}
	return s + "</" + e.tag + ">" // Add closing tag
}

// List of void elements for HTML generation. This allows us to know when not to close tags.
var voidElementsTags = map[string]struct{}{
	"area":   {},
	"base":   {},
	"br":     {},
	"col":    {},
	"embed":  {},
	"hr":     {},
	"img":    {},
	"input":  {},
	"link":   {},
	"meta":   {},
	"source": {},
	"track":  {},
	"wbr":    {},
}

// Create instanciates a new Element with the given tag.
func Create(t string, opts ...Modifier) *Element {
	e := &Element{tag: t}
	return e.Apply(opts...)
}

func (e *Element) Tag() string { return e.tag }

// Apply mutates this element using the given modifier(s).
func (e *Element) Apply(opts ...Modifier) *Element {
	for _, opt := range opts {
		opt(e)
	}
	return e
}

// Modifier performs a mutation on an element.
// Eg. Changing its internal state (for ex: adding an attribute)
type Modifier func(*Element)

func SetTag(t string) Modifier { return func(e *Element) { e.tag = t } }

func SetAttr(newAttr [2]string) Modifier {
	return func(e *Element) {
		for currI, currAttr := range e.attributes {
			// Replace existing attribute if exists
			if newAttr[0] == currAttr[0] {
				e.attributes[currI] = newAttr
				return
			}
		}
		// Append new attribute if does not exist yet
		e.attributes = append(e.attributes, newAttr)
	}
}

func SetAttrs(newAttrs ...[2]string) Modifier {
	return func(e *Element) {
		for _, newAttr := range newAttrs {
			e.Apply(SetAttr(newAttr))
		}
	}
}

func Wrap(s ...Stringer) Modifier {
	return func(e *Element) {
		if len(s) == 1 {
			e.innerHTML = s[0]
			return
		}
		e.innerHTML = NewFragment(s...)
	}
}

func WithString(s string) Modifier      { return Wrap(String(s)) }
func WithText(lines ...string) Modifier { return Wrap(Text(lines...)) }

func AddChild(s ...Stringer) Modifier {
	return func(e *Element) {
		if len(s) == 0 {
			return
		}
		e.innerHTML = NewFragment(append([]Stringer{e.innerHTML}, s...)...)
	}
}

// Sets the inline style of the element to the given CSS declarations.
func Style(d ...cssg.Decl) Modifier {
	return func(e *Element) { e.Apply(SetAttr([2]string{"style", cssg.DeclarationGroup(d).CSSString()})) }
}

// Fragment represents a group of HTMLStringers. This can be several HTML sibling elements.
type Fragment []Stringer

// NewFragment returns a fragment made of the given HTMLStringers.
func NewFragment(siblings ...Stringer) Fragment { return siblings }

// HTMLString returns the concatenated underyling HTMLStringers.
func (f Fragment) HTMLString() string {
	out := ""
	for _, item := range f {
		if item == nil {
			continue
		}
		out += item.HTMLString()
	}
	return out
}

// Send writes a HTML string to the given response writer.
// It sets the content-type header to "text/html; charset=utf-8" and writes the given status code.
func Send(w http.ResponseWriter, statusCode int, s string) (int, error) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(statusCode)
	return io.WriteString(w, s)
}

// SendStringer renders and writes a Stringer to the given response writer.
// It sets the content-type header to "text/html; charset=utf-8" and writes the given status code.
func SendStringer(w http.ResponseWriter, statusCode int, s Stringer) (int, error) {
	return Send(w, statusCode, s.HTMLString())
}
