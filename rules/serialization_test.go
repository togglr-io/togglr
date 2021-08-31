package rules_test

import (
	"encoding/json"
	"testing"

	"github.com/togglr-io/togglr/rules"
)

func Test_UnmarshalExpression(t *testing.T) {
	cases := []struct {
		name     string
		raw      string
		expected rules.Comparable
		metadata map[string]rules.Comparable
	}{
		{
			name:     "string expr",
			raw:      `{ "type": "string", "value": "hello" }`,
			expected: rules.NewString("hello"),
		},
		{
			name:     "int expr",
			raw:      `{ "type": "int", "value": 42 }`,
			expected: rules.NewInt(42),
		},
		{
			name:     "float expr",
			raw:      `{ "type": "float", "value": 42.5 }`,
			expected: rules.NewFloat(42.5),
		},
		{
			name:     "bool expr",
			raw:      `{ "type": "bool", "value": true }`,
			expected: rules.NewBool(true),
		},
		{
			name:     "ident",
			raw:      `{ "type": "ident", "value": "name" }`,
			expected: rules.NewString("toggle"),
			metadata: map[string]rules.Comparable{
				"name": rules.NewString("toggle"),
			},
		},
		{
			name: "binary expr",
			raw: `{
				"type": "binary",
				"op": "==",
				"left": {
					"type": "string",
					"value": "hello"
				},
				"right": {
					"type": "string",
					"value": "hello"
				}
			}`,
			expected: rules.NewBool(true),
		},
		{
			name: "nested binary expr",
			raw: `{
				"type": "binary",
				"op": "&&",
				"left": {
					"type": "binary",
					"op": "==",
					"left": {
						"type": "int",
						"value": 5
					},
					"right": {
						"type": "int",
						"value": 5
					}
				},
				"right": {
					"type": "bool",
					"value": true
				}
			}`,
			expected: rules.NewBool(true),
		},
		{
			name: "nested failing binary expr",
			raw: `{
				"type": "binary",
				"op": "&&",
				"left": {
					"type": "binary",
					"op": "==",
					"left": {
						"type": "int",
						"value": 5
					},
					"right": {
						"type": "int",
						"value": 3
					}
				},
				"right": {
					"type": "bool",
					"value": true
				}
			}`,
			expected: rules.NewBool(false),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var expr rules.Expression
			if err := json.Unmarshal([]byte(c.raw), &expr); err != nil {
				t.Fatalf("failed to unmarshal: %s", err)
			}

			// catch ident types and insert metadata
			if expr.Type == rules.ExprTypeIdent {
				expr.Ident.SetMetadata(c.metadata)
			}

			if !expr.Evaluate().Eq(c.expected) {
				t.Fatalf("expected expression to evaluate to %+v", c.expected)
			}
		})
	}
}

func Test_MarshalExpression(t *testing.T) {
	cases := []struct {
		name       string
		expression rules.Expression
		expected   string
	}{
		{
			name: "marshal binary",
			expression: rules.Expression{
				Type: rules.ExprTypeBinary,
				Binary: rules.NewBinary(
					rules.NewString("hello"),
					rules.NewString("hello"),
					rules.BinOpEq,
				),
			},
			expected: `{"type":"binary","left":{"type":"string","value":"hello"},"right":{"type":"string","value":"hello"},"op":"=="}`,
		},
		{
			name: "marshal binary with ident",
			expression: rules.Expression{
				Type: rules.ExprTypeBinary,
				Binary: rules.NewBinary(
					rules.NewIdent("hello", nil),
					rules.NewString("hello"),
					rules.BinOpEq,
				),
			},
			expected: `{"type":"binary","left":{"type":"ident","value":"hello"},"right":{"type":"string","value":"hello"},"op":"=="}`,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			data, err := json.Marshal(c.expression)
			if err != nil {
				t.Fatalf("failed to marshal expression: %s", err)
			}

			if string(data) != c.expected {
				t.Log(string(data))
				t.Fatalf("marshaled expression does not match expected")
			}
		})
	}
}
