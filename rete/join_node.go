package rete

import "container/list"

type TestAtJoinNode struct {
	field_of_arg1            int
	condition_number_of_arg2 int
	field_of_arg2            int
}

type JoinNode struct {
	parent   IReteNode
	children *list.List
	amem     *AlphaMemory
	tests    *list.List
	has      *Has
}

func (node JoinNode) get_node_type() string {
	return JOIN_NODE
}
func (node JoinNode) get_items() *list.List {
	return nil
}
func (node JoinNode) get_parent() IReteNode {
	return node.parent
}
func (node JoinNode) get_children() *list.List {
	return node.children
}
func (node *JoinNode) right_activation(w *WME) {
	parent := node.parent
	// dummy join
	if parent.get_parent().get_node_type() == BETA_MEMORY_NODE {
		b := node.make_binding(w)
		for _e := node.children.Front(); _e != nil; _e = _e.Next() {
			child := _e.Value.(IReteNode)
			child.left_activation(nil, w, b)
		}
		return
	}
	for e := parent.get_items().Front(); e != nil; e = e.Next() {
		t := e.Value.(*Token)
		if node.perform_join_tests(t, w) {
			b := node.make_binding(w)
			for _e := node.children.Front(); _e != nil; _e = _e.Next() {
				child := _e.Value.(IReteNode)
				child.left_activation(t, w, b)
			}
		}
	}
}
func (node *JoinNode) left_activation(t *Token, w *WME, b Binding) {
	for e := node.amem.items.Front(); e != nil; e = e.Next() {
		w := e.Value.(*WME)
		if node.perform_join_tests(t, w) {
			b := node.make_binding(w)
			for _e := node.children.Front(); _e != nil; _e = _e.Next() {
				child := _e.Value.(IReteNode)
				child.left_activation(t, w, b)
			}
		}
	}
}
func (node *JoinNode) perform_join_tests(t *Token, w *WME) bool {
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
func (node *JoinNode) make_binding(w *WME) Binding {
	b := make(Binding)
	for idx, v := range node.has.fields {
		if is_var(v) {
			b[v] = w.fields[idx]
		}
	}
	return b
}
