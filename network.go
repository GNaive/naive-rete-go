package rete

import (
	"container/list"
)

type IReteNode interface {
	get_node_type () string
	get_items () *list.List
	get_parent () IReteNode
	get_children () *list.List
	left_activation (t *Token, w *WME)
	right_activation (w *WME)
}

type Network struct {
	alpha_root *ConstantTestNode
	beta_root IReteNode
}
func CreateNetwork () Network {
	work_memory := &AlphaMemory{
		items: list.New(),
		successors: list.New(),
	}
	alpha_root := &ConstantTestNode{
		field_to_test: NO_TEST,
		field_must_equal: "",
		output_memory: work_memory,
		children: list.New(),
	}
	beta_root := &BetaMemory{
		items: list.New(),
		parent: nil,
		children: list.New(),
	}
	return Network{
		alpha_root:alpha_root,
		beta_root: beta_root,
	}
}
func CreateCondition(id, attr, value string) Condition {
	return Condition{
		fields: [3]string{id, attr, value},
		negative: false,
	}
}
func (n Network) AddProduction (lhs []Condition) IReteNode {
	current_node := n.build_or_share_network_for_conditions(n.beta_root, lhs, []Condition{})
	return n.build_or_share_beta_memory(current_node)
}
func (n Network) AddWME (w *WME) {
	n.alpha_root.activation(w)
}
func (n Network) build_or_share_network_for_conditions(
	parent IReteNode, conds []Condition, earlier_conds []Condition) IReteNode {
	current_node := parent
	conds_higher_up := earlier_conds
	for _, cond := range conds {
		if !cond.negative {
			current_node = n.build_or_share_beta_memory(current_node)
			tests := n.get_join_tests_from_condition(cond, conds_higher_up)
			am := n.build_or_share_alpha_memory(cond)
			current_node = n.build_or_share_join_node(current_node, am, tests)
		}
		conds_higher_up = append(conds_higher_up, cond)
	}
	return current_node
}
func (n Network) build_or_share_beta_memory(parent IReteNode) IReteNode {
	for e := parent.get_children().Front(); e != nil; e = e.Next() {
		if e.Value.(IReteNode).get_node_type() == BETA_MEMORY_NODE {
			return e.Value.(IReteNode)
		}
	}
	node := &BetaMemory{
		items: list.New(),
		parent: parent,
		children: list.New(),
	}
	parent.get_children().PushBack(node)
	n.update_new_node_with_matches_above(node)
	return node
}
func (n Network) build_or_share_join_node(parent IReteNode, amem *AlphaMemory, tests *list.List) IReteNode {
	for e := parent.get_children().Front(); e != nil; e = e.Next() {
		if e.Value.(IReteNode).get_node_type() != JOIN_NODE {
			continue
		}
		node := e.Value.(*JoinNode)
		if node.amem == amem && node.tests == tests {
			return node
		}
	}
	node := &JoinNode {
		parent: parent,
		children: list.New(),
		amem: amem,
		tests: tests,
	}
	parent.get_children().PushBack(node)
	amem.successors.PushBack(node)
	return node
}
func (n Network) build_or_share_alpha_memory(c Condition) *AlphaMemory {
	current_node := n.alpha_root
	for field, sym := range c.fields {
		if sym[0] != '$' {
			current_node = n.build_or_share_constant_test_node(current_node, field, sym)
		}
	}
	if current_node.output_memory != nil {
		return current_node.output_memory
	}
	am := &AlphaMemory{
		items: list.New(),
		successors: list.New(),
	}
	current_node.output_memory = am
	for e := n.alpha_root.output_memory.items.Front(); e != nil; e = e.Next() {
		w := e.Value.(*WME)
		if c.test_wme(w) {
			am.activation(w)
		}
	}
	return am
}
func (n Network) build_or_share_constant_test_node (
	parent *ConstantTestNode, field int, symbol string) *ConstantTestNode {
	for e := parent.children.Front(); e != nil; e = e.Next() {
		child := e.Value.(*ConstantTestNode)
		if child.field_to_test == field && child.field_must_equal == symbol {
			return child
		}
	}
	node := &ConstantTestNode {
		field_to_test: field,
		field_must_equal: symbol,
		output_memory: nil,
		children: list.New(),
	}
	parent.children.PushBack(node)
	return node
}
func (n Network) get_join_tests_from_condition(c Condition, earlier_conds []Condition) *list.List {
	ret := list.New()
	for v_field1, v := range c.fields {
		if v[0] != '$' {
			continue
		}
		for cond_idx, cond := range earlier_conds {
			v_field2 := cond.contain(v)
			if v_field2 == -1 || cond.negative {
				continue
			}
			node := &TestAtJoinNode{v_field1, cond_idx, v_field2}
			ret.PushBack(node)
		}
	}
	return ret
}
func (n Network) update_new_node_with_matches_above(node IReteNode) {
	parent := node.get_parent()
	if parent == nil {
		return
	}
	switch parent.get_node_type() {
	case BETA_MEMORY_NODE:
		for e := parent.get_items().Front(); e != nil; e = e.Next() {
			t := e.Value.(*Token)
			node.left_activation(t, nil)
		}
	case JOIN_NODE:
		parent := parent.(*JoinNode)
		saved_children := parent.children
		hack_children :=  list.New()
		hack_children.PushBack(node)
		parent.children = hack_children
		for e := parent.amem.items.Front(); e != nil; e = e.Next() {
			w := e.Value.(*WME)
			parent.right_activation(w)
		}
		parent.children = saved_children
	}
}
