// SPDX-Licence-Identifier: EUPL-1.2

package ctml

import (
	"bytes"
	"encoding/xml"
	"io"
	"strings"

	core "dappco.re/go"
	html "dappco.re/go/html"
)

// parser walks one document with encoding/xml.Decoder. eachStack tracks
// the currently-open <each as="..."> names so a {{path}} token can be
// rejected at its own position the moment its root name is unbound,
// rather than failing later, unhelpfully, at materialise time.
type parser struct {
	dec       *xml.Decoder
	eachStack []string
}

// Parse parses src into the node tree the Go builder API would produce for
// the same page -- see docs/ctml.md for the grammar. bindings is optional;
// an <each> whose items name is absent from Sequences renders as an empty
// list rather than failing to parse.
//
// Usage example: tree, err := ctml.Parse(src)
func Parse(src []byte, bindings ...Bindings) (html.Node, error) {
	root, bnd, err := parseRoot(src, bindings)
	if err != nil {
		return nil, err
	}
	return materialise(root, rootResolver, bnd), nil
}

// ParseLayout is Parse specialised to documents whose root is <layout>,
// returning the concrete *html.Layout so callers can chain further
// builder calls (Responsive.Add, further slot appends) or call
// Layout-specific methods (RenderTerm) directly.
//
// Usage example: layout, err := ctml.ParseLayout(src)
func ParseLayout(src []byte, bindings ...Bindings) (*html.Layout, error) {
	root, bnd, err := parseRoot(src, bindings)
	if err != nil {
		return nil, err
	}
	l, ok := root.(*astLayout)
	if !ok {
		msg := "root element must be <layout>"
		return nil, &ParseError{Line: 1, Col: 1, Msg: msg, Cause: core.E("ctml.ParseLayout", msg, nil)}
	}
	return materialiseLayout(l, rootResolver, bnd), nil
}

// rootResolver backs every {{path}} lookup outside any <each>. It is never
// actually reached by a tree Parse produced: validateEachRef refuses to
// parse a {{path}} token with no enclosing <each as="..."> in the first
// place, so this is defence in depth, not a live path.
func rootResolver(string) (any, bool) { return nil, false }

func parseRoot(src []byte, bindings []Bindings) (astNode, Bindings, error) {
	var bnd Bindings
	if len(bindings) > 0 {
		bnd = bindings[0]
	}

	p := &parser{dec: xml.NewDecoder(bytes.NewReader(src))}

	for {
		tok, err := p.dec.Token()
		if err != nil {
			return nil, bnd, p.wrapXMLErr(err)
		}
		switch t := tok.(type) {
		case xml.StartElement:
			root, err := p.parseElement(t)
			if err != nil {
				return nil, bnd, err
			}
			if err := p.expectTrailingEOF(); err != nil {
				return nil, bnd, err
			}
			return root, bnd, nil
		case xml.CharData:
			if strings.TrimSpace(string(t)) != "" {
				return nil, bnd, p.errAt("unexpected text before root element")
			}
		case xml.Comment, xml.ProcInst, xml.Directive:
			continue
		}
	}
}

// expectTrailingEOF confirms a document has exactly one root element --
// anything but insignificant whitespace/comments after it closes is a
// document defect, not silently-ignored trailing content.
func (p *parser) expectTrailingEOF() error {
	for {
		tok, err := p.dec.Token()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return p.wrapXMLErr(err)
		}
		switch t := tok.(type) {
		case xml.CharData:
			if strings.TrimSpace(string(t)) != "" {
				return p.errAt("unexpected content after root element")
			}
		case xml.Comment, xml.ProcInst, xml.Directive:
			continue
		default:
			return p.errAt("unexpected content after root element")
		}
	}
}

