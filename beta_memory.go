package rete

import "container/list"

type BetaMemory struct {
	items    *list.List
	parent   IReteNode
	children *list.List
}

func (node BetaMemory) get_node_type() string {
	return BETA_MEMORY_NODE
}
func (node BetaMemory) get_items() *list.List {
	return node.items
}
func (node BetaMemory) get_parent() IReteNode {
	return node.parent
}
func (node BetaMemory) get_children() *list.List {
	return node.children
}
func (node BetaMemory) left_activation(t *Token, w *WME) {
	new_token := make_token(node, t, w)
	node.items.PushBack(&new_token)
	for e := node.children.Front(); e != nil; e = e.Next() {
		e.Value.(IReteNode).left_activation(&new_token, nil)
	}
}
func (node BetaMemory) right_activation(w *WME) {
}
