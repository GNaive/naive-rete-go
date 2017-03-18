package rete

import (
	"container/list"
	"github.com/beevik/etree"
	"errors"
	"encoding/json"
	"fmt"
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

func FromXML(s string) (result []Production, err error) {
	doc := etree.NewDocument()
	err = doc.ReadFromString(s)
	if err != nil {
		return result, err
	}
	root := doc.Root()
	if root == nil {
		return result, errors.New("Not XML")
	}

	for _, ep := range root.ChildElements() {
		if ep.Tag != "production" {
			continue
		}
		p := Production{
			rhs: make(map[string]interface{}),
		}
		for idx, hand := range ep.ChildElements() {
			if idx == 0 {
				p.lhs = XMLParseLHS(hand)
			} else if idx == 1 {
				for _, attr := range hand.Attr {
					p.rhs[attr.Key] = attr.Value
				}
			}
		}
		result = append(result, p)
	}
	return result, nil
}

func XMLParseLHS(root *etree.Element) Rule {
	r := NewRule()
	for _, e := range root.ChildElements() {
		switch e.Tag {
		case "has", "neg":
			class_name, identity, attribute, value := "", "", "", ""
			for _, attr := range e.Attr {
				if attr.Key == "classname" {
					class_name = attr.Value
				} else if attr.Key == "identifier" {
					identity = attr.Value
				} else if attr.Key == "attribute" {
					attribute = attr.Value
				} else if attr.Key == "value" {
					value = attr.Value
				}
			}
			var has Has
			if e.Tag == "has" {
				has = NewHas(class_name, identity, attribute, value)
			} else {
				has = NewNeg(class_name, identity, attribute, value)
			}
			r.items = append(r.items, has)
		case "filter":
			f := Filter{tmpl: e.Text()}
			r.items = append(r.items, f)
		case "ncc":
			_rule := XMLParseLHS(e)
			_rule.negative = true
			r.items = append(r.items, _rule)
		}
	}
	return r
}

func FromJSON (s string) (r []Production, err error) {
	root := make(map[string]interface{})
	err = json.Unmarshal([]byte(s), &root)
	if err != nil {
		return r, err
	}
	if root["productions"] == nil {
		return r, errors.New("no productions")
	}
	ps, ok := root["productions"].([]interface{})
	if !ok {
		return r, errors.New("productions not List")
	}
	for _, p := range ps {
		production := Production{}
		p, ok := p.(map[string]interface{})
		if !ok {
			message := fmt.Sprintf("production not Object: %s", p)
			return r, errors.New(message)
		}
		production.rhs, ok = p["rhs"].(map[string]interface{})
		if !ok {
			message := fmt.Sprintf("rhs not Object: %s", p["rhs"])
			return r, errors.New(message)
		}
		lhs, ok := p["lhs"].([]interface{})
		if !ok {
			message := fmt.Sprintf("lhs not List: %s", p["lhs"])
			return r, errors.New(message)
		}
		production.lhs, err = JSONParseLHS(lhs)
		if err != nil {
			return r, err
		}
		r = append(r, production)
	}
	return r, err
}

func JSONParseLHS (lhs []interface{}) (r Rule, err error) {
	for _, e := range lhs {
		cond, ok := e.(map[string]interface{})
		if !ok {
			message := fmt.Sprintf("lhs element not Object: %s", e)
			return r, errors.New(message)
		}
		switch cond["tag"] {
		case "has", "neg":
			class, ok_0 := cond["classname"].(string)
			id, ok_1 := cond["identifier"].(string)
			attr, ok_2 := cond["attribute"].(string)
			value, ok_3 := cond["value"].(string)
			if !ok_0 || !ok_1 || !ok_2 || !ok_3 {
				message := fmt.Sprintf("condition missing fields: %s", cond)
				return r, errors.New(message)
			}
			if cond["tag"] == "has" {
				r.items = append(r.items, NewHas(class, id, attr, value))
			} else {
				r.items = append(r.items, NewNeg(class, id, attr, value))
			}
		case "filter":
			tmpl, ok := cond["tmpl"].(string)
			if !ok {
				message := fmt.Sprintf("filter tmpl not string: %s", cond)
				return r, errors.New(message)
			}
			r.items = append(r.items, Filter{tmpl:tmpl})
		case "ncc":
			ncc, ok := cond["items"].([]interface{})
			if !ok {
				message := fmt.Sprintf("lhs not List: %s", cond["items"])
				return r, errors.New(message)
			}
			_rule, err := JSONParseLHS(ncc)
			_rule.negative = true
			if err != nil {
				return r, err
			}
			r.items = append(r.items, _rule)
		default:
			message := fmt.Sprintf("tag error: %s", cond["tag"])
			return r, errors.New(message)

		}
	}
	return r, err
}
