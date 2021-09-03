package rules

// An Ident represents some identifier that exists within the metadata object
type Ident struct {
	Type  ExprType `json:"type"`
	Value string   `json:"value"`
}

// NewIdent returns a new Ident expression with a given key and metadata map
func NewIdent(key string) Ident {
	return Ident{ExprTypeIdent, key}
}

// Evaluate returns the Ident expression as a Comparable.
func (i Ident) Evaluate(md Metadata) Comparable {
	if val, ok := md[i.Value]; ok {
		return val
	}

	// TODO (etate): What should we resolve if the key doesn't exist?
	return NewBool(false)
}
