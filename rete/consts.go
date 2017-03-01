package rete

const (
	IDENTIFIER = iota
	ATTRIBUTE
	VALUE
	NO_TEST
)

const (
	BETA_MEMORY_NODE = "beta_memory"
	JOIN_NODE        = "join_node"
	NEGATIVE_NODE    = "negative_node"
)

var FIELDS = []int{IDENTIFIER, ATTRIBUTE, VALUE}
