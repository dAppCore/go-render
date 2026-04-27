package html

import (
	"hash/fnv"
	"reflect"
	"strconv"
)

const (
	defaultGrammarImprintMaxDepth   = 128
	defaultGrammarImprintMaxPathLen = 256
)

// Stamp is a structural fingerprint for a node position in the HLCRF tree.
type Stamp struct {
	Path string
	Hash uint64
	Tags []string
}

// GrammarImprint classifies node structure without reading rendered content.
type GrammarImprint struct {
	maxDepth   int
	maxPathLen int
}

// Imprint returns a structural stamp for node without rendering it or reading
// text/raw content.
func (g *GrammarImprint) Imprint(node Node, ctx Context) Stamp {
	if isNilNode(node) {
		return Stamp{}
	}

	w := grammarImprintWalker{
		ctx:        &ctx,
		maxDepth:   g.configuredMaxDepth(),
		maxPathLen: g.configuredMaxPathLen(),
	}
	return w.imprint(node, "", 0)
}

func (g *GrammarImprint) configuredMaxDepth() int {
	if g != nil && g.maxDepth > 0 {
		return g.maxDepth
	}
	return defaultGrammarImprintMaxDepth
}

func (g *GrammarImprint) configuredMaxPathLen() int {
	if g != nil && g.maxPathLen > 0 {
		return g.maxPathLen
	}
	return defaultGrammarImprintMaxPathLen
}

type grammarImprintWalker struct {
	ctx        *Context
	maxDepth   int
	maxPathLen int
}

func (w grammarImprintWalker) imprint(node Node, path string, depth int) Stamp {
	if isNilNode(node) {
		return Stamp{}
	}
	if depth >= w.maxDepth {
		return w.stamp(node, path, true)
	}

	switch n := node.(type) {
	case *Layout:
		return w.imprintLayout(n, path, depth)
	case *Responsive:
		return w.imprintResponsive(n, path, depth)
	case *ifNode:
		if n == nil || n.cond == nil || n.node == nil || !n.cond(w.ctx) {
			return w.emptyStamp(node, path)
		}
		return w.imprint(n.node, path, depth+1)
	case *unlessNode:
		if n == nil || n.cond == nil || n.node == nil || n.cond(w.ctx) {
			return w.emptyStamp(node, path)
		}
		return w.imprint(n.node, path, depth+1)
	case *entitledNode:
		if n == nil || n.node == nil || w.ctx == nil || w.ctx.Entitlements == nil || !w.ctx.Entitlements(n.feature) {
			return w.emptyStamp(node, path)
		}
		return w.imprint(n.node, path, depth+1)
	case *switchNode:
		if n == nil || n.selector == nil || n.cases == nil {
			return w.emptyStamp(node, path)
		}
		child := n.cases[n.selector(w.ctx)]
		if child == nil {
			return w.emptyStamp(node, path)
		}
		return w.imprint(child, path, depth+1)
	default:
		return w.stamp(node, path, false)
	}
}

func (w grammarImprintWalker) imprintLayout(l *Layout, path string, depth int) Stamp {
	if l == nil {
		return Stamp{}
	}

	slotCounts := make(map[byte]int)
	slotOrdinal := 0

	for i := range len(l.variant) {
		slot := l.variant[i]
		if _, ok := slotRegistry[slot]; !ok {
			continue
		}

		count := slotOrdinal
		slotOrdinal++

		children := l.slots[slot]
		if len(children) == 0 {
			continue
		}

		if path == "" {
			count = slotCounts[slot]
			slotCounts[slot] = count + 1
		}

		blockPath := w.layoutBlockPath(path, slot, count)
		for childIndex, child := range children {
			if isNilNode(child) {
				continue
			}
			childPath := w.joinPath(blockPath, strconv.Itoa(childIndex))
			return w.imprint(child, childPath, depth+1)
		}

		return w.slotStamp(blockPath)
	}

	return w.stamp(l, path, false)
}

func (w grammarImprintWalker) imprintResponsive(r *Responsive, path string, depth int) Stamp {
	if r == nil {
		return Stamp{}
	}
	for _, variant := range r.variants {
		if variant.layout == nil {
			continue
		}
		return w.imprint(variant.layout, path, depth+1)
	}
	return w.stamp(r, path, false)
}

func (w grammarImprintWalker) stamp(node Node, path string, truncated bool) Stamp {
	path = w.normalizedPath(node, path)
	childCount := structuralChildCount(node, w.ctx)
	tags := structuralTags(childCount, truncated, structuralEmpty(node, w.ctx))
	nodeType := structuralNodeType(node)

	return Stamp{
		Path: path,
		Hash: grammarStampHash(path, nodeType, childCount),
		Tags: tags,
	}
}

func (w grammarImprintWalker) emptyStamp(node Node, path string) Stamp {
	path = w.normalizedPath(node, path)
	return Stamp{
		Path: path,
		Hash: grammarStampHash(path, structuralNodeType(node), 0),
		Tags: []string{"empty"},
	}
}

func (w grammarImprintWalker) slotStamp(path string) Stamp {
	return Stamp{
		Path: path,
		Hash: grammarStampHash(path, "layout-slot", 0),
		Tags: []string{"empty"},
	}
}

func (w grammarImprintWalker) normalizedPath(node Node, path string) string {
	if path != "" || isCoordinateContainer(node) {
		return w.clampPath(path)
	}
	return "0"
}

func isCoordinateContainer(node Node) bool {
	switch node.(type) {
	case *Layout, *Responsive:
		return true
	default:
		return false
	}
}

func (w grammarImprintWalker) layoutBlockPath(base string, slot byte, rendered int) string {
	if base == "" {
		if rendered == 0 {
			return w.clampPath(string(slot))
		}
		return w.joinPath(string(slot), strconv.Itoa(rendered))
	}
	if rendered == 0 {
		return w.clampPath(base)
	}
	return w.joinPath(base, strconv.Itoa(rendered))
}

