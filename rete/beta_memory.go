package rete

import (
	"container/list"
	"reflect"
	"errors"
)

type BetaMemory struct {
	items    *list.List
	parent   IReteNode
	children *list.List
	RHS      *RHS
}

func (node BetaMemory) GetNodeType() string {
	return BETA_MEMORY_NODE
}

func (node BetaMemory) GetItems() *list.List {
	return node.items
}

func (node BetaMemory) GetParent() IReteNode {
	return node.parent
}

func (node BetaMemory) GetChildren() *list.List {
	return node.children
}

func (node *BetaMemory) LeftActivation(t *Token, w *WME, b Env) {
	new_token := make_token(node, t, w, b)
	node.items.PushBack(new_token)
	for e := node.children.Front(); e != nil; e = e.Next() {
		e.Value.(IReteNode).LeftActivation(new_token, nil, nil)
	}
}

func (node BetaMemory) RightActivation(w *WME) {
}

func (node BetaMemory) GetExecuteParam(s string) interface{} {
	return node.RHS.Extra[s]
}

func (node BetaMemory) Eval(t *Token, env Env) (result []reflect.Value, err error) {
	if node.RHS == nil || len(node.RHS.tmpl) == 0 {
		err = errors.New("no tmpl to eval")
		return
	}
	if t == nil {
		err = errors.New("token is nil")
		return
	}
	all_binding := t.AllBinding()
	for k, v := range all_binding {
		env[k] = v
	}
	return EvalFromString(node.RHS.tmpl, env)
}

func (node BetaMemory) PopToken() *Token {
	e := node.items.Front()
	if e == nil {
		return nil
	}
	node.items.Remove(e)
	return e.Value.(*Token)
}
