package htmlg

type HTMLStringer interface {
	HTMLString() string
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

type Page struct{ Root *ElementNode }

func NewPage(attrs Attrs, children ...HTMLStringer) *Page {
	return &Page{Root: TagHTML.With(attrs, children...)}
}

func (p *Page) HTMLString() string { return "<!DOCTYPE html>\n" + p.Root.HTMLString() }

type PageHead struct {
	Title       string
	FaviconURL  string
	Description string
}

func NewPageHead(ph PageHead) *ElementNode {
	children := []HTMLStringer{}
	children = append(children, TagTitle.Text(ph.Title))
	children = append(children, TagLink.With(Attrs{"rel": "icon", "href": ph.FaviconURL}))
	children = append(children, TagMeta.With(Attrs{"description": ph.Description}))
	return TagHead.With(nil, children...)
}

type Fragment []HTMLStringer

func NewFragment(els ...HTMLStringer) Fragment { return Fragment(els) }

func (frag Fragment) HTMLString() string {
	out := ""
	for _, el := range frag {
		out += el.HTMLString()
	}
	return out
}
