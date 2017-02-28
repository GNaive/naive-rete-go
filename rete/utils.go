package rete

import "container/list"

func contain(l *list.List, value interface{}) *list.Element {
	if l == nil {
		return nil
	}
	for e := l.Front(); e != nil; e = e.Next() {
		if e.Value == value {
			return e
		}
	}
	return nil
}

func remove_by_value(l *list.List, value interface{}) {
	if e := contain(l, value); e != nil {
		l.Remove(e)
	}
}
