package rete

import "container/list"

type NccNode struct {
	parent   IReteNode
	children *list.List
	items    *list.List
	partner  *NccPartnerNode
}

func (node *NccNode) GetNodeType() string {
	return NccNodeTy
}
func (node *NccNode) GetParent() IReteNode {
	return node.parent
}
func (node *NccNode) GetItems() *list.List {
	return node.items
}
func (node *NccNode) GetChildren() *list.List {
	return node.children
}
func (node *NccNode) LeftActivation(t *Token, w *WME, b Env) {
	newToken := makeToken(node, t, w, b)
	node.items.PushBack(newToken)

	newToken.nccResults = list.New()
	buffer := node.partner.newResultBuffer
	for e := buffer.Front(); e != nil; e = e.Next() {
		result := e.Value.(*Token)
		result.owner = newToken
		newToken.nccResults.PushBack(result)
		buffer.Remove(e)
	}
	if newToken.nccResults.Len() > 0 {
		return
	}
	for e := node.children.Front(); e != nil; e = e.Next() {
		e.Value.(IReteNode).LeftActivation(newToken, nil, nil)
	}
}
func (node NccNode) RightActivation(w *WME) {
}
