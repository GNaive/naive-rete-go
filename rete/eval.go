package rete

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"runtime/debug"
	"strconv"
)

func EvalFromString(s string, env Env) (result []reflect.Value, err error) {
	exp, err := parser.ParseExpr(s)
	if err != nil {
		return
	}
	return Eval(exp, env)
}

func Eval(exp ast.Expr, env Env) (result []reflect.Value, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprintf("%s %s", e, debug.Stack()))
		}
	}()
	switch exp := exp.(type) {
	case *ast.BinaryExpr:
		return EvalBinaryExpr(exp, env)
	case *ast.BasicLit:
		var r interface{}
		switch exp.Kind {
		case token.INT, token.FLOAT:
			r, err = strconv.ParseFloat(exp.Value, 64)
		case token.STRING:
			r, err = strconv.Unquote(exp.Value)
		}
		if err != nil {
			return
		}
		result = append(result, reflect.ValueOf(r))
		return result, nil
	case *ast.UnaryExpr:
		switch exp.Op {
		case token.ADD:
			return Eval(exp.X, env)
		case token.SUB:
			result, err = Eval(exp.X, env)
			if err != nil {
				return
			}
			val := -result[0].Float()
			result = []reflect.Value{reflect.ValueOf(val)}
			return
		}
	case *ast.ParenExpr:
		return Eval(exp.X, env)
	case *ast.Ident:
		return EvalIdent(exp, env)
	case *ast.CallExpr:
		return EvalCall(exp, env)
	}

	return
}

func EvalCall(exp *ast.CallExpr, env Env) (result []reflect.Value, err error) {
	f_val, err := Eval(exp.Fun, env)
	if err != nil {
		return
	}
	var args_val []reflect.Value
	for _, arg := range exp.Args {
		arg_val, err := Eval(arg, env)
		if err != nil {
			return result, err
		}
		args_val = append(args_val, arg_val[0])
	}
	result = f_val[0].Call(args_val)
	return
}

func EvalIdent(exp *ast.Ident, env Env) (result []reflect.Value, err error) {
	v := env[exp.Name]
	if v == nil {
		err = errors.New(fmt.Sprintf("Ident `%s` undefined", exp))
	} else {
		result = append(result, reflect.ValueOf(v))
	}
	return
}

func EvalBinaryExpr(exp *ast.BinaryExpr, env map[string]interface{}) (result []reflect.Value, err error) {
	leftResult, err := Eval(exp.X, env)
	if err != nil {
		return
	}
	rightResult, err := Eval(exp.Y, env)
	if err != nil {
		return
	}
	var r interface{}
	left := value2float(leftResult[0])
	right := value2float(rightResult[0])

	switch exp.Op {
	case token.GTR:
		r = left > right
	case token.LSS:
		r = left < right
	case token.GEQ:
		r = left >= right
	case token.LEQ:
		r = left <= right
	case token.EQL:
		r = left == right
	case token.ADD:
		r = left + right
	case token.SUB:
		r = left - right
	case token.MUL:
		r = left * right
	case token.QUO:
		r = left / right
	default:
		err = errors.New(fmt.Sprintf("OP `%s` undefined", exp.Op))
		return
	}
	result = append(result, reflect.ValueOf(r))
	return result, nil
}

func value2float(v reflect.Value) float64 {
	switch v.Kind() {
	case reflect.Float64:
		return v.Float()
	case reflect.String:
		r, _ := strconv.ParseFloat(v.String(), 64)
		return r
	}
	return v.Float()
}
