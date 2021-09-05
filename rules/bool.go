package rules

// A Bool expression represents a boolean literal during rule evaluation
type Bool struct {
	Value bool `json:"value"`
}

// NewBool returns a new Bool expression
func NewBool(val bool) Bool {
	return Bool{val}
}

// Eq checks if the other Comparable is a Bool with the same value
func (b Bool) Eq(other Comparable) bool {
	if val, ok := other.(Bool); ok {
		return b.Value == val.Value
	}

	return false
}

// Gt always returns false for Bool expressions since they don't really make sense for bools
func (b Bool) Gt(other Comparable) bool {
	return false
}

// IsTrue just returns the value of the Bool expression
func (b Bool) IsTrue() bool {
	return b.Value
}

// Evaluate returns the Bool expression as a Comparable
func (b Bool) Evaluate(md Metadata) Comparable {
	return b
}