func (w grammarImprintWalker) joinPath(path, coord string) string {
	if coord == "" {
		return w.clampPath(path)
	}
	if path == "" {
		return w.clampPath(coord)
	}
	if len(path)+1+len(coord) <= w.maxPathLen {
		return path + "." + coord
	}
	if len(path) >= w.maxPathLen {
		return path[:w.maxPathLen]
	}
	remaining := w.maxPathLen - len(path)
	if remaining <= 1 {
		return path
	}
	return path + "." + coord[:remaining-1]
}

func (w grammarImprintWalker) clampPath(path string) string {
	if len(path) <= w.maxPathLen {
		return path
	}
	return path[:w.maxPathLen]
}

func structuralChildCount(node Node, ctx *Context) int {
	switch n := node.(type) {
	case *rawNode, *textNode:
		return 0
	case *elNode:
		if n == nil {
			return 0
		}
		return countNodes(n.children)
	case *ifNode:
		if n == nil || n.cond == nil || n.node == nil || !n.cond(ctx) {
			return 0
		}
		return 1
	case *unlessNode:
		if n == nil || n.cond == nil || n.node == nil || n.cond(ctx) {
			return 0
		}
		return 1
	case *entitledNode:
		if n == nil || n.node == nil || ctx == nil || ctx.Entitlements == nil || !ctx.Entitlements(n.feature) {
			return 0
		}
		return 1
	case *switchNode:
		if n == nil {
			return 0
		}
		return countMapNodes(n.cases)
	case *Layout:
		if n == nil {
			return 0
		}
		return countLayoutChildren(n)
	case *Responsive:
		if n == nil {
			return 0
		}
		count := 0
		for _, variant := range n.variants {
			if variant.layout != nil {
				count++
			}
		}
		return count
	default:
		return structuralEachChildCount(node)
	}
}

func structuralEachChildCount(node Node) int {
	n, ok := node.(interface{ isNilHTMLNode() bool })
	if !ok || n.isNilHTMLNode() {
		return 0
	}

	value := reflect.ValueOf(node)
	if !value.IsValid() || value.Kind() != reflect.Pointer || value.IsNil() {
		return 0
	}
	elem := value.Elem()
	if !elem.IsValid() || elem.Kind() != reflect.Struct {
		return 0
	}
	items := elem.FieldByName("items")
	if items.IsValid() && items.Kind() == reflect.Slice {
		return items.Len()
	}
	return 0
}

func countLayoutChildren(l *Layout) int {
	if l == nil {
		return 0
	}
	count := 0
	for i := range len(l.variant) {
		slot := l.variant[i]
		if _, ok := slotRegistry[slot]; !ok {
			continue
		}
		count += countNodes(l.slots[slot])
	}
	return count
}

func countNodes(nodes []Node) int {
	count := 0
	for _, node := range nodes {
		if !isNilNode(node) {
			count++
		}
	}
	return count
}

func countMapNodes(nodes map[string]Node) int {
	count := 0
	for _, node := range nodes {
		if !isNilNode(node) {
			count++
		}
	}
	return count
}

func structuralEmpty(node Node, ctx *Context) bool {
	switch n := node.(type) {
	case *ifNode:
		return n == nil || n.cond == nil || n.node == nil || !n.cond(ctx)
	case *unlessNode:
		return n == nil || n.cond == nil || n.node == nil || n.cond(ctx)
	case *entitledNode:
		return n == nil || n.node == nil || ctx == nil || ctx.Entitlements == nil || !ctx.Entitlements(n.feature)
	case *switchNode:
		if n == nil || n.selector == nil || n.cases == nil {
			return true
		}
		return isNilNode(n.cases[n.selector(ctx)])
	case *Layout:
		return countLayoutChildren(n) == 0
	case *Responsive:
		return structuralChildCount(n, ctx) == 0
	default:
		t := reflect.TypeOf(node)
		if t == nil || t.Kind() != reflect.Pointer || t.Elem().Kind() != reflect.Struct {
			return false
		}
		if _, ok := node.(interface{ isNilHTMLNode() bool }); !ok {
			return false
		}
		return structuralEachChildCount(node) == 0
	}
}

func structuralTags(childCount int, truncated, empty bool) []string {
	if truncated {
		return []string{"branch", "truncated"}
	}
	if empty {
		return []string{"empty"}
	}
	if childCount == 0 {
		return []string{"leaf"}
	}
	return []string{"branch"}
}

func structuralNodeType(node Node) string {
	switch n := node.(type) {
	case *rawNode:
		return "raw"
	case *textNode:
		return "text"
	case *elNode:
		if n == nil {
			return "element"
		}
		return "element:" + n.tag
	case *ifNode:
		return "if"
	case *unlessNode:
		return "unless"
	case *entitledNode:
		return "entitled"
	case *switchNode:
		return "switch"
	case *Layout:
		return "layout"
	case *Responsive:
		return "responsive"
	default:
		t := reflect.TypeOf(node)
		if t == nil {
			return ""
		}
		return t.String()
	}
}

func grammarStampHash(path, nodeType string, childCount int) uint64 {
	// Sanctioned exception: core.Hash64 is unavailable in the current core
	// module, so RFC §4's deterministic 64-bit structural hash uses stdlib FNV.
	h := fnv.New64()
	_, _ = h.Write([]byte(path))
	_, _ = h.Write([]byte{0})
	_, _ = h.Write([]byte(nodeType))
	_, _ = h.Write([]byte{0})
	_, _ = h.Write([]byte(strconv.Itoa(childCount)))
	return h.Sum64()
}
