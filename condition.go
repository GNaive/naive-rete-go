package rete

type Condition struct {
	fields [3]string
	negative bool
}
func (cond Condition) contain (s string) int {
	for idx, v := range cond.fields {
		if v == s {
			return idx
		}
	}
	return -1
}
func (cond Condition) test_wme (w *WME) bool {
	for idx, v := range cond.fields {
		if v[0] == '$' {
			continue
		}
		if v != w.fields[idx] {
			return false
		}
	}
	return true
}