package null

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

// String is the null string value for the
type String struct {
	sql.NullString
}

// compile time check for the json.Marshaler interface.
var _ json.Marshaler = String{}

// MarshalJSON implements json.Marshaler interface.
// Encodes "null" if the n.Valid is false
func (n String) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(n.NullString.String)
}

// compile time check for the json.Unmarshaler interface.
var _ json.Unmarshaler = &String{}

// UnmarshalJSON implements json.Unmarshaler interface.
func (n *String) UnmarshalJSON(data []byte) (err error) {
	var v interface{}
	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}

	switch x := v.(type) {
	case string:
		n.String = x
	case map[string]interface{}:
		err = json.Unmarshal(data, &n.NullString)
	case nil:
		n.Valid = false
		return nil
	default:
		err = fmt.Errorf("json: cannot unmarshal #{v} into Go value of type null.String")
	}
	n.Valid = err == nil
	return err
}
