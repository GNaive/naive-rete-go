package rete

import (
	"container/list"
)

type BetaMemory struct {
	items    *list.List
	parent   IReteNode
	children *list.List
	RHS      *RHS
}

func (node BetaMemory) GetNodeType() string {
	return BetaMemoryNodeTy
}

func (node BetaMemory) GetItems() *list.List {
	return node.items
}

func (node BetaMemory) GetParent() IReteNode {
	return node.parent
}

func (node BetaMemory) GetChildren() *list.List {
	return node.children
}

func (node *BetaMemory) LeftActivation(t *Token, w *WME, b Env) {
	newToken := makeToken(node, t, w, b)
	node.items.PushBack(newToken)
	for e := node.children.Front(); e != nil; e = e.Next() {
		e.Value.(IReteNode).LeftActivation(newToken, nil, nil)
	}
}

func (node BetaMemory) RightActivation(w *WME) {
}

func (node BetaMemory) GetExecuteParam(s string) interface{} {
	if node.RHS == nil {
		return nil
	}
	return node.RHS.Extra[s]
}

func (node BetaMemory) PopToken() *Token {
	e := node.items.Front()
	if e == nil {
		return nil
	}
	node.items.Remove(e)
	return e.Value.(*Token)
}
