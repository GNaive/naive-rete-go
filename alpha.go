package rete

import (
	"container/list"
)

type AlphaMemory struct {
	items *list.List
	successors *list.List
}

type ConstantTestNode struct {
	field_to_test int
	field_must_equal string
	output_memory *AlphaMemory
	children *list.List
}

func (node ConstantTestNode) activation(w *WME) {
	if node.field_to_test != NO_TEST {
		if w.fields[node.field_to_test] != node.field_must_equal {
			return
		}
	}
	if node.output_memory != nil {
		node.output_memory.activation(w)
	}
	for e := node.children.Front(); e != nil; e = e.Next() {
		e.Value.(*ConstantTestNode).activation(w)
	}
}

func (node *AlphaMemory) activation(w *WME) {
	node.items.PushBack(w)
	w.alpha_mems.PushBack(node)
	for e := node.successors.Front(); e != nil; e = e.Next() {
		e.Value.(IReteNode).right_activation(w)
	}
}