package rules

// A String expression represents a string literal during rule evaluation
type String struct {
	Type  ExprType `json:"type"`
	Value string   `json:"value"`
}

// NewString returns a new String expression
func NewString(str string) String {
	return String{ExprTypeString, str}
}

// Eq checks if the other Comparable is a String with the same value
func (s String) Eq(other Comparable) bool {
	if val, ok := other.(String); ok {
		return s.Value == val.Value
	}

	return false
}

// Gt checks if the other Comparable is a String that is lexigraphically less
func (s String) Gt(other Comparable) bool {
	if val, ok := other.(String); ok {
		return s.Value > val.Value
	}

	return false
}

// IsTrue is a truthiness check that treats an empty string as false and all other
// values as true
func (s String) IsTrue() bool {
	return s.Value != ""
}

// Evaluate returns the String expression as a Comparable
func (s String) Evaluate() Comparable {
	return s
}
