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

func NewElementNode(t Tag) *ElementNode { return &ElementNode{Tag: t, Attrs: make(Attrs)} }

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
func (n *ElementNode) Wrap(children ...HTMLStringer) *ElementNode {
	n.Children = children
	return n
}

func (n *ElementNode) WrapText(s string) *ElementNode {
	n.Children = []HTMLStringer{TextNode(s)}
	return n
}

func (n *ElementNode) AppendChildren(children ...HTMLStringer) *ElementNode {
	n.Children = append(n.Children, children...)
	return n
}

type Page struct{ Root *ElementNode }

func NewPage(attrs Attrs, children ...HTMLStringer) *Page {
	return &Page{Root: TagHTML.El().WithAttrs(attrs).Wrap(children...)}
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

type Tag string

func (t Tag) El() *ElementNode { return NewElementNode(t) }

const (
	TagA          Tag = "a"
	TagAddress    Tag = "address"
	TagAbbr       Tag = "abbr"
	TagArticle    Tag = "article"
	TagAside      Tag = "aside"
	TagArea       Tag = "area"
	TagAudio      Tag = "audio"
	TagB          Tag = "base"
	TagBase       Tag = "base"
	TagBdi        Tag = "bdi"
	TagBdo        Tag = "bdo"
	TagBlockquote Tag = "blockquote"
	TagBody       Tag = "body"
	TagBr         Tag = "br"
	TagButton     Tag = "button"
	TagCite       Tag = "cite"
	TagCol        Tag = "col"
	TagCode       Tag = "code"
	TagData       Tag = "data"
	TagDd         Tag = "dd"
	TagDiv        Tag = "div"
	TagDl         Tag = "dl"
	TagDfn        Tag = "dfn"
	TagDt         Tag = "dt"
	TagEm         Tag = "em"
	TagEmbed      Tag = "embed"
	TagFooter     Tag = "footer"
	TagForm       Tag = "form"
	TagFigcaption Tag = "figcaption"
	TagFigure     Tag = "figure"
	TagH1         Tag = "h1"
	TagH2         Tag = "h2"
	TagH3         Tag = "h3"
	TagH4         Tag = "h4"
	TagH5         Tag = "h5"
	TagH6         Tag = "h6"
	TagHgroup     Tag = "hgroup"
	TagHTML       Tag = "html"
	TagHead       Tag = "head"
	TagHeader     Tag = "header"
	TagHr         Tag = "hr"
	TagI          Tag = "i"
	TagIframe     Tag = "iframe"
	TagImg        Tag = "img"
	TagInput      Tag = "input"
	TagKbd        Tag = "kbd"
	TagLabel      Tag = "label"
	TagLegend     Tag = "legend"
	TagLink       Tag = "link"
	TagLi         Tag = "li"
	TagMain       Tag = "main"
	TagMap        Tag = "map"
	TagMath       Tag = "math"
	TagMark       Tag = "mark"
	TagMeta       Tag = "meta"
	TagMenu       Tag = "menu"
	TagNav        Tag = "nav"
	TagOl         Tag = "ol"
	TagObject     Tag = "object"
	TagP          Tag = "p"
	TagPicture    Tag = "picture"
	TagPortal     Tag = "portal"
	TagPre        Tag = "pre"
	TagQ          Tag = "q"
	TagU          Tag = "u"
	TagUl         Tag = "ul"
	TagS          Tag = "s"
	TagSamp       Tag = "samp"
	TagSource     Tag = "source"
	TagSection    Tag = "section"
	TagSmall      Tag = "small"
	TagSpan       Tag = "span"
	TagStyle      Tag = "style"
	TagStrong     Tag = "strong"
	TagSub        Tag = "sub"
	TagSup        Tag = "sup"
	TagSvg        Tag = "svg"
	TagTextarea   Tag = "textarea"
	TagTime       Tag = "time"
	TagTitle      Tag = "title"
	TagTrack      Tag = "track"
	TagVar        Tag = "var"
	TagVideo      Tag = "video"
	TagWbr        Tag = "wbr"
)

// Self-closing tags
var VoidElementsTags = map[Tag]struct{}{
	TagArea:   {},
	TagBase:   {},
	TagBr:     {},
	TagCol:    {},
	TagEmbed:  {},
	TagHr:     {},
	TagImg:    {},
	TagInput:  {},
	TagLink:   {},
	TagMeta:   {},
	TagSource: {},
	TagTrack:  {},
	TagWbr:    {},
}
