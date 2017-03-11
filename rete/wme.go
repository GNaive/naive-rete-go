package rete

import (
	"container/list"
	"fmt"
)

type WME struct {
	fields                [3]string
	alpha_mems            *list.List
	tokens                *list.List
	negative_join_results *list.List
}

func RemoveWME(w *WME) {
	for e := w.alpha_mems.Front(); e != nil; e = e.Next() {
		amem := e.Value.(*AlphaMemory)
		remove_by_value(amem.items, w)
	}
	for w.tokens != nil && w.tokens.Len() > 0 {
		e := w.tokens.Front()
		t := e.Value.(*Token)
		t.delete_token_and_descendents()
		w.tokens.Remove(e)
	}
	for e := w.negative_join_results.Front(); e != nil; e = e.Next() {
		jr := e.Value.(*NegativeJoinResult)
		remove_by_value(jr.owner.join_results, jr)
		if jr.owner.join_results.Len() == 0 {
			for i := jr.owner.node.get_children().Front(); i != nil; i = i.Next() {
				child := i.Value.(IReteNode)
				child.left_activation(jr.owner, nil, nil)
			}
		}
	}
}

func CreateWME(id, attr, value string) WME {
	return WME{
		fields:                [3]string{id, attr, value},
		alpha_mems:            list.New(),
		tokens:                list.New(),
		negative_join_results: list.New(),
	}
}

func (wme *WME) Equal(w *WME) bool {
	return wme.fields == w.fields
}

func (wme *WME) String() string {
	return fmt.Sprintf("%s", wme.fields)
}
