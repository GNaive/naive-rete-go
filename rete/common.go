package rete

func isVar(v string) bool {
	return len(v) > 0 && v[0] == '$'
}

func varKey(v string) string {
	if isVar(v) {
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
	tmpl  string
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

func (has Has) testWme(w *WME) bool {
	for idx, v := range has.fields {
		if isVar(v) {
			continue
		}
		if v != w.fields[idx] {
			return false
		}
	}
	return true
}

func NewHas(className, id, attr, value string) Has {
	return Has{
		fields:   [4]string{className, id, attr, value},
		negative: false,
	}
}

func NewNeg(className, id, attr, value string) Has {
	return Has{
		fields:   [4]string{className, id, attr, value},
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
		tmpl:  "",
		Extra: make(map[string]interface{}),
	}
}

func NewNccRule(items ...interface{}) LHS {
	return LHS{
		items:    items,
		negative: true,
	}
}
