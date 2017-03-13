package rete

import (
	"container/list"
	"github.com/beevik/etree"
)

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

func remove_by_value(l *list.List, value interface{}) bool {
	if e := contain(l, value); e != nil {
		l.Remove(e)
		return true
	}
	return false
}

func FromXML(s string) []Production {
	doc := etree.NewDocument()
	doc.ReadFromString(s)
	root := doc.Root()

	result := []Production{}
	for _, ep := range root.ChildElements() {
		if ep.Tag != "production" {continue}
		p := Production{
			rhs: make(map[string]interface{}),
		}
		for idx, hand := range ep.ChildElements() {
			if idx == 0 {
				p.lhs = parse_lhs(hand)
			} else if idx == 1 {
				for _, attr := range hand.Attr {
					p.rhs[attr.Key] = attr.Value
				}
			}
		}
		result = append(result, p)
	}
	return result
}

func parse_lhs(root *etree.Element) Rule {
	r := CreateRule()
	for _, e := range root.ChildElements() {
		switch e.Tag {
		case "has":
			identity, attribute, value := "", "", ""
			for _, attr := range e.Attr {
				if attr.Key == "identifier" {
					identity = attr.Value
				} else if attr.Key == "attribute" {
					attribute = attr.Value
				} else if attr.Key == "value" {
					value = attr.Value
				}
			}
			has := CreateHas(identity, attribute, value)
			r.items = append(r.items, has)
		case "filter":
			f := Filter{tmpl: e.Text()}
			r.items = append(r.items, f)
		}
	}
	return r
}
