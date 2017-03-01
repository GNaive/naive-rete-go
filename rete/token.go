package rete

import (
	"container/list"
	"fmt"
	"strings"
)

type Token struct {
	parent       *Token
	wme          *WME
	node         IReteNode
	children     *list.List
	join_results *list.List
}

func (t *Token) get_wmes() []*WME {
	ret := []*WME{}
	_ws := list.New()
	_ws.PushFront(t.wme)
	for t.parent != nil {
		t = t.parent
		_ws.PushFront(t.wme)
	}
	for e := _ws.Front(); e != nil; e = e.Next() {
		ret = append(ret, e.Value.(*WME))
	}
	return ret
}
func make_token(node IReteNode, parent *Token, w *WME) *Token {
	tok := &Token{parent: parent, wme: w, node: node, children: list.New()}
	if parent != nil {
		parent.children.PushBack(tok)
	}
	if w != nil {
		w.tokens.PushBack(tok)
	}
	return tok
}
func (tok *Token) delete_token_and_descendents() {
	for tok.children != nil && tok.children.Len() > 0 {
		e := tok.children.Front()
		child := e.Value.(*Token)
		child.delete_token_and_descendents()
		tok.children.Remove(e)
	}
	remove_by_value(tok.node.get_items(), tok)
	if tok.wme != nil {
		//remove_by_value(tok.wme.tokens, tok)
	}
	if tok.parent != nil {
		//remove_by_value(tok.parent.children, tok)
	}
}
func (tok Token) String() string {
	ret := []string{}
	wmes := tok.get_wmes()
	for _, v := range wmes {
		s := fmt.Sprintf("%s", v)
		ret = append(ret, s)
	}
	return fmt.Sprintf("<Token %s>", strings.Join(ret, ", "))
}
