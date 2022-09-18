package null

import (
	"encoding/json"
	"fmt"

	uuid "github.com/satori/go.uuid"
)

var empty uuid.UUID

// NewUUIDValid returns new nullable uuid with a valid value
func NewUUIDValid(id uuid.UUID) UUID {
	return newNull(id, true)
}

// NewUUIDInvalid returns new nullable uuid with invalid value.
func NewUUIDInvalid() UUID {
	return newNull(*(&empty), false)
}

func newNull(id uuid.UUID, valid bool) UUID {
	return UUID{NullUUID: uuid.NullUUID{UUID: id, Valid: valid}}
}

// UUID is the wrapper around satori uuid.NullUUID that is capable of being
// marshaled and unmarshaled into json encoding format.
type UUID struct {
	uuid.NullUUID
}

// compile time check for the json.Marshaler interface.
var _ json.Marshaler = UUID{}

// MarshalJSON implements json.Marshaler interface.
// Encodes "null" if the n.Valid is false
func (n UUID) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(n.UUID.String())
}

// compile time check for the json.Unmarshaler interface.
var _ json.Unmarshaler = &UUID{}

// UnmarshalJSON implements json.Unmarshaler interface.
// Decodes "null" into non valid UUID.
func (n *UUID) UnmarshalJSON(data []byte) (err error) {
	var v interface{}
	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}

	switch x := v.(type) {
	case string:
		n.UUID, err = uuid.FromString(x)
	case map[string]interface{}:
		err = json.Unmarshal(data, &n.NullUUID)
	case nil:
		n.Valid = false
		return nil
	default:
		err = fmt.Errorf("json: cannot unmarhsal %T into Go value of type NullUUID", v)
	}
	n.Valid = err == nil
	return err
}
