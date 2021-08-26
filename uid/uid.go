package uid

import (
	"encoding/json"

	"github.com/google/uuid"
)

// A UID represents a unique identifier. Other packages can leverage this type
// without caring about the underlying implementation driving it
type UID struct {
	uuid.NullUUID
}

// New creates a new UID
func New() UID {
	return UID{
		NullUUID: uuid.NullUUID{UUID: uuid.New(), Valid: true},
	}
}

// IsNull returns whether or not a given UID is actually null
func (u UID) IsNull() bool {
	return !u.Valid
}

// FromString attempts to parse a UID from a string
func FromString(id string) (UID, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return UID{}, err
	}

	return UID{
		NullUUID: uuid.NullUUID{UUID: uid, Valid: true},
	}, nil
}

// String returns the string representation of a UID
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

// UnmarshalJSON implements the json.Unmarshaller interface for UID
func (u *UID) UnmarshalJSON(data []byte) error {
	// treat empty strings as null
	if len(data) == 0 || string(data) == "" || string(data) == "\"\"" || string(data) == "\"null\"" {
		u.NullUUID = uuid.NullUUID{Valid: false}
		return nil
	}

	return json.Unmarshal(data, &u.NullUUID)
}
