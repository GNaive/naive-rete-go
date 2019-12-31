package rete

import (
	"container/list"
	"testing"
)

func Test_beta_memory_items(t *testing.T) {
	bm := BetaMemory{
		items: list.New(),
	}
	newToken := &Token{}
	bm.GetItems().PushBack(newToken)
	e := bm.GetItems().Front()
	token := e.Value.(*Token)
	if token != newToken {
		t.Error("token error")
	}
}
