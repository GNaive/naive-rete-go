package rete

import (
	"bytes"
	"container/list"
	"log"
	"runtime/debug"
)

type IReteNode interface {
	GetNodeType() string
	GetItems() *list.List
	GetParent() IReteNode
	GetChildren() *list.List
	LeftActivation(t *Token, w *WME, b Env)
	RightActivation(w *WME)
}

type Network struct {
	alphaRoot *ConstantTestNode
	betaRoot  IReteNode
	objects   Env // for rhs result
	PNodes    []*BetaMemory
	halt      bool
	LogBuf    *bytes.Buffer
}

func NewNetwork() *Network {
	workMemory := &AlphaMemory{
		items:      list.New(),
		successors: list.New(),
	}
	alphaRoot := &ConstantTestNode{
		fieldToTest:    NoTest,
		fieldMustEqual: "",
		outputMemory:   workMemory,
		children:       list.New(),
	}
	betaRoot := &BetaMemory{
		items:    list.New(),
		parent:   nil,
		children: list.New(),
	}
	return &Network{
		alphaRoot: alphaRoot,
		betaRoot:  betaRoot,
		objects:   make(Env),
		PNodes:    []*BetaMemory{},
		halt:      false,
		LogBuf:    &bytes.Buffer{},
	}
}

func (n *Network) AddObject(key string, obj interface{}) {
	n.objects[key] = obj
}

func (n Network) GetObjects() Env {
	return n.objects
}

func (n Network) GetObject(key string) interface{} {
	return n.objects[key]
}

func (n *Network) Halt() {
	n.halt = true
}

func (n *Network) ExecuteRules(env Env) (err error) {
	for _, pNode := range n.PNodes {
		for elem := pNode.GetItems().Front(); elem != nil; elem = elem.Next() {
			token := elem.Value.(*Token)
			if pNode.RHS == nil || len(pNode.RHS.tmpl) == 0 {
				continue
			}
			handler := env[pNode.RHS.tmpl]
			if handler == nil {
				continue
			}
			func() {
				defer func() {
					l := log.New(n.LogBuf, "RHS `"+pNode.RHS.tmpl+"` ", log.Lshortfile)
					if e := recover(); e != nil {
						l.Printf("%s %s", e, debug.Stack())
					}

				}()
				handler.(func(network *Network, token *Token))(
					n, token,
				)
			}()
			if n.halt {
				return nil
			}
		}
	}
	return nil
}

func (n *Network) AddProduction(lhs LHS, rhs RHS) *BetaMemory {
	currentNode := n.buildOrShareNetworkForConditions(n.betaRoot, lhs, LHS{})
	node := n.buildOrShareBetaMemory(currentNode)
	memory := node.(*BetaMemory)
	memory.RHS = &rhs
	n.PNodes = append(n.PNodes, memory)
	return memory
}

func (n *Network) AddProductionFromXML(s string) (result []*BetaMemory, err error) {
	ps, err := FromXML(s)
	if err != nil {
		return result, err
	}
	for _, p := range ps {
		result = append(result, n.AddProduction(p.lhs, p.rhs))
	}
	return result, nil
}

func (n *Network) AddWME(w *WME) {
	n.alphaRoot.activation(w)
}

func (n Network) buildOrShareNetworkForConditions(
	parent IReteNode, rule LHS, earlierConds LHS) IReteNode {
	currentNode := parent
	condsHigherUp := earlierConds
	for _, cond := range rule.items {
		switch cond.(type) {
		case Has:
			cond := cond.(Has)
			if !cond.negative {
				currentNode = n.buildOrShareBetaMemory(currentNode)
				tests := n.getJoinTestsFromCondition(cond, condsHigherUp)
				am := n.buildOrShareAlphaMemory(cond)
				currentNode = n.buildOrShareJoinNode(currentNode, am, tests, &cond)
			} else {
				tests := n.getJoinTestsFromCondition(cond, condsHigherUp)
				am := n.buildOrShareAlphaMemory(cond)
				currentNode = n.buildOrShareNegativeNode(currentNode, am, tests)
			}
		case Filter:
			cond := cond.(Filter)
			currentNode = n.buildOrShareFilterNode(currentNode, cond)
		case LHS:
			cond := cond.(LHS)
			if cond.negative {
				currentNode = n.buildOrShareNccNodes(currentNode, cond, condsHigherUp)
			}
		}
		condsHigherUp.items = append(condsHigherUp.items, cond)
	}
	return currentNode
}

