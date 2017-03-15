package rete

import (
	"fmt"
	"testing"
)

func TestNetworkAddWME(t *testing.T) {
	n := CreateNetwork()
	c0 := CreateHas("$x", "on", "$y")
	c1 := CreateHas("$y", "left_of", "$z")
	am0 := n.build_or_share_alpha_memory(c0)
	am1 := n.build_or_share_alpha_memory(c1)
	wmes := []WME{
		CreateWME("B1", "on", "B2"),
		CreateWME("B2", "left_of", "B3"),
		CreateWME("B2", "on", "table"),
	}
	for idx := range wmes {
		n.AddWME(&wmes[idx])
	}
	if am0.items.Len() != 2 || am1.items.Len() != 1 {
		t.Error("add wme error")
	}
}

func TestCase0(t *testing.T) {
	n := CreateNetwork()
	c0 := CreateHas("$x", "on", "$y")
	c1 := CreateHas("$y", "left_of", "$z")
	c2 := CreateHas("$z", "color", "red")
	p := n.AddProduction(CreateRule(c0, c1, c2), nil)
	wmes := []WME{
		CreateWME("B1", "on", "B2"),
		CreateWME("B1", "on", "B3"),
		CreateWME("B1", "color", "red"),
		CreateWME("B2", "on", "table"),
		CreateWME("B2", "left_of", "B3"),
		CreateWME("B2", "color", "blue"),
		CreateWME("B3", "left_of", "B4"),
		CreateWME("B3", "on", "table"),
		CreateWME("B3", "color", "red"),
	}
	for idx := range wmes {
		n.AddWME(&wmes[idx])
	}

	// am0 := n.build_or_share_alpha_memory(c0)
	// j0 := am0.successors.Front().Value.(*JoinNode)
	// b1 := n.build_or_share_beta_memory(j0)

	// am1 := n.build_or_share_alpha_memory(c1)
	// j1 := am1.successors.Front().Value.(*JoinNode)
	// b2 := n.build_or_share_beta_memory(j1)

	// am2 := n.build_or_share_alpha_memory(c2)
	// j2 := am2.successors.Front().Value.(*JoinNode)
	// b3 := n.build_or_share_beta_memory(j2)

	// fmt.Println(am0.items, am1.items, am2.items)
	// fmt.Println(b1.GetItems(), b2.GetItems(), b3.GetItems())
	expect := "<Token [B1 on B2], [B2 left_of B3], [B3 color red]>"
	for e := p.GetItems().Front(); e != nil; e = e.Next() {
		tok := e.Value.(*Token)
		if fmt.Sprint(tok) != expect {
			t.Error("error result")
		}
		x, y, z := tok.GetBinding("$x"), tok.GetBinding("$y"), tok.GetBinding("$z")
		if x != "B1" || y != "B2" || z != "B3" {
			t.Error("error binding")
		}
	}
}

func TestNegativeNode(t *testing.T) {
	n := CreateNetwork()
	c0 := CreateHas("$x", "on", "$y")
	c1 := CreateNeg("$y", "color", "blue")
	p := n.AddProduction(CreateRule(c0, c1), nil)

	wmes := []WME{
		CreateWME("B1", "on", "B2"),
		CreateWME("B1", "on", "B3"),
		CreateWME("B2", "color", "blue"),
		CreateWME("B3", "color", "red"),
	}
	for idx := range wmes {
		n.AddWME(&wmes[idx])
	}

	expect := "<Token [B1 on B3], <nil>>"
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
	n := CreateNetwork()
	c0 := CreateHas("$x", "on", "$y")
	c1 := CreateHas("$y", "left_of", "$z")
	c2 := CreateHas("$z", "color", "red")
	c3 := CreateHas("$z", "on", "$w")
	p := n.AddProduction(CreateRule(c0, c1, CreateNccRule(c2, c3)), nil)
	wmes := []WME{
		CreateWME("B1", "on", "B2"),
		CreateWME("B1", "on", "B3"),
		CreateWME("B1", "color", "red"),
		CreateWME("B2", "on", "table"),
		CreateWME("B2", "left_of", "B3"),
		CreateWME("B2", "color", "blue"),
		CreateWME("B3", "left_of", "B4"),
		CreateWME("B3", "on", "table"),
		CreateWME("B3", "color", "red"),
	}
	for idx := range wmes {
		n.AddWME(&wmes[idx])
	}
	expect := "<Token [B1 on B3], [B3 left_of B4], <nil>>"
	for e := p.GetItems().Front(); e != nil; e = e.Next() {
		tok := e.Value.(*Token)
		if fmt.Sprint(tok) != expect {
			t.Error("error")
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
		    <has identifier="user" attribute="id" value="$uid"/>
		    <has identifier="spu:1" attribute="quantity" value="$quantity"/>
		    <filter><![CDATA[$quantity > 1]]></filter>
		</lhs>
		<rhs action="block"></rhs>
	    </production>
	    <production>
		<lhs>
		    <has identifier="user" attribute="id" value="$uid"/>
		    <has identifier="spu:2" attribute="quantity" value="$quantity"/>
		    <filter><![CDATA[$quantity > 10]]></filter>
		</lhs>
		<rhs action="block"></rhs>
	    </production>
	</data>
	`
	n := CreateNetwork()
	pnodes, err := n.AddProductionFromXML(data)
	if err != nil {
		fmt.Println(err)
		return
	}
	p0 := pnodes[0]
	p1 := pnodes[1]

	wmes := []WME{
		CreateWME("user", "id", "100001"),
		CreateWME("spu:1", "quantity", "2"),
		CreateWME("spu:2", "quantity", "6"),
	}
	for idx := range wmes {
		n.AddWME(&wmes[idx])
	}
	expect := "<Token [user id 100001], [spu:1 quantity 2]>"
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
