package rete

import (
	"container/list"
	"fmt"
)

type WME struct {
	fields     [3]string
	alpha_mems *list.List
	tokens     *list.List
}

func RemoveWme(w *WME) {
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
}

func CreateWME(id, attr, value string) WME {
	return WME{
		fields:     [3]string{id, attr, value},
		alpha_mems: list.New(),
		tokens:     list.New(),
	}
}

func (wme *WME) Equal(w *WME) bool {
	return wme.fields == w.fields
}

func (wme *WME) String() string {
	return fmt.Sprintf("%s", wme.fields)
}
