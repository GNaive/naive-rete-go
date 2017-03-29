package rete

import (
	"testing"
	"fmt"
)

func F(a string, b string) map[string]interface{} {
	result := make(map[string]interface{})
	result["a"] = a
	result["b"] = b
	return result
}

func TestEvalFromString(t *testing.T) {
	env := make(map[string]interface {})
	env["F"] = F
	result, err := EvalFromString(`F("hello", "world")`, env)
	if err != nil {
		t.Error(err)
	} else {
		fmt.Println("eval result:", result[0])
	}
}
