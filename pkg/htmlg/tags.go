package htmlg

type Tag string

func (t Tag) Element() *ElementNode { return &ElementNode{Tag: t} }

func (t Tag) With(attrs Attrs, children ...HTMLStringer) *ElementNode {
	return &ElementNode{Tag: t, Attrs: attrs, Children: children}
}

// Creates a new element with this tag and set the given text as a child text node.
func (t Tag) Text(s string) *ElementNode { return t.With(nil, TextNode(s)) }

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
