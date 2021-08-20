package rules

// An Int expression represents an int literal during rule evaluation
type Int struct {
	Value int `json:"value"`
}

// NewInt returns a new Int expression
func NewInt(val int) Int {
	return Int{val}
}

// Eq checks if the other Comparable is an Int with the same value
func (i Int) Eq(other Comparable) bool {
	if val, ok := other.(Int); ok {
		return i.Value == val.Value
	}

	return false
}

// Gt checks if the other Comparable is an Int that is smaller
func (i Int) Gt(other Comparable) bool {
	if val, ok := other.(Int); ok {
		return i.Value > val.Value
	}

	return false
}

// IsTrue is a truthiness check that treats any postive int as true and any 0 or negative
// int as false
func (i Int) IsTrue() bool {
	return i.Value > 0
}

// Evaluate returns the Int expression as a Comparable
func (i Int) Evaluate() Comparable {
	return i
}
