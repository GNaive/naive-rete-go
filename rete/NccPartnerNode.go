package rete

import "container/list"

type NccPartnerNode struct {
	parent              IReteNode
	children            *list.List
	ncc_node            IReteNode
	number_of_conjuncts int
	new_result_buffer   *list.List
}

func (node NccPartnerNode) get_node_type() string {
	return NCC_PARTNER_NODE
}
func (node NccPartnerNode) get_parent() IReteNode {
	return node.parent
}
func (node NccPartnerNode) get_items() *list.List {
	return nil
}
func (node NccPartnerNode) get_children() *list.List {
	return node.children
}
func (node NccPartnerNode) left_activation(t *Token, w *WME, b Binding) {
	ncc_node := node.ncc_node
	new_result := make_token(node, t, w, b)
	owners_t := t
	owners_w := w
	for i := 1; i <= node.number_of_conjuncts; i++ {
		owners_w = owners_t.wme
		owners_t = owners_t.parent
	}
	for e := ncc_node.get_items().Front(); e != nil; e = e.Next() {
		item := e.Value.(*Token)
		if item.parent == owners_t && item.wme == owners_w {
			item.ncc_results.PushBack(item)
			new_result.owner = item
			item.delete_token_and_descendents()
			return
		}
	}
	node.new_result_buffer.PushBack(new_result)
}
func (node NccPartnerNode) right_activation(w *WME) {
}