func (n Network) buildOrShareFilterNode(parent IReteNode, f Filter) IReteNode {
	for e := parent.GetChildren().Front(); e != nil; e = e.Next() {
		child := e.Value.(IReteNode)
		if child.GetNodeType() == FilterNodeTy {
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

func (n Network) buildOrShareNccNodes(parent IReteNode, ncc LHS, earlier LHS) IReteNode {
	bottomOfSubnetwork := n.buildOrShareNetworkForConditions(parent, ncc, earlier)
	for e := parent.GetChildren().Front(); e != nil; e = e.Next() {
		child := e.Value.(IReteNode)
		if child.GetNodeType() == NccNodeTy {
			child := child.(*NccNode)
			if child.partner.parent == bottomOfSubnetwork {
				return child
			}
		}
	}
	nccNode := &NccNode{
		parent:   parent,
		children: list.New(),
		items:    list.New(),
	}
	nccPartnerNode := &NccPartnerNode{
		parent:            bottomOfSubnetwork,
		children:          list.New(),
		newResultBuffer:   list.New(),
		numberOfConjuncts: len(ncc.items),
		nccNode:           nccNode,
	}
	nccNode.partner = nccPartnerNode
	parent.GetChildren().PushBack(nccNode)
	bottomOfSubnetwork.GetChildren().PushBack(nccPartnerNode)
	n.updateNewNodeWithMatchesAbove(nccNode)
	n.updateNewNodeWithMatchesAbove(nccPartnerNode)
	return nccNode
}

func (n Network) buildOrShareBetaMemory(parent IReteNode) IReteNode {
	for e := parent.GetChildren().Front(); e != nil; e = e.Next() {
		if e.Value.(IReteNode).GetNodeType() == BetaMemoryNodeTy {
			return e.Value.(IReteNode)
		}
	}
	node := &BetaMemory{
		items:    list.New(),
		parent:   parent,
		children: list.New(),
	}
	parent.GetChildren().PushBack(node)
	n.updateNewNodeWithMatchesAbove(node)
	return node
}

func (n Network) buildOrShareJoinNode(
	parent IReteNode, amem *AlphaMemory, tests *list.List, h *Has) IReteNode {
	for e := parent.GetChildren().Front(); e != nil; e = e.Next() {
		if e.Value.(IReteNode).GetNodeType() != JoinNodeTy {
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

func (n Network) buildOrShareNegativeNode(parent IReteNode, amem *AlphaMemory, tests *list.List) IReteNode {
	for e := parent.GetChildren().Front(); e != nil; e = e.Next() {
		if e.Value.(IReteNode).GetNodeType() != NegativeNodeTy {
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
	n.updateNewNodeWithMatchesAbove(node)
	return node
}

func (n Network) buildOrShareAlphaMemory(c Has) *AlphaMemory {
	currentNode := n.alphaRoot
	for field, sym := range c.fields {
		if !isVar(sym) {
			currentNode = n.buildOrShareConstantTestNode(currentNode, field, sym)
		}
	}
	if currentNode.outputMemory != nil {
		return currentNode.outputMemory
	}
	am := &AlphaMemory{
		items:      list.New(),
		successors: list.New(),
	}
	currentNode.outputMemory = am
	for e := n.alphaRoot.outputMemory.items.Front(); e != nil; e = e.Next() {
		w := e.Value.(*WME)
		if c.testWme(w) {
			am.activation(w)
		}
	}
	return am
}

func (n Network) buildOrShareConstantTestNode(
	parent *ConstantTestNode, field int, symbol string) *ConstantTestNode {
	for e := parent.children.Front(); e != nil; e = e.Next() {
		child := e.Value.(*ConstantTestNode)
		if child.fieldToTest == field && child.fieldMustEqual == symbol {
			return child
		}
	}
	node := &ConstantTestNode{
		fieldToTest:    field,
		fieldMustEqual: symbol,
		outputMemory:   nil,
		children:       list.New(),
	}
	parent.children.PushBack(node)
	return node
}

func (n Network) getJoinTestsFromCondition(c Has, earlierConds LHS) *list.List {
	ret := list.New()
	for vField1, v := range c.fields {
		if !isVar(v) {
			continue
		}
		for condIdx, cond := range earlierConds.items {
			switch cond.(type) {
			case Has:
				cond := cond.(Has)
				vField2 := cond.contain(v)
				if vField2 == -1 || cond.negative {
					continue
				}
				node := &TestAtJoinNode{vField1, condIdx, vField2}
				ret.PushBack(node)
			}
		}
	}
	return ret
}

func (n Network) updateNewNodeWithMatchesAbove(node IReteNode) {
	parent := node.GetParent()
	if parent == nil {
		return
	}
	switch parent.GetNodeType() {
	case BetaMemoryNodeTy:
		for e := parent.GetItems().Front(); e != nil; e = e.Next() {
			t := e.Value.(*Token)
			node.LeftActivation(t, nil, nil)
		}
	case JoinNodeTy:
		parent := parent.(*JoinNode)
		savedChildren := parent.children
		hackChildren := list.New()
		hackChildren.PushBack(node)
		parent.children = hackChildren
		for e := parent.amem.items.Front(); e != nil; e = e.Next() {
			w := e.Value.(*WME)
			parent.RightActivation(w)
		}
		parent.children = savedChildren
	case NegativeNodeTy:
		for e := parent.GetItems().Front(); e != nil; e = e.Next() {
			t := e.Value.(*Token)
			node.LeftActivation(t, nil, nil)
		}
	case NccNodeTy:
		for e := parent.GetItems().Front(); e != nil; e = e.Next() {
			t := e.Value.(*Token)
			node.LeftActivation(t, nil, nil)
		}
	}
}
