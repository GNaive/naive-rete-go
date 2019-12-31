package rete

import (
	"container/list"
)

type FilterNode struct {
	parent   IReteNode
	children *list.List
	tmpl     string
}

func (node FilterNode) GetNodeType() string {
	return FilterNodeTy
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
func (node *FilterNode) LeftActivation(t *Token, w *WME, b Env) {
	all_binding := t.AllBinding()
	for k, v := range b {
		all_binding[k] = v
	}
	result, err := EvalFromString(node.tmpl, all_binding)
	if err != nil || len(result) == 0 {
		return
	}
	if !result[0].Bool() {
		return
	}
	for e := node.children.Front(); e != nil; e = e.Next() {
		child := e.Value.(IReteNode)
		child.LeftActivation(t, w, b)
	}
}