// parseElement dispatches on the fifteen reserved tag names (S:S3); every
// other tag is a literal element (parseEl).
func (p *parser) parseElement(start xml.StartElement) (astNode, error) {
	switch start.Name.Local {
	case "if":
		return p.parseIf(start)
	case "unless":
		return p.parseUnless(start)
	case "switch":
		return p.parseSwitch(start)
	case "case":
		return nil, p.errAt("<case> is only valid as a direct child of <switch>")
	case "entitled":
		return p.parseEntitled(start)
	case "each":
		return p.parseEach(start)
	case "raw":
		return p.parseRaw(start)
	case "layout":
		return p.parseLayout(start)
	case "h", "l", "c", "r", "f":
		return nil, p.errAt("<" + start.Name.Local + "> is only valid as a direct child of <layout>")
	case "responsive":
		return p.parseResponsive(start)
	case "variant":
		return nil, p.errAt("<variant> is only valid as a direct child of <responsive>")
	default:
		return p.parseEl(start)
	}
}

func (p *parser) parseEl(start xml.StartElement) (astNode, error) {
	var attrs []astAttr
	argsAttr, hasArgs := "", false
	for _, a := range start.Attr {
		if a.Name.Local == "args" {
			argsAttr, hasArgs = a.Value, true
			continue
		}
		attrs = append(attrs, astAttr{Key: a.Name.Local, Value: a.Value})
	}

	children, err := p.parseContent(start.Name)
	if err != nil {
		return nil, err
	}
	if hasArgs {
		if err := p.applyArgs(children, argsAttr); err != nil {
			return nil, err
		}
	}
	return &astEl{Tag: start.Name.Local, Attrs: attrs, Children: children}, nil
}

// applyArgs attaches args="..." to the element's sole text/bind child --
// S:S6.4 restricts args to single-run content, so mixed or element content
// is a parse error rather than a silently-dropped attribute.
func (p *parser) applyArgs(children []astNode, raw string) error {
	if len(children) != 1 {
		return p.errAt("args attribute requires exactly one text child")
	}
	tokens, err := p.parseArgTokens(raw)
	if err != nil {
		return err
	}
	switch t := children[0].(type) {
	case *astText:
		t.Args = tokens
	case *astBind:
		t.Args = tokens
	default:
		return p.errAt("args attribute requires text content, not an element child")
	}
	return nil
}

func (p *parser) parseArgTokens(raw string) ([]argToken, error) {
	parts := strings.Split(raw, ",")
	tokens := make([]argToken, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed == "" {
			continue
		}
		if path, ok := matchBindToken(trimmed); ok {
			if err := p.validateEachRef(path); err != nil {
				return nil, err
			}
			tokens = append(tokens, argToken{Path: path, IsPath: true})
			continue
		}
		tokens = append(tokens, argToken{Lit: trimmed})
	}
	return tokens, nil
}

// parseContent reads children until the EndElement matching startName,
// dispatching StartElements recursively and turning each CharData run
// into a text or bind node (S:S2). It is shared by every construct whose
// children are ordinary mixed content; <switch>, <raw>, <layout>, and
// <responsive> have their own loops because each restricts its children
// to one specific reserved element.
func (p *parser) parseContent(startName xml.Name) ([]astNode, error) {
	var nodes []astNode
	for {
		tok, err := p.dec.Token()
		if err != nil {
			return nil, p.wrapXMLErr(err)
		}
		switch t := tok.(type) {
		case xml.StartElement:
			child, err := p.parseElement(t)
			if err != nil {
				return nil, err
			}
			nodes = append(nodes, child)
		case xml.EndElement:
			if t.Name == startName {
				return nodes, nil
			}
			return nil, p.errAt("mismatched closing tag")
		case xml.CharData:
			raw := string(t)
			if strings.TrimSpace(raw) == "" {
				continue // pure structural whitespace between siblings
			}
			node, err := p.charDataNode(normaliseRun(raw))
			if err != nil {
				return nil, err
			}
			nodes = append(nodes, node)
		case xml.Comment, xml.ProcInst, xml.Directive:
			continue
		}
	}
}

