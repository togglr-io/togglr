package rules_test

import (
	"testing"

	"github.com/eriktate/toggle/rules"
)

func Test_IntComparable(t *testing.T) {
	cases := []struct {
		name           string
		val1           int
		val2           int
		expectedEq     bool
		expectedGt     bool
		expectedIsTrue bool
	}{
		{
			name:           "equal ints",
			val1:           42,
			val2:           42,
			expectedEq:     true,
			expectedGt:     false,
			expectedIsTrue: true,
		},
		{
			name:           "larger int",
			val1:           50,
			val2:           42,
			expectedEq:     false,
			expectedGt:     true,
			expectedIsTrue: true,
		},
		{
			name:           "smaller int",
			val1:           42,
			val2:           50,
			expectedEq:     false,
			expectedGt:     false,
			expectedIsTrue: true,
		},
		{
			name:           "zeroes",
			val1:           0,
			val2:           0,
			expectedEq:     true,
			expectedGt:     false,
			expectedIsTrue: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			val1 := rules.NewInt(c.val1)
			val2 := rules.NewInt(c.val2)

			if val1.Eq(val2) != c.expectedEq {
				t.Fatalf("expected equality check to be %t for %d == %d", c.expectedEq, val1, val2)
			}

			if val1.Gt(val2) != c.expectedGt {
				t.Fatalf("expected greater-than check to be %t for %d > %d", c.expectedGt, val1, val2)
			}

			if val1.IsTrue() != c.expectedIsTrue {
				t.Fatalf("expected truthiness check to be %t for %d", c.expectedIsTrue, val1)
			}
		})
	}
}

func Test_IntEvaluate(t *testing.T) {
	str := rules.NewInt(42)
	evaluated := str.Evaluate()
	if val, ok := evaluated.(rules.Int); ok {
		if val.Eq(rules.NewInt(42)) {
			return
		}
		t.Fatalf("evaluated int does not equal original int")
	}
	t.Fatalf("evaluated int is not a Comparable of type Int")
}
