package rete

import (
	"container/list"
	"strings"
)

type FilterNode struct {
	parent   IReteNode
	children *list.List
	tmpl     string
}

func (node FilterNode) GetNodeType() string {
	return FILTER_NODE
}
func (node FilterNode) GetItems() *list.List {
	return nil
}
func (node FilterNode) GetParent() IReteNode {
	return node.parent
}
func (node FilterNode) GetChildren() *list.List {
	return node.children
}
func (node *FilterNode) RightActivation(w *WME) {
}
func (node *FilterNode) LeftActivation(t *Token, w *WME, b Binding) {
	all_binding := t.AllBinding()
	for k, v := range b {
		all_binding[k] = v
	}
	code := node.tmpl
	for k, v := range all_binding {
		code = strings.Replace(code, k, v, -1)
	}
	result := EvalFromString(code)
	if result != nil && !result.(bool) {
		return
	}
	for e := node.children.Front(); e != nil; e = e.Next() {
		child := e.Value.(IReteNode)
		child.LeftActivation(t, w, b)
	}
}
