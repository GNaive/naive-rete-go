package rete

type Rule struct {
	items    []interface{}
	negative bool
}
type Has struct {
	fields   [3]string
	negative bool
}

func (has Has) contain(s string) int {
	for idx, v := range has.fields {
		if v == s {
			return idx
		}
	}
	return -1
}
func (has Has) test_wme(w *WME) bool {
	for idx, v := range has.fields {
		if is_var(v) {
			continue
		}
		if v != w.fields[idx] {
			return false
		}
	}
	return true
}
func CreateHas(id, attr, value string) Has {
	return Has{
		fields:   [3]string{id, attr, value},
		negative: false,
	}
}
func CreateNeg(id, attr, value string) Has {
	return Has{
		fields:   [3]string{id, attr, value},
		negative: true,
	}
}
func CreateRule(items ...interface{}) Rule {
	return Rule{
		items: items,
	}
}
func CreateNccRule(items ...interface{}) Rule {
	return Rule{
		items:    items,
		negative: true,
	}
}
