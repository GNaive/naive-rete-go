package rete

import (
	"container/list"
	"testing"
)

func Test_make_token(t *testing.T) {
	node := BetaMemory{
		items:    list.New(),
		parent:   nil,
		children: list.New(),
	}
	w := &WME{
		fields:    [4]string{"Object", "B1", "on", "table"},
		alphaMems: list.New(),
		tokens:    list.New(),
	}
	token := makeToken(&node, nil, w, nil)
	if token.node.GetNodeType() != BetaMemoryNodeTy {
		t.Error("token node type error")
	}
	if token.wme != w {
		t.Error("token wme error")
	}
}
