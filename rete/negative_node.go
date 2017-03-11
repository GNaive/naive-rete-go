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

func (node NegativeNode) get_node_type() string {
	return NEGATIVE_NODE
}
func (node NegativeNode) get_parent() IReteNode {
	return node.parent
}
func (node NegativeNode) get_items() *list.List {
	return node.items
}
func (node *NegativeNode) get_children() *list.List {
	return node.children
}
func (node *NegativeNode) left_activation(t *Token, w *WME, b Binding) {
	new_token := make_token(node, t, w, b)
	node.items.PushBack(new_token)

	new_token.join_results = list.New()
	for e := node.amem.items.Front(); e != nil; e = e.Next() {
		w := e.Value.(*WME)
		if node.perform_join_tests(new_token, w) {
			jr := &NegativeJoinResult{
				owner: new_token,
				wme:   w,
			}
			new_token.join_results.PushBack(jr)
			w.negative_join_results = list.New()
			w.negative_join_results.PushBack(jr)
		}
	}
	if new_token.join_results.Len() == 0 {
		for e := node.children.Front(); e != nil; e = e.Next() {
			child := e.Value.(IReteNode)
			child.left_activation(new_token, nil, nil)
		}
	}
}
func (node *NegativeNode) right_activation(w *WME) {
	for e := node.items.Front(); e != nil; e = e.Next() {
		t := e.Value.(*Token)
		if node.perform_join_tests(t, w) {
			if t.join_results.Len() == 0 {
				t.delete_token_and_descendents()
			}
			jr := &NegativeJoinResult{
				owner: t,
				wme:   w,
			}
			t.join_results.PushBack(jr)
			w.negative_join_results = list.New()
			w.negative_join_results.PushBack(jr)
		}
	}
}
func (node *NegativeNode) perform_join_tests(t *Token, w *WME) bool {
	for e := node.tests.Front(); e != nil; e = e.Next() {
		test := e.Value.(*TestAtJoinNode)
		arg1 := w.fields[test.field_of_arg1]
		wme2 := t.get_wmes()[test.condition_number_of_arg2]
		arg2 := wme2.fields[test.field_of_arg2]
		if arg1 != arg2 {
			return false
		}
	}
	return true
}