func (p *parser) charDataNode(text string) (astNode, error) {
	if path, ok := matchBindToken(text); ok {
		if err := p.validateEachRef(path); err != nil {
			return nil, err
		}
		return &astBind{Path: path}, nil
	}
	return &astText{Key: text}, nil
}

// normaliseRun strips a CharData run's leading/trailing whitespace only
// when that edge whitespace contains a newline -- i.e. only when it is
// source-formatting indentation. A plain inline space at the edge (as in
// "Hello " before a following <strong>) is significant text-flow spacing
// and is left untouched, so mixed content round-trips through .ctml with
// the same spacing an author typed. A run that is whitespace-only never
// reaches this function -- parseContent drops it entirely first.
func normaliseRun(s string) string {
	start := 0
	for start < len(s) && isRunWS(s[start]) {
		start++
	}
	end := len(s)
	for end > start && isRunWS(s[end-1]) {
		end--
	}
	if !strings.Contains(s[:start], "\n") {
		start = 0
	}
	if !strings.Contains(s[end:], "\n") {
		end = len(s)
	}
	return s[start:end]
}

func isRunWS(b byte) bool {
	return b == ' ' || b == '\t' || b == '\r' || b == '\n'
}

// singleChild collapses parsed children to the one Node that If, Unless,
// Entitled, and <case> require, wrapping multiples in astFragment (S:S6.3).
func (p *parser) singleChild(nodes []astNode) (astNode, error) {
	switch len(nodes) {
	case 0:
		return nil, p.errAt("expected at least one child")
	case 1:
		return nodes[0], nil
	default:
		return &astFragment{Children: nodes}, nil
	}
}

func (p *parser) parseIf(start xml.StartElement) (astNode, error) {
	cond, ok := attrValue(start, "cond")
	if !ok {
		return nil, p.errAt("<if> requires a cond attribute")
	}
	children, err := p.parseContent(start.Name)
	if err != nil {
		return nil, err
	}
	child, err := p.singleChild(children)
	if err != nil {
		return nil, err
	}
	return &astIf{CondKey: cond, Child: child}, nil
}

func (p *parser) parseUnless(start xml.StartElement) (astNode, error) {
	cond, ok := attrValue(start, "cond")
	if !ok {
		return nil, p.errAt("<unless> requires a cond attribute")
	}
	children, err := p.parseContent(start.Name)
	if err != nil {
		return nil, err
	}
	child, err := p.singleChild(children)
	if err != nil {
		return nil, err
	}
	return &astUnless{CondKey: cond, Child: child}, nil
}

func (p *parser) parseEntitled(start xml.StartElement) (astNode, error) {
	feature, ok := attrValue(start, "feature")
	if !ok {
		return nil, p.errAt("<entitled> requires a feature attribute")
	}
	children, err := p.parseContent(start.Name)
	if err != nil {
		return nil, err
	}
	child, err := p.singleChild(children)
	if err != nil {
		return nil, err
	}
	return &astEntitled{Feature: feature, Child: child}, nil
}

func (p *parser) parseEach(start xml.StartElement) (astNode, error) {
	items, ok := attrValue(start, "items")
	if !ok {
		return nil, p.errAt("<each> requires an items attribute")
	}
	as, ok := attrValue(start, "as")
	if !ok {
		return nil, p.errAt("<each> requires an as attribute")
	}

	p.eachStack = append(p.eachStack, as)
	children, err := p.parseContent(start.Name)
	p.eachStack = p.eachStack[:len(p.eachStack)-1]
	if err != nil {
		return nil, err
	}
	body, err := p.singleChild(children)
	if err != nil {
		return nil, err
	}
	return &astEach{ItemsName: items, AsName: as, Body: body}, nil
}

