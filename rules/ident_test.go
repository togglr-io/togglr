package rules_test

import (
	"testing"

	"github.com/togglr-io/togglr/rules"
)

func Test_IdentEvaluate(t *testing.T) {
	metadata := map[string]rules.Comparable{
		"hello":     rules.NewString("world"),
		"age":       rules.NewInt(29),
		"cost":      rules.NewFloat(100.50),
		"toggleOn":  rules.NewBool(true),
		"toggleOff": rules.NewBool(false),
	}

	cases := []struct {
		name   string
		ident  rules.Ident
		expect rules.Comparable
	}{
		{
			name:   "equal strings",
			ident:  rules.NewIdent("hello"),
			expect: rules.NewString("world"),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if !c.ident.Evaluate(metadata).Eq(c.expect) {
				t.Fatalf("expected ident %s to equal %s", c.ident.Value, c.expect)
			}
		})
	}
}
