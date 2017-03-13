package rete

import (
	"go/parser"
	"testing"
)

func TestEval(t *testing.T) {
	exp, _ := parser.ParseExpr("1 > 2")
	// ast.Print(token.NewFileSet(), exp)
	// fmt.Println(Eval(exp))
	if Eval(exp).(bool) {
		t.Error("error")
	}
}