func (p *parser) parseSwitch(start xml.StartElement) (astNode, error) {
	on, ok := attrValue(start, "on")
	if !ok {
		return nil, p.errAt("<switch> requires an on attribute")
	}

	cases := make(map[string]astNode)
	for {
		tok, err := p.dec.Token()
		if err != nil {
			return nil, p.wrapXMLErr(err)
		}
		switch t := tok.(type) {
		case xml.StartElement:
			if t.Name.Local != "case" {
				return nil, p.errAt("<switch> only accepts <case> children")
			}
			value, ok := attrValue(t, "value")
			if !ok {
				return nil, p.errAt("<case> requires a value attribute")
			}
			if _, dup := cases[value]; dup {
				return nil, p.errAt("duplicate <case value=\"" + value + "\">")
			}
			children, err := p.parseContent(t.Name)
			if err != nil {
				return nil, err
			}
			child, err := p.singleChild(children)
			if err != nil {
				return nil, err
			}
			cases[value] = child
		case xml.EndElement:
			if t.Name == start.Name {
				return &astSwitch{OnKey: on, Cases: cases}, nil
			}
			return nil, p.errAt("mismatched closing tag in <switch>")
		case xml.CharData:
			if strings.TrimSpace(string(t)) != "" {
				return nil, p.errAt("<switch> only accepts <case> children")
			}
		case xml.Comment, xml.ProcInst, xml.Directive:
			continue
		}
	}
}

func (p *parser) parseRaw(start xml.StartElement) (astNode, error) {
	var b strings.Builder
	for {
		tok, err := p.dec.Token()
		if err != nil {
			return nil, p.wrapXMLErr(err)
		}
		switch t := tok.(type) {
		case xml.CharData:
			b.Write(t)
		case xml.EndElement:
			if t.Name == start.Name {
				return &astRaw{Content: b.String()}, nil
			}
			return nil, p.errAt("mismatched closing tag in <raw>")
		case xml.StartElement:
			return nil, p.errAt("<raw> cannot contain child elements")
		case xml.Comment, xml.ProcInst, xml.Directive:
			continue
		}
	}
}

func (p *parser) parseLayout(start xml.StartElement) (astNode, error) {
	variant, ok := attrValue(start, "variant")
	if !ok {
		return nil, p.errAt("<layout> requires a variant attribute")
	}

	slots := make(map[byte][]astNode)
	for {
		tok, err := p.dec.Token()
		if err != nil {
			return nil, p.wrapXMLErr(err)
		}
		switch t := tok.(type) {
		case xml.StartElement:
			slot := slotLetter(t.Name.Local)
			if slot == 0 {
				return nil, p.errAt("<layout> only accepts <h> <l> <c> <r> <f> children")
			}
			children, err := p.parseContent(t.Name)
			if err != nil {
				return nil, err
			}
			slots[slot] = append(slots[slot], children...)
		case xml.EndElement:
			if t.Name == start.Name {
				return &astLayout{Variant: variant, Slots: slots}, nil
			}
			return nil, p.errAt("mismatched closing tag in <layout>")
		case xml.CharData:
			if strings.TrimSpace(string(t)) != "" {
				return nil, p.errAt("<layout> only accepts <h> <l> <c> <r> <f> children")
			}
		case xml.Comment, xml.ProcInst, xml.Directive:
			continue
		}
	}
}

func (p *parser) parseResponsive(start xml.StartElement) (astNode, error) {
	var variants []astVariant
	for {
		tok, err := p.dec.Token()
		if err != nil {
			return nil, p.wrapXMLErr(err)
		}
		switch t := tok.(type) {
		case xml.StartElement:
			if t.Name.Local != "variant" {
				return nil, p.errAt("<responsive> only accepts <variant> children")
			}
			name, ok := attrValue(t, "name")
			if !ok {
				return nil, p.errAt("<variant> requires a name attribute")
			}
			media, _ := attrValue(t, "media")
			children, err := p.parseContent(t.Name)
			if err != nil {
				return nil, err
			}
			layout, err := p.singleLayoutChild(children)
			if err != nil {
				return nil, err
			}
			variants = append(variants, astVariant{Name: name, Media: media, Layout: layout})
		case xml.EndElement:
			if t.Name == start.Name {
				return &astResponsive{Variants: variants}, nil
			}
			return nil, p.errAt("mismatched closing tag in <responsive>")
		case xml.CharData:
			if strings.TrimSpace(string(t)) != "" {
				return nil, p.errAt("<responsive> only accepts <variant> children")
			}
		case xml.Comment, xml.ProcInst, xml.Directive:
			continue
		}
	}
}

