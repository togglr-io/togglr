package rules

// A UnaryOp represents all possible unary operations.
type UnaryOp string

// All available UnaryOps
const (
	UnaryOpNot   = UnaryOp("!")
	UnaryOpExist = UnaryOp("!!")
)

type Unary struct {
	expr Expr
	op   UnaryOp
}

func (u Unary) Evaluate(md Metadata) Comparable {
	val := u.expr.Evaluate(md)
	switch u.op {
	case UnaryOpNot:
		return NewBool(!val.IsTrue())
	case UnaryOpExist:
		// TODO (etate): figure out what this means
		return NewBool(true)
	}

	return NewBool(val.IsTrue())
}
