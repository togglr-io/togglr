package rules_test

import (
	"testing"

	"github.com/eriktate/toggle/rules"
)

func Test_IdentComparable(t *testing.T) {
	metadata := map[string]rules.Comparable{
		"hello":     rules.NewString("world"),
		"age":       rules.NewInt(29),
		"cost":      rules.NewFloat(100.50),
		"toggleOn":  rules.NewBool(true),
		"toggleOff": rules.NewBool(false),
	}

	cases := []struct {
		name           string
		ident          rules.Ident
		compare        rules.Comparable
		expectedEq     bool
		expectedGt     bool
		expectedIsTrue bool
	}{
		{
			name:           "equal strings",
			ident:          rules.NewIdent("hello", metadata),
			compare:        rules.NewString("world"),
			expectedEq:     true,
			expectedGt:     false,
			expectedIsTrue: true,
		},
		{
			name:           "larger int",
			ident:          rules.NewIdent("age", metadata),
			compare:        rules.NewInt(18),
			expectedEq:     false,
			expectedGt:     true,
			expectedIsTrue: true,
		},
		{
			name:           "smaller float",
			ident:          rules.NewIdent("cost", metadata),
			compare:        rules.NewFloat(101),
			expectedEq:     false,
			expectedGt:     false,
			expectedIsTrue: true,
		},
		{
			name:           "true bool",
			ident:          rules.NewIdent("toggleOn", metadata),
			compare:        rules.NewBool(true),
			expectedEq:     true,
			expectedGt:     false,
			expectedIsTrue: true,
		},
		{
			name:           "false bool",
			ident:          rules.NewIdent("toggleOff", metadata),
			compare:        rules.NewBool(false),
			expectedEq:     true,
			expectedGt:     false,
			expectedIsTrue: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.ident.Eq(c.compare) != c.expectedEq {
				t.Fatalf("expected equality check to be %t for %+v == %+v", c.expectedEq, c.ident, c.compare)
			}

			if c.ident.Gt(c.compare) != c.expectedGt {
				t.Fatalf("expected greater-than check to be %t for %+v > %+v", c.expectedGt, c.ident, c.compare)
			}

			if c.ident.IsTrue() != c.expectedIsTrue {
				t.Fatalf("expected truthiness check to be %t for %+v", c.expectedIsTrue, c.ident)
			}
		})
	}
}

func Test_IdentEvaluate(t *testing.T) {
	metadata := map[string]rules.Comparable{
		"hello": rules.NewString("world"),
	}
	ident := rules.NewIdent("hello", metadata)

	evaluated := ident.Evaluate()
	if val, ok := evaluated.(rules.String); ok {
		if val.Eq(rules.NewString("world")) {
			return
		}
		t.Fatalf("evaluated ident does not equal original string")
	}
	t.Fatalf("evaluated ident is not a Comparable of type String")
}
