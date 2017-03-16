package rete

type Production struct {
	lhs Rule
	rhs map[string]interface{}
}
type Rule struct {
	items    []interface{}
	negative bool
}
type Has struct {
	fields   [4]string
	negative bool
}
type Filter struct {
	tmpl string
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
func NewHas(class_name, id, attr, value string) Has {
	return Has{
		fields:   [4]string{class_name, id, attr, value},
		negative: false,
	}
}
func NewNeg(class_name, id, attr, value string) Has {
	return Has{
		fields:   [4]string{class_name, id, attr, value},
		negative: true,
	}
}
func NewRule(items ...interface{}) Rule {
	return Rule{
		items: items,
	}
}
func NewNccRule(items ...interface{}) Rule {
	return Rule{
		items:    items,
		negative: true,
	}
}
