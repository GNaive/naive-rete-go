package rete

import "container/list"

type NccNode struct {
	parent   IReteNode
	children *list.List
	items    *list.List
	partner  *NccPartnerNode
}

func (node NccNode) GetNodeType() string {
	return NCC_NODE
}
func (node NccNode) GetParent() IReteNode {
	return node.parent
}
func (node NccNode) GetItems() *list.List {
	return node.items
}
func (node NccNode) GetChildren() *list.List {
	return node.children
}
func (node NccNode) LeftActivation(t *Token, w *WME, b Env) {
	new_token := make_token(node, t, w, b)
	node.items.PushBack(new_token)

	new_token.ncc_results = list.New()
	buffer := node.partner.new_result_buffer
	for e := buffer.Front(); e != nil; e = e.Next() {
		result := e.Value.(*Token)
		result.owner = new_token
		new_token.ncc_results.PushBack(result)
		buffer.Remove(e)
	}
	if new_token.ncc_results.Len() > 0 {
		return
	}
	for e := node.children.Front(); e != nil; e = e.Next() {
		e.Value.(IReteNode).LeftActivation(new_token, nil, nil)
	}
}
func (node NccNode) RightActivation(w *WME) {
}
