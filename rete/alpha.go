package rete

import (
	"container/list"
)

type AlphaMemory struct {
	items      *list.List
	successors *list.List
}

type ConstantTestNode struct {
	fieldToTest    int
	fieldMustEqual string
	outputMemory   *AlphaMemory
	children       *list.List
}

func (node ConstantTestNode) activation(w *WME) {
	if node.fieldToTest != NoTest {
		if w.fields[node.fieldToTest] != node.fieldMustEqual {
			return
		}
	}
	if node.outputMemory != nil {
		node.outputMemory.activation(w)
	}
	for e := node.children.Front(); e != nil; e = e.Next() {
		e.Value.(*ConstantTestNode).activation(w)
	}
}

func (node *AlphaMemory) activation(w *WME) {
	node.items.PushBack(w)
	w.alphaMems.PushBack(node)
	for e := node.successors.Front(); e != nil; e = e.Next() {
		e.Value.(IReteNode).RightActivation(w)
	}
}
