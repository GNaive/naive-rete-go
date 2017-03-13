package rete

import (
	"container/list"
)

type IReteNode interface {
	GetNodeType() string
	GetItems() *list.List
	GetParent() IReteNode
	GetChildren() *list.List
	LeftActivation(t *Token, w *WME, b Binding)
	RightActivation(w *WME)
}

type Network struct {
	alpha_root *ConstantTestNode
	beta_root  IReteNode
}

func CreateNetwork() Network {
	work_memory := &AlphaMemory{
		items:      list.New(),
		successors: list.New(),
	}
	alpha_root := &ConstantTestNode{
		field_to_test:    NO_TEST,
		field_must_equal: "",
		output_memory:    work_memory,
		children:         list.New(),
	}
	beta_root := &BetaMemory{
		items:    list.New(),
		parent:   nil,
		children: list.New(),
	}
	return Network{
		alpha_root: alpha_root,
		beta_root:  beta_root,
	}
}

func (n Network) AddProduction(lhs Rule, rhs map[string]interface{}) *BetaMemory {
	current_node := n.build_or_share_network_for_conditions(n.beta_root, lhs, Rule{})
	pnode := n.build_or_share_beta_memory(current_node)
	pnode.(*BetaMemory).execute_params = rhs
	return pnode.(*BetaMemory)
}
func (n Network) AddProductionFromXML(s string) []*BetaMemory {
	result := []*BetaMemory{}
	ps := FromXML(s)
	for _, p := range ps {
		result = append(result, n.AddProduction(p.lhs, p.rhs))
	}
	return result
}
func (n Network) AddWME(w *WME) {
	n.alpha_root.activation(w)
}
func (n Network) build_or_share_network_for_conditions(
	parent IReteNode, rule Rule, earlier_conds Rule) IReteNode {
	current_node := parent
	conds_higher_up := earlier_conds
	for _, cond := range rule.items {
		switch cond.(type) {
		case Has:
			cond := cond.(Has)
			if !cond.negative {
				current_node = n.build_or_share_beta_memory(current_node)
				tests := n.get_join_tests_from_condition(cond, conds_higher_up)
				am := n.build_or_share_alpha_memory(cond)
				current_node = n.build_or_share_join_node(current_node, am, tests, &cond)
			} else {
				tests := n.get_join_tests_from_condition(cond, conds_higher_up)
				am := n.build_or_share_alpha_memory(cond)
				current_node = n.build_or_share_negative_node(current_node, am, tests)
			}
		case Filter:
			cond := cond.(Filter)
			current_node = n.build_or_share_filter_node(current_node, cond)
		case Rule:
			cond := cond.(Rule)
			if cond.negative {
				current_node = n.build_or_share_ncc_nodes(current_node, cond, conds_higher_up)
			}
		}
		conds_higher_up.items = append(conds_higher_up.items, cond)
	}
	return current_node
}
func (n Network) build_or_share_filter_node(parent IReteNode, f Filter) IReteNode {
	for e := parent.GetChildren().Front(); e != nil; e = e.Next() {
		child := e.Value.(IReteNode)
		if child.GetNodeType() == FILTER_NODE {
			child := child.(*FilterNode)
			if child.tmpl == f.tmpl {
				return child
			}
		}
	}
	filter_node := &FilterNode{
		parent:   parent,
		children: list.New(),
		tmpl:     f.tmpl,
	}
	parent.GetChildren().PushBack(filter_node)
	return filter_node
}
func (n Network) build_or_share_ncc_nodes(parent IReteNode, ncc Rule, earlier Rule) IReteNode {
	bottom_of_subnetwork := n.build_or_share_network_for_conditions(parent, ncc, earlier)
	for e := parent.GetChildren().Front(); e != nil; e = e.Next() {
		child := e.Value.(IReteNode)
		if child.GetNodeType() == NCC_NODE {
			child := child.(*NccNode)
			if child.partner.parent == bottom_of_subnetwork {
				return child
			}
		}
	}
	ncc_node := &NccNode{
		parent:   parent,
		children: list.New(),
		items:    list.New(),
	}
	ncc_partner_node := &NccPartnerNode{
		parent:              bottom_of_subnetwork,
		children:            list.New(),
		new_result_buffer:   list.New(),
		number_of_conjuncts: len(ncc.items),
		ncc_node:            ncc_node,
	}
	ncc_node.partner = ncc_partner_node
	parent.GetChildren().PushBack(ncc_node)
	bottom_of_subnetwork.GetChildren().PushBack(ncc_partner_node)
	n.update_new_node_with_matches_above(ncc_node)
	n.update_new_node_with_matches_above(ncc_partner_node)
	return ncc_node
}
func (n Network) build_or_share_beta_memory(parent IReteNode) IReteNode {
	for e := parent.GetChildren().Front(); e != nil; e = e.Next() {
		if e.Value.(IReteNode).GetNodeType() == BETA_MEMORY_NODE {
			return e.Value.(IReteNode)
		}
	}
	node := &BetaMemory{
		items:    list.New(),
		parent:   parent,
		children: list.New(),
	}
	parent.GetChildren().PushBack(node)
	n.update_new_node_with_matches_above(node)
	return node
}
func (n Network) build_or_share_join_node(
	parent IReteNode, amem *AlphaMemory, tests *list.List, h *Has) IReteNode {
	for e := parent.GetChildren().Front(); e != nil; e = e.Next() {
		if e.Value.(IReteNode).GetNodeType() != JOIN_NODE {
			continue
		}
		node := e.Value.(*JoinNode)
		if node.amem == amem && node.tests == tests {
			return node
		}
	}
	node := &JoinNode{
		parent:   parent,
		children: list.New(),
		amem:     amem,
		tests:    tests,
		has:      h,
	}
	parent.GetChildren().PushBack(node)
	amem.successors.PushBack(node)
	return node
}
func (n Network) build_or_share_negative_node(parent IReteNode, amem *AlphaMemory, tests *list.List) IReteNode {
	for e := parent.GetChildren().Front(); e != nil; e = e.Next() {
		if e.Value.(IReteNode).GetNodeType() != NEGATIVE_NODE {
			continue
		}
		node := e.Value.(*NegativeNode)
		if node.amem == amem && node.tests == tests {
			return node
		}
	}
	node := &NegativeNode{
		parent:   parent,
		children: list.New(),
		amem:     amem,
		tests:    tests,
		items:    list.New(),
	}
	parent.GetChildren().PushBack(node)
	amem.successors.PushBack(node)
	n.update_new_node_with_matches_above(node)
	return node
}
func (n Network) build_or_share_alpha_memory(c Has) *AlphaMemory {
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
		items:      list.New(),
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
func (n Network) build_or_share_constant_test_node(
	parent *ConstantTestNode, field int, symbol string) *ConstantTestNode {
	for e := parent.children.Front(); e != nil; e = e.Next() {
		child := e.Value.(*ConstantTestNode)
		if child.field_to_test == field && child.field_must_equal == symbol {
			return child
		}
	}
	node := &ConstantTestNode{
		field_to_test:    field,
		field_must_equal: symbol,
		output_memory:    nil,
		children:         list.New(),
	}
	parent.children.PushBack(node)
	return node
}
func (n Network) get_join_tests_from_condition(c Has, earlier_conds Rule) *list.List {
	ret := list.New()
	for v_field1, v := range c.fields {
		if !is_var(v) {
			continue
		}
		for cond_idx, cond := range earlier_conds.items {
			switch cond.(type) {
			case Has:
				cond := cond.(Has)
				v_field2 := cond.contain(v)
				if v_field2 == -1 || cond.negative {
					continue
				}
				node := &TestAtJoinNode{v_field1, cond_idx, v_field2}
				ret.PushBack(node)
			}
		}
	}
	return ret
}
func (n Network) update_new_node_with_matches_above(node IReteNode) {
	parent := node.GetParent()
	if parent == nil {
		return
	}
	switch parent.GetNodeType() {
	case BETA_MEMORY_NODE:
		for e := parent.GetItems().Front(); e != nil; e = e.Next() {
			t := e.Value.(*Token)
			node.LeftActivation(t, nil, nil)
		}
	case JOIN_NODE:
		parent := parent.(*JoinNode)
		saved_children := parent.children
		hack_children := list.New()
		hack_children.PushBack(node)
		parent.children = hack_children
		for e := parent.amem.items.Front(); e != nil; e = e.Next() {
			w := e.Value.(*WME)
			parent.RightActivation(w)
		}
		parent.children = saved_children
	case NEGATIVE_NODE:
		for e := parent.GetItems().Front(); e != nil; e = e.Next() {
			t := e.Value.(*Token)
			node.LeftActivation(t, nil, nil)
		}
	case NCC_NODE:
		for e := parent.GetItems().Front(); e != nil; e = e.Next() {
			t := e.Value.(*Token)
			node.LeftActivation(t, nil, nil)
		}
	}
}
