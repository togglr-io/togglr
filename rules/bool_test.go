package rules_test

import (
	"testing"

	"github.com/eriktate/toggle/rules"
)

func Test_BoolComparable(t *testing.T) {
	cases := []struct {
		name           string
		val1           bool
		val2           bool
		expectedEq     bool
		expectedGt     bool
		expectedIsTrue bool
	}{
		{
			name:           "both true",
			val1:           true,
			val2:           true,
			expectedEq:     true,
			expectedGt:     false,
			expectedIsTrue: true,
		},
		{
			name:           "both false",
			val1:           false,
			val2:           false,
			expectedEq:     true,
			expectedGt:     false,
			expectedIsTrue: false,
		},
		{
			name:           "true and false",
			val1:           true,
			val2:           false,
			expectedEq:     false,
			expectedGt:     false,
			expectedIsTrue: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			val1 := rules.NewBool(c.val1)
			val2 := rules.NewBool(c.val2)

			if val1.Eq(val2) != c.expectedEq {
				t.Fatalf("expected equality check to be %t for %t == %t", c.expectedEq, val1, val2)
			}

			if val1.Gt(val2) != c.expectedGt {
				t.Fatalf("expected greater-than check to be %t for %t > %t", c.expectedGt, val1, val2)
			}

			if val1.IsTrue() != c.expectedIsTrue {
				t.Fatalf("expected truthiness check to be %t for %t", c.expectedIsTrue, val1)
			}
		})
	}
}

func Test_BoolEvaluate(t *testing.T) {
	str := rules.NewBool(true)
	evaluated := str.Evaluate()
	if val, ok := evaluated.(rules.Bool); ok {
		if val.Eq(rules.NewBool(true)) {
			return
		}
		t.Fatalf("evaluated bool does not equal original bool")
	}
	t.Fatalf("evaluated bool is not a Comparable of type Bool")
}
