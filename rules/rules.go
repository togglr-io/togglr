package rules

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

// Metadata is included with the initial request for Toggles when a new client initializes. It's used to
// evaluate rules and determine the final value for each toggle.
type Metadata map[string]Comparable

// Enums

// An ExprType represents all of the types of expressions possible.
type ExprType string

// All available ExprTypes
const (
	ExprTypeBinary = ExprType("binary")
	ExprTypeUnary  = ExprType("unary")
	ExprTypeString = ExprType("string")
	ExprTypeIdent  = ExprType("ident")
	ExprTypeInt    = ExprType("int")
	ExprTypeFloat  = ExprType("float")
	ExprTypeBool   = ExprType("bool")
	ExprTypeNoop   = ExprType("noop")
)

// A Comparable can be compared with another Comparable to evaluate to a bool.
type Comparable interface {
	Eq(other Comparable) bool
	Gt(other Comparable) bool
	IsTrue() bool
}

// An Expr is anything that can Evaluate to a Comparable.
type Expr interface {
	Evaluate(md Metadata) Comparable
}

// An Expression is the physical (i.e. serializable) representation for everything implementing the Expr interface. It's essentially a discriminated union.
type Expression struct {
	Binary Binary
	Unary  Unary
	Ident  Ident
	String String
	Int    Int
	Float  Float
	Bool   Bool
	Type   ExprType `json:"type"`
}

func ExpressionFromExpr(expr Expr) Expression {
	switch v := expr.(type) {
	case Binary:
		return Expression{Binary: v, Type: ExprTypeBinary}
	case Unary:
		return Expression{Unary: v, Type: ExprTypeUnary}
	case Ident:
		return Expression{Ident: v, Type: ExprTypeIdent}
	case String:
		return Expression{String: v, Type: ExprTypeString}
	case Int:
		return Expression{Int: v, Type: ExprTypeInt}
	case Float:
		return Expression{Float: v, Type: ExprTypeFloat}
	case Bool:
		return Expression{Bool: v, Type: ExprTypeBool}
	case Expression:
		return v // if we find an Expression, just return it as is
	}

	return Expression{Type: ExprTypeNoop}
}

func (e Expression) Evaluate(md Metadata) Comparable {
	switch e.Type {
	case ExprTypeBinary:
		return e.Binary.Evaluate(md)
	case ExprTypeUnary:
		return e.Unary.Evaluate(md)
	case ExprTypeIdent:
		return e.Ident.Evaluate(md)
	case ExprTypeString:
		return e.String.Evaluate(md)
	case ExprTypeInt:
		return e.Int.Evaluate(md)
	case ExprTypeFloat:
		return e.Float.Evaluate(md)
	case ExprTypeBool:
		return e.Bool.Evaluate(md)
	}

	// TODO (etate): consider adding errors or a noop for cases like these
	return NewBool(false)
}

func (e Expression) MarshalJSON() ([]byte, error) {
	switch e.Type {
	case ExprTypeBinary:
		return json.Marshal(e.Binary)
	case ExprTypeUnary:
		return json.Marshal(e.Unary)
	case ExprTypeIdent:
		return json.Marshal(e.Ident)
	case ExprTypeString:
		return json.Marshal(e.String)
	case ExprTypeInt:
		return json.Marshal(e.Int)
	case ExprTypeFloat:
		return json.Marshal(e.Float)
	case ExprTypeBool:
		return json.Marshal(e.Bool)
	}

	return nil, fmt.Errorf("failed to marshal invalid Expression type %s", e.Type)
}

func (e *Expression) UnmarshalJSON(data []byte) error {
	// need to peel off the type to figure out how to unmarshal
	exprType := struct {
		Type ExprType `json:"type"`
	}{}

	if err := json.Unmarshal(data, &exprType); err != nil {
		return err
	}

	e.Type = exprType.Type

	switch e.Type {
	case ExprTypeBinary:
		return json.Unmarshal(data, &e.Binary)
	case ExprTypeUnary:
		return json.Unmarshal(data, &e.Unary)
	case ExprTypeIdent:
		return json.Unmarshal(data, &e.Ident)
	case ExprTypeString:
		return json.Unmarshal(data, &e.String)
	case ExprTypeInt:
		return json.Unmarshal(data, &e.Int)
	case ExprTypeFloat:
		return json.Unmarshal(data, &e.Float)
	case ExprTypeBool:
		return json.Unmarshal(data, &e.Bool)
	}

	return fmt.Errorf("failed to unmarshal invalid Expression type %s", e.Type)
}

// A Rule evaluates against Metadata to determine a value for a particular Toggle.
type Rule struct {
	Op   BinOp      `json:"op"`
	Expr Expression `json:"expression"`
}

func (r Rule) Validate() error {
	if r.Op != BinOpAnd && r.Op != BinOpOr {
		return errors.New("A Rule Op can only be logical (&& ||)")
	}

	return nil
}

// Rules is an alias to a Rule slice that we can implement some interfaces on
type Rules []Rule

// Value impelments the sql.Valuer interface
func (r Rules) Value() (driver.Value, error) {
	data, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// Scan impelements the sql.Scanner interface
func (r *Rules) Scan(src interface{}) error {
	var source []byte
	switch val := src.(type) {
	case string:
		source = []byte(val)
	case []byte:
		source = val
	case nil:
		source = nil
	default:
		return errors.New("incompatible type for Rule")
	}

	if err := json.Unmarshal(source, r); err != nil {
		return fmt.Errorf("failed to unmarshal database JSON into Rules: %w", err)
	}

	return nil
}

func EvaluateRules(md Metadata, rules ...Rule) bool {
	var prevExpr Expr
	prevExpr = NewBool(true)
	for _, rule := range rules {
		prevExpr = NewBinary(prevExpr, rule.Expr, rule.Op)
	}

	return prevExpr.Evaluate(md).IsTrue()
}

// MetaFromRaw does a typeswitch on each metadata value in order to create a map of Comparables.
func MetaFromRaw(raw map[string]interface{}) Metadata {
	md := make(Metadata, len(raw))
	for key, value := range raw {
		switch v := value.(type) {
		case string:
			md[key] = NewString(v)
		case int:
			md[key] = NewInt(v)
		case float32:
			md[key] = NewFloat(v)
		case bool:
			md[key] = NewBool(v)
		}
	}

	return md
}
