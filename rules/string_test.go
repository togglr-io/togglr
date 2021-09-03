package rules_test

import (
	"testing"

	"github.com/togglr-io/togglr/rules"
)

func Test_StringComparable(t *testing.T) {
	cases := []struct {
		name           string
		val1           string
		val2           string
		expectedEq     bool
		expectedGt     bool
		expectedIsTrue bool
	}{
		{
			name:           "equal strings",
			val1:           "hello",
			val2:           "hello",
			expectedEq:     true,
			expectedGt:     false,
			expectedIsTrue: true,
		},
		{
			name:           "lexigraphically larger",
			val1:           "xyz",
			val2:           "abc",
			expectedEq:     false,
			expectedGt:     true,
			expectedIsTrue: true,
		},
		{
			name:           "lexigraphically smaller",
			val1:           "abc",
			val2:           "xyz",
			expectedEq:     false,
			expectedGt:     false,
			expectedIsTrue: true,
		},
		{
			name:           "empty strings",
			val1:           "",
			val2:           "",
			expectedEq:     true,
			expectedGt:     false,
			expectedIsTrue: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			val1 := rules.NewString(c.val1)
			val2 := rules.NewString(c.val2)

			if val1.Eq(val2) != c.expectedEq {
				t.Fatalf("expected equality check to be %t for %s == %s", c.expectedEq, val1, val2)
			}

			if val1.Gt(val2) != c.expectedGt {
				t.Fatalf("expected greater-than check to be %t for %s > %s", c.expectedGt, val1, val2)
			}

			if val1.IsTrue() != c.expectedIsTrue {
				t.Fatalf("expected truthiness check to be %t for %s", c.expectedIsTrue, val1)
			}
		})
	}
}

func Test_StringEvaluate(t *testing.T) {
	str := rules.NewString("hello, world!")
	evaluated := str.Evaluate(nil)
	if val, ok := evaluated.(rules.String); ok {
		if val.Eq(rules.NewString("hello, world!")) {
			return
		}
		t.Fatalf("evaluated string does not equal original string")
	}
	t.Fatalf("evaluated string is not a Comparable of type String")
}
