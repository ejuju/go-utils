package htmlg

type HTMLStringer interface {
	HTMLString() string
}

type TextNode string

func (n TextNode) HTMLString() string { return string(n) }

type ElementNode struct {
	Tag      Tag
	Attrs    Attributes
	Children []HTMLStringer
}

type Attributes map[string]string

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

func NewPage(quirksMode bool, attrs map[string]string, children ...HTMLStringer) *Page {
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
	children = append(children, TagLink.With(Attributes{"rel": "icon", "href": ph.FaviconURL}))
	children = append(children, TagMeta.With(Attributes{"description": ph.Description}))
	return TagHead.With(nil, children...)
}
