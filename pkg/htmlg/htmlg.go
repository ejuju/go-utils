package htmlg

import "net/http"

type HTMLStringer interface {
	HTMLString() string
}

func Render(hs HTMLStringer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(hs.HTMLString()))
	}
}

type TextNode string

func (n TextNode) HTMLString() string { return string(n) }

// Element attribute
type Attrs map[string]string

type ElementNode struct {
	Tag      Tag
	Attrs    Attrs
	Children []HTMLStringer
}

func (n *ElementNode) HTMLString() string {
	// Add opening tag and attributes if any
	out := "<" + string(n.Tag)
	for k, v := range n.Attrs {
		out += " " + k + "=\"" + v + "\""
	}
	out += ">"

	// Stop execution if tag is self-closing (void element)
	if _, isVoidEl := VoidElementsTags[Tag(n.Tag)]; isVoidEl {
		return out
	}

	// Add children
	for _, child := range n.Children {
		out += child.HTMLString()
	}

	// Add closing tag
	return out + "</" + string(n.Tag) + ">"
}

func (n *ElementNode) WithAttr(k, v string) *ElementNode  { n.Attrs[k] = v; return n }
func (n *ElementNode) WithAttrs(attrs Attrs) *ElementNode { n.Attrs = attrs; return n }
func (n *ElementNode) WithChildren(children ...HTMLStringer) *ElementNode {
	n.Children = children
	return n
}

func (n *ElementNode) AppendChildren(children ...HTMLStringer) *ElementNode {
	n.Children = append(n.Children, children...)
	return n
}

func (n *ElementNode) WrapText(s string) *ElementNode {
	n.Children = []HTMLStringer{TextNode(s)}
	return n
}

type Page struct{ Root *ElementNode }

func NewPage(attrs Attrs, children ...HTMLStringer) *Page {
	return &Page{Root: TagHTML.With(attrs, children...)}
}

func (p *Page) HTMLString() string { return "<!DOCTYPE html>\n" + p.Root.HTMLString() }

type Fragment []HTMLStringer

func NewFragment(els ...HTMLStringer) Fragment { return Fragment(els) }

func (frag Fragment) HTMLString() string {
	out := ""
	for _, el := range frag {
		out += el.HTMLString()
	}
	return out
}
