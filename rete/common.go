package rete

func is_var(v string) bool {
	return len(v) > 0 && v[0] == '$'
}
