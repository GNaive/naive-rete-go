package rete

import (
	"container/list"
	"testing"
)

func Test_beta_memory_items(t *testing.T) {
	bm := BetaMemory{
		items: list.New(),
	}
	new_token := &Token{}
	bm.GetItems().PushBack(new_token)
	e := bm.GetItems().Front()
	token := e.Value.(*Token)
	if token != new_token {
		t.Error("token error")
	}
}
