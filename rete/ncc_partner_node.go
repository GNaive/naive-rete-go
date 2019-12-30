package rete

import "container/list"

type NccPartnerNode struct {
	parent            IReteNode
	children          *list.List
	nccNode           IReteNode
	numberOfConjuncts int
	newResultBuffer   *list.List
}

func (node NccPartnerNode) GetNodeType() string {
	return NccPartnerNodeTy
}
func (node NccPartnerNode) GetParent() IReteNode {
	return node.parent
}
func (node NccPartnerNode) GetItems() *list.List {
	return nil
}
func (node NccPartnerNode) GetChildren() *list.List {
	return node.children
}
func (node NccPartnerNode) LeftActivation(t *Token, w *WME, b Env) {
	nccNode := node.nccNode
	newResult := makeToken(node, t, w, b)
	ownersT := t
	ownersW := w
	for i := 1; i <= node.numberOfConjuncts; i++ {
		ownersW = ownersT.wme
		ownersT = ownersT.parent
	}
	for e := nccNode.GetItems().Front(); e != nil; e = e.Next() {
		item := e.Value.(*Token)
		if item.parent == ownersT && item.wme == ownersW {
			item.nccResults.PushBack(item)
			newResult.owner = item
			item.deleteTokenAndDescendents()
			return
		}
	}
	node.newResultBuffer.PushBack(newResult)
}
func (node NccPartnerNode) RightActivation(w *WME) {
}
