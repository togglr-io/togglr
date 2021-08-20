package rules

// A Float expression represents a float literal during rule evaluation
type Float struct {
	Value float32 `json:"value"`
}

// NewFloat returns a new Float expression
func NewFloat(val float32) Float {
	return Float{val}
}

// Eq checks if the other Comparable is a Float with the same value
func (f Float) Eq(other Comparable) bool {
	if val, ok := other.(Float); ok {
		// TODO (etate): Might have to do an approximation here
		return f.Value == val.Value
	}

	return false
}

// Gt checks if the other Comparable is a Float that is smaller
func (f Float) Gt(other Comparable) bool {
	if val, ok := other.(Float); ok {
		// TODO (etate): Might have to do an approximation here
		return f.Value > val.Value
	}

	return false
}

// IsTrue is a truthiness check that treats any postive float as true and any 0 or negative
// float as false
func (f Float) IsTrue() bool {
	return f.Value > 0.0
}

// Evaluate returns the Float expression as a Comparable
func (f Float) Evaluate() Comparable {
	return f
}
