package rules_test

import (
	"testing"

	"github.com/eriktate/toggle/rules"
)

func Test_FloatComparable(t *testing.T) {
	cases := []struct {
		name           string
		val1           float32
		val2           float32
		expectedEq     bool
		expectedGt     bool
		expectedIsTrue bool
	}{
		{
			name:           "equal floats",
			val1:           42.0,
			val2:           42.0,
			expectedEq:     true,
			expectedGt:     false,
			expectedIsTrue: true,
		},
		{
			name:           "larger float",
			val1:           42.5,
			val2:           42.0,
			expectedEq:     false,
			expectedGt:     true,
			expectedIsTrue: true,
		},
		{
			name:           "smaller float",
			val1:           42.0,
			val2:           42.5,
			expectedEq:     false,
			expectedGt:     false,
			expectedIsTrue: true,
		},
		{
			name:           "zeroes",
			val1:           0.0,
			val2:           0.0,
			expectedEq:     true,
			expectedGt:     false,
			expectedIsTrue: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			val1 := rules.NewFloat(c.val1)
			val2 := rules.NewFloat(c.val2)

			if val1.Eq(val2) != c.expectedEq {
				t.Fatalf("expected equality check to be %t for %f == %f", c.expectedEq, val1, val2)
			}

			if val1.Gt(val2) != c.expectedGt {
				t.Fatalf("expected greater-than check to be %t for %f > %f", c.expectedGt, val1, val2)
			}

			if val1.IsTrue() != c.expectedIsTrue {
				t.Fatalf("expected truthiness check to be %t for %f", c.expectedIsTrue, val1)
			}
		})
	}
}

func Test_FloatEvaluate(t *testing.T) {
	str := rules.NewFloat(42)
	evaluated := str.Evaluate()
	if val, ok := evaluated.(rules.Float); ok {
		if val.Eq(rules.NewFloat(42)) {
			return
		}
		t.Fatalf("evaluated float does not equal original float")
	}
	t.Fatalf("evaluated float is not a Comparable of type Float")
}
