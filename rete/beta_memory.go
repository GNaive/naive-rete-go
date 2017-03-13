package rete

import "container/list"

type BetaMemory struct {
	items          *list.List
	parent         IReteNode
	children       *list.List
	execute_params map[string]interface{} //use BetaMemory as PNode
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
func (node *BetaMemory) LeftActivation(t *Token, w *WME, b Binding) {
	new_token := make_token(node, t, w, b)
	node.items.PushBack(new_token)
	for e := node.children.Front(); e != nil; e = e.Next() {
		e.Value.(IReteNode).LeftActivation(new_token, nil, nil)
	}
}
func (node BetaMemory) RightActivation(w *WME) {
}
func (node BetaMemory) GetExecuteParam(s string) interface{} {
	return node.execute_params[s]
}
