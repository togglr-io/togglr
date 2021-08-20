package uid

import "github.com/google/uuid"

// A UID represents a unique identifier. Other packages can leverage this type
// without caring about the underlying implementation driving it.
type UID struct {
	uuid.NullUUID
}

// New creates a new UID.
func New() UID {
	return UID{
		NullUUID: uuid.NullUUID{UUID: uuid.New(), Valid: true},
	}
}

// FromString attempts to parse a UID from a string.
func FromString(id string) (UID, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return UID{}, err
	}

	return UID{
		NullUUID: uuid.NullUUID{UUID: uid, Valid: true},
	}, nil
}

// String returns the string representation of a UID.
func (u UID) String() string {
	if u.Valid {
		return u.UUID.String()
	}

	return "null"
}

// Equals returns whether or not two UIDs match
func (u UID) Equals(uid UID) bool {
	return u.String() == uid.String()
}
