package rete

import (
	"container/list"
)

type NegativeJoinResult struct {
	owner *Token
	wme   *WME
}
type NegativeNode struct {
	parent   IReteNode
	children *list.List
	items    *list.List
	amem     *AlphaMemory
	tests    *list.List
}

func (node NegativeNode) GetNodeType() string {
	return NegativeNodeTy
}
func (node NegativeNode) GetParent() IReteNode {
	return node.parent
}
func (node NegativeNode) GetItems() *list.List {
	return node.items
}
func (node *NegativeNode) GetChildren() *list.List {
	return node.children
}
func (node *NegativeNode) LeftActivation(t *Token, w *WME, b Env) {
	newToken := makeToken(node, t, w, b)
	node.items.PushBack(newToken)

	newToken.joinResults = list.New()
	for e := node.amem.items.Front(); e != nil; e = e.Next() {
		w := e.Value.(*WME)
		if node.perform_join_tests(newToken, w) {
			jr := &NegativeJoinResult{
				owner: newToken,
				wme:   w,
			}
			newToken.joinResults.PushBack(jr)
			w.negativeJoinResults = list.New()
			w.negativeJoinResults.PushBack(jr)
		}
	}
	if newToken.joinResults.Len() == 0 {
		for e := node.children.Front(); e != nil; e = e.Next() {
			child := e.Value.(IReteNode)
			child.LeftActivation(newToken, nil, nil)
		}
	}
}
func (node *NegativeNode) RightActivation(w *WME) {
	for e := node.items.Front(); e != nil; e = e.Next() {
		t := e.Value.(*Token)
		if node.perform_join_tests(t, w) {
			if t.joinResults.Len() == 0 {
				t.deleteTokenAndDescendents()
			}
			jr := &NegativeJoinResult{
				owner: t,
				wme:   w,
			}
			t.joinResults.PushBack(jr)
			w.negativeJoinResults = list.New()
			w.negativeJoinResults.PushBack(jr)
		}
	}
}
func (node *NegativeNode) perform_join_tests(t *Token, w *WME) bool {
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
