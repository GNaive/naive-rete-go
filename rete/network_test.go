package rete

import (
	"fmt"
	"testing"
)

func TestNetworkAddWME(t *testing.T) {
	n := NewNetwork()
	c0 := NewHas("Object", "$x", "on", "$y")
	c1 := NewHas("Object", "$y", "left_of", "$z")
	am0 := n.build_or_share_alpha_memory(c0)
	am1 := n.build_or_share_alpha_memory(c1)
	wmes := []*WME{
		NewWME("Object", "B1", "on", "B2"),
		NewWME("Object", "B2", "left_of", "B3"),
		NewWME("Object", "B2", "on", "table"),
	}
	for idx := range wmes {
		n.AddWME(wmes[idx])
	}
	if am0.items.Len() != 2 || am1.items.Len() != 1 {
		t.Error("add wme error")
	}
}

func TestCase0(t *testing.T) {
	n := NewNetwork()
	c0 := NewHas("Object", "$x", "on", "$y")
	c1 := NewHas("Object", "$y", "left_of", "$z")
	c2 := NewHas("Object", "$z", "color", "red")
	p := n.AddProduction(NewRule(c0, c1, c2), nil)
	wmes := []*WME{
		NewWME("Object", "B1", "on", "B2"),
		NewWME("Object", "B1", "on", "B3"),
		NewWME("Object", "B1", "color", "red"),
		NewWME("Object", "B2", "on", "table"),
		NewWME("Object", "B2", "left_of", "B3"),
		NewWME("Object", "B2", "color", "blue"),
		NewWME("Object", "B3", "left_of", "B4"),
		NewWME("Object", "B3", "on", "table"),
		NewWME("Object", "B3", "color", "red"),
	}
	for idx := range wmes {
		n.AddWME(wmes[idx])
	}

	if p.GetItems().Len() != 1 {
		t.Error()
	}
	expect := "<Token [Object B1 on B2], [Object B2 left_of B3], [Object B3 color red]>"
	for e := p.GetItems().Front(); e != nil; e = e.Next() {
		tok := e.Value.(*Token)
		if fmt.Sprint(tok) != expect {
			t.Error(tok)
		}
		x, y, z := tok.GetBinding("$x"), tok.GetBinding("$y"), tok.GetBinding("$z")
		if x != "B1" || y != "B2" || z != "B3" {
			t.Error("error binding")
		}
	}
}

func TestNegativeNode(t *testing.T) {
	n := NewNetwork()
	c0 := NewHas("Object","$x", "on", "$y")
	c1 := NewNeg("Object", "$y", "color", "blue")
	p := n.AddProduction(NewRule(c0, c1), nil)

	wmes := []*WME{
		NewWME("Object", "B1", "on", "B2"),
		NewWME("Object", "B1", "on", "B3"),
		NewWME("Object", "B2", "color", "blue"),
		NewWME("Object", "B3", "color", "red"),
	}
	for idx := range wmes {
		n.AddWME(wmes[idx])
	}

	expect := "<Token [Object B1 on B3], <nil>>"
	for e := p.GetItems().Front(); e != nil; e = e.Next() {
		tok := e.Value.(*Token)
		if fmt.Sprint(tok) != expect {
			t.Error("error result")
		}
		x, y := tok.GetBinding("$x"), tok.GetBinding("$y")
		if x != "B1" || y != "B3" {
			t.Error("error binding")
		}
	}
}

func TestNccNode(t *testing.T) {
	n := NewNetwork()
	c0 := NewHas("Object","$x", "on", "$y")
	c1 := NewHas("Object","$y", "left_of", "$z")
	c2 := NewHas("Object","$z", "color", "red")
	c3 := NewHas("Object","$z", "on", "$w")
	p := n.AddProduction(NewRule(c0, c1, NewNccRule(c2, c3)), nil)
	wmes := []*WME{
		NewWME("Object","B1", "on", "B2"),
		NewWME("Object","B1", "on", "B3"),
		NewWME("Object", "B1", "color", "red"),
		NewWME("Object", "B2", "on", "table"),
		NewWME("Object", "B2", "left_of", "B3"),
		NewWME("Object", "B2", "color", "blue"),
		NewWME("Object", "B3", "left_of", "B4"),
		NewWME("Object", "B3", "on", "table"),
		NewWME("Object", "B3", "color", "red"),
	}
	for idx := range wmes {
		n.AddWME(wmes[idx])
	}
	expect := "<Token [Object B1 on B3], [Object B3 left_of B4], <nil>>"
	for e := p.GetItems().Front(); e != nil; e = e.Next() {
		tok := e.Value.(*Token)
		if fmt.Sprint(tok) != expect {
			t.Error(tok)
		}
		x, y := tok.GetBinding("$x"), tok.GetBinding("$y")
		if x != "B1" || y != "B3" {
			t.Error("error binding")
		}
	}
}

func TestFromXML(t *testing.T) {
	data := `
	<?xml version="1.0"?>
	<data>
	    <production>
		<lhs>
		    <has classname="user" identifier="$uid" attribute="id" value="$uid"/>
		    <has classname="spu" identifier="1" attribute="quantity" value="$quantity"/>
		    <filter><![CDATA[$quantity > 1]]></filter>
		</lhs>
		<rhs action="block"></rhs>
	    </production>
	    <production>
		<lhs>
		    <has classname="user" identifier="$uid" attribute="id" value="$uid"/>
		    <has classname="spu" identifier="2" attribute="quantity" value="$quantity"/>
		    <filter><![CDATA[$quantity > 10]]></filter>
		</lhs>
		<rhs action="block"></rhs>
	    </production>
	</data>
	`
	n := NewNetwork()
	pnodes, err := n.AddProductionFromXML(data)
	if err != nil {
		t.Error(err)
		return
	}
	p0 := pnodes[0]
	p1 := pnodes[1]

	wmes := []*WME{
		NewWME("user", "100001", "id", "100001"),
		NewWME("spu", "1", "quantity", "2"),
		NewWME("spu", "2", "quantity", "6"),
	}
	for idx := range wmes {
		n.AddWME(wmes[idx])
	}
	expect := "<Token [user 100001 id 100001], [spu 1 quantity 2]>"
	for e := p0.GetItems().Front(); e != nil; e = e.Next() {
		tok := e.Value.(*Token)
		if val := fmt.Sprint(tok); val != expect {
			t.Error(val)
		}
	}
	if p0.GetExecuteParam("action") != "block" {t.Error()}
	if p0.GetExecuteParam("dummy") != nil {t.Error()}
	if p1.GetItems().Len() > 0 {
		t.Error()
	}
}
