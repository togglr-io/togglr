package rules

import "fmt"

// An Ident represents some identifier that exists within the metadata object
type Ident struct {
	Type     ExprType `json:"type"`
	Value    string   `json:"value"`
	metadata map[string]Comparable
}

// NewIdent returns a new Ident expression with a given key and metadata map
func NewIdent(key string, metadata map[string]Comparable) Ident {
	return Ident{ExprTypeIdent, key, metadata}
}

// SetMetadata can be used to overwrite the metadata map that an Ident uses
func (i *Ident) SetMetadata(metadata map[string]Comparable) {
	i.metadata = metadata
}

// Eq checks if the other Comparable is equal to whatever Comparable Ident refers to
func (i Ident) Eq(other Comparable) bool {
	return i.Evaluate().Eq(other)
}

// Gt checks if the other Comparable is greater than whatever Comparable Ident refers to.
func (i Ident) Gt(other Comparable) bool {
	return i.Evaluate().Gt(other)
}

// IsTrue checks if the Comparable that the Ident refers to evaluates to a truthy value.
func (i Ident) IsTrue() bool {
	return i.Evaluate().IsTrue()
}

// Evaluate returns the Ident expression as a Comparable.
func (i Ident) Evaluate() Comparable {
	return i.metadata[i.Value]
}

// String representation of the Ident (for formatting)
func (i Ident) String() string {
	return fmt.Sprintf("%+v", i.Evaluate())
}
