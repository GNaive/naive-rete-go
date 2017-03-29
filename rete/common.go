package rete

func is_var(v string) bool {
	return len(v) > 0 && v[0] == '$'
}

func var_key(v string) string {
	if is_var(v) {
		return v[1:]
	}
	return ""
}

type Production struct {
	lhs LHS
	rhs RHS
}

type LHS struct {
	items    []interface{}
	negative bool
}

type RHS struct {
	tmpl string
	Extra map[string]interface{}
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

func NewLHS(items ...interface{}) LHS {
	return LHS{
		items: items,
	}
}

func NewRHS() RHS {
	return RHS{
		tmpl: "",
		Extra: make(map[string]interface{}),
	}
}

func NewNccRule(items ...interface{}) LHS {
	return LHS{
		items:    items,
		negative: true,
	}
}
