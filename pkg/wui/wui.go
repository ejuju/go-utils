package wui

import (
	"strings"

	"github.com/ejuju/go-utils/pkg/cssg"
	"github.com/ejuju/go-utils/pkg/htmlg"
)

// Page implements HTMLStringer and prepends the HTML doctype declaration
// to the underlying content.
// Expects HTML tag as direct child.
type Page struct{ root *htmlg.Element }

func NewPage(attrs [][2]string, head, body *htmlg.Element) *Page {
	return &Page{root: HTML(htmlg.SetAttrs(attrs...), htmlg.Wrap(head, body))}
}

// HTMLString returns the underlying content with the HTML doctype declaration prefixed.
func (p *Page) HTMLString() string { return "<!DOCTYPE html>" + p.root.HTMLString() }

// NewFavicon creates a new <link> element with the rel icon attribute and href set to url.
func NewFavicon(url string) *htmlg.Element {
	return Link(htmlg.SetAttrs(AttrRel("icon"), AttrHref(url)))
}

func NewMeta(name, content string) *htmlg.Element {
	return Meta(htmlg.SetAttrs(AttrName(name), AttrContent(content)))
}

func NewTitle(s string) *htmlg.Element              { return Title(htmlg.Wrap(htmlg.String(s))) }
func NewMetaDescription(desc string) *htmlg.Element { return NewMeta("description", desc) }
func NewMetaColorScheme(s string) *htmlg.Element    { return NewMeta("color-scheme", s) }

// NewGlobalStyle creates a new <style> element with the provided rules.
func NewGlobalStyle(rules cssg.RuleGroup) *htmlg.Element {
	return Style(htmlg.Wrap(htmlg.String(" " + rules.CSSString())))
}

func NewHyperlink(anchor, href string) *htmlg.Element {
	return A(htmlg.SetAttr(AttrHref(href)), htmlg.Wrap(htmlg.String(anchor)))
}

func inputWithOptions(typ, name, value string) *htmlg.Element {
	return Input(htmlg.SetAttrs(AttrType(typ), AttrName(name), AttrValue(value)))
}

func NewInputEmail(name, placeholder, value string) *htmlg.Element {
	return inputWithOptions("email", name, value).Apply(htmlg.SetAttr(AttrPlaceholder(placeholder)))
}

func NewInputMultilineText(name string, placeholder string, value ...string) *htmlg.Element {
	return Textarea(
		htmlg.Wrap(htmlg.String(strings.Join(value, "\n"))),
		htmlg.SetAttrs(AttrName(name), AttrPlaceholder(placeholder)),
	)
}

func NewInputSubmit(txt string) *htmlg.Element {
	return Input(htmlg.SetAttrs(AttrType("submit"), AttrValue(txt)))
}

func NewForm(action string, method string) *htmlg.Element {
	return Form(htmlg.SetAttrs(AttrAction(action), AttrMethod(method)))
}

func NewButton(text string, onclick string) *htmlg.Element {
	return Button(htmlg.Wrap(htmlg.String(text)), htmlg.SetAttr(AttrOnclick(onclick)))
}

// Wrap returns a HTML fragment containing the given tag with one children in each.
func Wrap(children []htmlg.Stringer, tag string, modifiers ...htmlg.Modifier) htmlg.Fragment {
	items := make(htmlg.Fragment, len(children))
	for i, child := range children {
		e := htmlg.Create(tag, htmlg.Wrap(child))
		for _, modify := range modifiers {
			modify(e)
		}
		items[i] = e
	}
	return items
}
