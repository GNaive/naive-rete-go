package rete

import (
	"go/token"
	"go/ast"
	"strconv"
)


func Eval(exp ast.Expr) interface{} {
	switch exp := exp.(type) {
	case *ast.BinaryExpr:
		return EvalBinaryExpr(exp)
	case *ast.BasicLit:
		switch exp.Kind {
		case token.INT, token.FLOAT:
			i, _ := strconv.ParseFloat(exp.Value, 32)
			return i
		}
	case *ast.UnaryExpr:
		switch exp.Op {
		case token.ADD:
			return Eval(exp.X).(float64)
		case token.SUB:
			return -Eval(exp.X).(float64)
		}
	case *ast.ParenExpr:
		return Eval(exp.X)
	case *ast.Ident:
		return EvalIdent(exp.Name)
	}

	return nil
}

func EvalIdent(exp string) interface{} {
	switch exp {
	case "true":
		return true
	case "false":
		return false
	}
	return nil
}

func EvalBinaryExpr(exp *ast.BinaryExpr) interface{} {
	left := Eval(exp.X)
	right := Eval(exp.Y)
	switch exp.Op {
	case token.GTR:
		return left.(float64) > right.(float64)
	case token.LSS:
		return left.(float64) < right.(float64)
	case token.GEQ:
		return left.(float64) >= right.(float64)
	case token.LEQ:
		return left.(float64) <= right.(float64)
	case token.ADD:
		return left.(float64) + right.(float64)
	case token.SUB:
		return left.(float64) - right.(float64)
	case token.MUL:
		return left.(float64) * right.(float64)
	case token.QUO:
		return left.(float64) / right.(float64)
	case token.REM:
		return int(left.(float64)) % int(right.(float64))
	case token.LAND:
		return left.(bool) && right.(bool)
	}
	return nil
}
