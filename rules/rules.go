package rules

import (
	"encoding/json"
	"errors"

	"github.com/togglr-io/togglr/uid"
)

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
	Evaluate() Comparable
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
	}

	return Expression{Type: ExprTypeNoop}
}

func (e Expression) Evaluate() Comparable {
	switch e.Type {
	case ExprTypeBinary:
		return e.Binary.Evaluate()
	case ExprTypeUnary:
		return e.Unary.Evaluate()
	case ExprTypeIdent:
		return e.Ident.Evaluate()
	case ExprTypeString:
		return e.String.Evaluate()
	case ExprTypeInt:
		return e.Int.Evaluate()
	case ExprTypeFloat:
		return e.Float.Evaluate()
	case ExprTypeBool:
		return e.Bool.Evaluate()
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

	return nil, errors.New("failed to marshal invalid Expression type")
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

	return errors.New("failed to unmarshal invalid Expression type")
}

// A Rule evaluates against Metadata to determine a value for a particular Toggle.
type Rule struct {
	ID       uid.UID    `json:"id"`
	ToggleID uid.UID    `json:"toggleID"`
	Expr     Expression `json:"expression"`
}
