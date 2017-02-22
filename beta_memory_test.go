package rete

import (
	"testing"
	"container/list"
)

func Test_beta_memory_items(t *testing.T) {
	bm := BetaMemory{
		items: list.New(),
	}
	new_token := &Token{}
	bm.get_items().PushBack(new_token)
	e := bm.get_items().Front()
	token := e.Value.(*Token)
	if token != new_token {
		t.Error("token error")
	}
}