func (p *parser) singleLayoutChild(nodes []astNode) (*astLayout, error) {
	if len(nodes) != 1 {
		return nil, p.errAt("<variant> requires exactly one <layout> child")
	}
	l, ok := nodes[0].(*astLayout)
	if !ok {
		return nil, p.errAt("<variant> requires a <layout> child")
	}
	return l, nil
}

// validateEachRef rejects a {{path}} token at parse time when its root
// name has no enclosing <each as="...">, per docs/ctml.md S:S8.3 -- an
// unbound reference is a document defect, not a runtime miss.
func (p *parser) validateEachRef(path string) error {
	root := path
	if i := strings.IndexByte(path, '.'); i >= 0 {
		root = path[:i]
	}
	for _, name := range p.eachStack {
		if name == root {
			return nil
		}
	}
	return p.errAt("unbound reference {{" + path + "}}: no enclosing each as=\"" + root + "\"")
}

func (p *parser) errAt(msg string) error {
	line, col := p.dec.InputPos()
	return &ParseError{Line: line, Col: col, Msg: msg, Cause: core.E("ctml.Parse", msg, nil)}
}

func (p *parser) wrapXMLErr(err error) error {
	line, col := p.dec.InputPos()
	msg := "malformed XML"
	// A clean end of input surfaces as the io.EOF sentinel; a document that
	// runs off the end mid-tag surfaces as a *xml.SyntaxError wrapping the
	// text "unexpected EOF" instead -- both are the same defect from a
	// .ctml author's point of view (nothing to close the document with).
	if err == io.EOF || strings.Contains(err.Error(), "unexpected EOF") {
		msg = "unexpected end of document"
	}
	return &ParseError{Line: line, Col: col, Msg: msg, Cause: core.E("ctml.Parse", msg, err)}
}

func attrValue(start xml.StartElement, name string) (string, bool) {
	for _, a := range start.Attr {
		if a.Name.Local == name {
			return a.Value, true
		}
	}
	return "", false
}

func slotLetter(tag string) byte {
	switch tag {
	case "h":
		return 'H'
	case "l":
		return 'L'
	case "c":
		return 'C'
	case "r":
		return 'R'
	case "f":
		return 'F'
	default:
		return 0
	}
}

// matchBindToken recognises a whole trimmed run of the shape {{ident(.ident)*}}.
// Anything else -- including a near-miss like {{oops!}} -- is left as
// literal text (S:S8.4 is lenient by design: prose that merely contains
// double braces should not become a parse error).
func matchBindToken(s string) (string, bool) {
	if len(s) < 4 || !strings.HasPrefix(s, "{{") || !strings.HasSuffix(s, "}}") {
		return "", false
	}
	inner := strings.TrimSpace(s[2 : len(s)-2])
	if !isValidPath(inner) {
		return "", false
	}
	return inner, true
}

func isValidPath(s string) bool {
	if s == "" {
		return false
	}
	for _, step := range strings.Split(s, ".") {
		if !isIdent(step) {
			return false
		}
	}
	return true
}

func isIdent(s string) bool {
	if s == "" {
		return false
	}
	for i := 0; i < len(s); i++ {
		c := s[i]
		letter := (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
		digit := c >= '0' && c <= '9'
		if i == 0 && !letter {
			return false
		}
		if i > 0 && !letter && !digit {
			return false
		}
	}
	return true
}
