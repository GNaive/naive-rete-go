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
	env := make(Env)
	env["F"] = func (network *Network, token *Token) {
		x := token.GetBinding("x")
		y := token.GetBinding("y")
		z := token.GetBinding("z")
		dummy := token.GetRHSParam("dummy")
		ret := fmt.Sprintf("%s, %s, %s, %v", x, y, z, dummy)
		network.AddObject("result", ret)
		network.Halt()
	}

	c0 := NewHas("Object", "$x", "on", "$y")
	c1 := NewHas("Object", "$y", "left_of", "$z")
	c2 := NewHas("Object", "$z", "color", "red")
	m := make(map[string]interface{})
	m["dummy"] = 1
	n.AddProduction(NewLHS(c0, c1, c2), RHS{
		tmpl: `F`,
		Extra: m,
	})
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
	err := n.ExecuteRules(env)
	if err != nil {
		t.Error(err)
	}
	if n.GetObject("result") != "B1, B2, B3, 1" {
		t.Error(n.GetObject("result"))
	}

}

func TestNegativeNode(t *testing.T) {
	n := NewNetwork()
	c0 := NewHas("Object","$x", "on", "$y")
	c1 := NewNeg("Object", "$y", "color", "blue")
	p := n.AddProduction(NewLHS(c0, c1), NewRHS())

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
		x, y := tok.GetBinding("x"), tok.GetBinding("y")
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
	p := n.AddProduction(NewLHS(c0, c1, NewNccRule(c2, c3)), NewRHS())
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
		x, y := tok.GetBinding("x"), tok.GetBinding("y")
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
		    <has classname="Request" identifier="$nil" attribute="user_id" value="$uid"/>
		    <has classname="SKU" identifier="$sku_id" attribute="quantity" value="$quantity"/>
		    <filter><![CDATA[quantity > 1]]></filter>
		</lhs>
		<rhs>Handler</rhs>
	   </production>
	   <production>
		<lhs>
		    <has classname="Request" identifier="$uid" attribute="id" value="$uid"/>
		    <has classname="SKU" identifier="$sku_id" attribute="quantity" value="$quantity"/>
		    <filter><![CDATA[quantity > 10]]></filter>
		</lhs>
		<rhs>Handler</rhs>
	   </production>
	</data>`
	n := NewNetwork()
	env := make(Env)
	env["Handler"] = func (network *Network, token *Token) {
		fmt.Println(token)
	}
	_, err := n.AddProductionFromXML(data)
	if err != nil {
		t.Error(err)
		return
	}

	wmes := []*WME{
		NewWME("Request", "nil", "user_id", "100001"),
		NewWME("SKU", "1", "quantity", "2"),
		NewWME("SKU", "2", "quantity", "20"),
	}
	for idx := range wmes {
		n.AddWME(wmes[idx])
	}
	n.ExecuteRules(env)
}

func TestFromJSON(t *testing.T) {
	data := `
	{
	  "productions": [
	    {
	      "lhs": [
	        {
	          "attribute": "quantity",
	          "classname": "ProductSKU",
	          "identifier": "1",
	          "tag": "has",
	          "value": "$quantity"
	        }
	      ],
	      "rhs": {
	        "kind": "quota:user:sku",
	        "quota": 1
	      }
	    },
	    {
	      "lhs": [
	        {
	          "attribute": "quantity",
	          "classname": "ProductSKU",
	          "identifier": "$sku_id",
	          "tag": "has",
	          "value": "$quantity"
	        }
	      ],
	      "rhs": {
	        "kind": "quota:user:sku",
	        "quota": 3
	      }
	    }
	  ]
	}`
	_, err := FromJSON(data)
	if err != nil {
		t.Error(err)
	}
}
