package rete

const (
	ClassName = iota
	Identifier
	Attribute
	Value
	NoTest
)

const (
	BetaMemoryNodeTy = "beta_memory"
	JoinNodeTy       = "join_node"
	NegativeNodeTy   = "negative_node"
	NccNodeTy        = "ncc_node"
	NccPartnerNodeTy = "ncc_parter_node"
	FilterNodeTy     = "filter_node"
)

var FIELDS = []int{ClassName, Identifier, Attribute, Value}
