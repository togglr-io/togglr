package rules

import (
	"encoding/json"
)

// A BinOp represents all possible binary operations.
type BinOp string

// All available BinOps
const (
	BinOpEq    = BinOp("==")
	BinOpNotEq = BinOp("!=")
	BinOpGt    = BinOp(">")
	BinOpLt    = BinOp("<")
	BinOpGtEq  = BinOp(">=")
	BinOpLtEq  = BinOp("<=")
	BinOpAnd   = BinOp("&&")
	BinOpOr    = BinOp("||")
)

// A Binary expression that compares a left Expr with a right Expr using a particular operator
type Binary struct {
	left  Expr
	right Expr
	op    BinOp
}

// NewBinary returns a new Binary expression
func NewBinary(left, right Expr, op BinOp) Binary {
	return Binary{left, right, op}
}

// have to provide an alternate struct as a marshal target, otherwise there's a cycling issue between Binary and Expression
type marshalTarget struct {
	Type  ExprType   `json:"type"`
	Left  Expression `json:"left"`
	Right Expression `json:"right"`
	Op    BinOp      `json:"op"`
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (b *Binary) UnmarshalJSON(data []byte) error {
	var target marshalTarget
	if err := json.Unmarshal(data, &target); err != nil {
		return err
	}

	b.left = target.Left
	b.right = target.Right
	b.op = target.Op

	return nil
}

// MarshalJSON implements the json.Marshaler interface
func (b Binary) MarshalJSON() ([]byte, error) {
	target := marshalTarget{
		Type:  ExprTypeBinary,
		Left:  ExpressionFromExpr(b.left),
		Right: ExpressionFromExpr(b.right),
		Op:    b.op,
	}

	return json.Marshal(target)
}

// Eq checks if the evaluated result of the Binary expression is equal to the other Comparable
func (b Binary) Eq(other Comparable) bool {
	val := b.Evaluate()
	return val.Eq(other)
}

// Gt always returns false since Binary expression always resolve to a Bool
func (b Binary) Gt(other Comparable) bool {
	// since Binary exprs technically evaluate to Bools, implementing Gt doesn't really make sense
	return false
}

// IsTrue is a truthiness check that relies on the evaluate of the Binary expression
func (b Binary) IsTrue() bool {
	return b.Evaluate().IsTrue()
}

// Evaluate resolves the Binary expression to the resulting Bool expression
func (b Binary) Evaluate() Comparable {
	left := b.left.Evaluate()
	right := b.right.Evaluate()
	switch b.op {
	case BinOpEq:
		return NewBool(left.Eq(right))
	case BinOpNotEq:
		return NewBool(!left.Eq(right))
	case BinOpGt:
		return NewBool(left.Gt(right))
	case BinOpGtEq:
		return NewBool(left.Gt(right) || left.Eq(right))
	case BinOpLt:
		return NewBool(!left.Gt(right))
	case BinOpLtEq:
		return NewBool(left.Eq(right) || !left.Gt(right))
	case BinOpAnd:
		return NewBool(left.IsTrue() && right.IsTrue())
	case BinOpOr:
		return NewBool(left.IsTrue() || right.IsTrue())
	}

	// TODO (etate): Consider making Comparables error.
	return NewBool(false)
}
