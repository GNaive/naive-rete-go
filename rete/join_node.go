package rete

import "container/list"

type TestAtJoinNode struct {
	fieldOfArg1           int
	conditionNumberOfArg2 int
	fieldOfArg2           int
}

type JoinNode struct {
	parent   IReteNode
	children *list.List
	amem     *AlphaMemory
	tests    *list.List
	has      *Has
}

func (node JoinNode) GetNodeType() string {
	return JoinNodeTy
}
func (node JoinNode) GetItems() *list.List {
	return nil
}
func (node JoinNode) GetParent() IReteNode {
	return node.parent
}
func (node JoinNode) GetChildren() *list.List {
	return node.children
}
func (node *JoinNode) RightActivation(w *WME) {
	parent := node.parent
	// dummy join
	if parent.GetParent().GetNodeType() == BetaMemoryNodeTy {
		b := node.makeBinding(w)
		for _e := node.children.Front(); _e != nil; _e = _e.Next() {
			child := _e.Value.(IReteNode)
			child.LeftActivation(nil, w, b)
		}
		return
	}
	for e := parent.GetItems().Front(); e != nil; e = e.Next() {
		t := e.Value.(*Token)
		if node.performJoinTests(t, w) {
			b := node.makeBinding(w)
			for _e := node.children.Front(); _e != nil; _e = _e.Next() {
				child := _e.Value.(IReteNode)
				child.LeftActivation(t, w, b)
			}
		}
	}
}
func (node *JoinNode) LeftActivation(t *Token, w *WME, b Env) {
	for e := node.amem.items.Front(); e != nil; e = e.Next() {
		w := e.Value.(*WME)
		if node.performJoinTests(t, w) {
			b := node.makeBinding(w)
			for _e := node.children.Front(); _e != nil; _e = _e.Next() {
				child := _e.Value.(IReteNode)
				child.LeftActivation(t, w, b)
			}
		}
	}
}
func (node *JoinNode) performJoinTests(t *Token, w *WME) bool {
	for e := node.tests.Front(); e != nil; e = e.Next() {
		test := e.Value.(*TestAtJoinNode)
		arg1 := w.fields[test.fieldOfArg1]
		wme2 := t.get_wmes()[test.conditionNumberOfArg2]
		arg2 := wme2.fields[test.fieldOfArg2]
		if arg1 != arg2 {
			return false
		}
	}
	return true
}
func (node *JoinNode) makeBinding(w *WME) Env {
	b := make(Env)
	for idx, v := range node.has.fields {
		if isVar(v) {
			b[varKey(v)] = w.fields[idx]
		}
	}
	return b
}
