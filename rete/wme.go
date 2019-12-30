package rete

import (
	"container/list"
	"fmt"
)

type WME struct {
	fields              [4]string
	alphaMems           *list.List
	tokens              *list.List
	negativeJoinResults *list.List
}

func RemoveWME(w *WME) {
	for e := w.alphaMems.Front(); e != nil; e = e.Next() {
		amem := e.Value.(*AlphaMemory)
		removeByValue(amem.items, w)
	}
	for w.tokens != nil && w.tokens.Len() > 0 {
		e := w.tokens.Front()
		t := e.Value.(*Token)
		t.deleteTokenAndDescendents()
		w.tokens.Remove(e)
	}
	for e := w.negativeJoinResults.Front(); e != nil; e = e.Next() {
		jr := e.Value.(*NegativeJoinResult)
		removeByValue(jr.owner.joinResults, jr)
		if jr.owner.joinResults.Len() == 0 {
			for i := jr.owner.node.GetChildren().Front(); i != nil; i = i.Next() {
				child := i.Value.(IReteNode)
				child.LeftActivation(jr.owner, nil, nil)
			}
		}
	}
}

func NewWME(className, id, attr, value string) *WME {
	return &WME{
		fields:              [4]string{className, id, attr, value},
		alphaMems:           list.New(),
		tokens:              list.New(),
		negativeJoinResults: list.New(),
	}
}

func (wme *WME) Equal(w *WME) bool {
	return wme.fields == w.fields
}

func (wme *WME) String() string {
	return fmt.Sprintf("%s", wme.fields)
}
