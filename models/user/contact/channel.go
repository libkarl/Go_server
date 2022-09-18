package contact

import (
	"database/sql"
	"database/sql/driver"
	"encoding"
	"strings"

	"github.com/investapp/backend/pkg/errdef"
)

// Channel representing type of contact
type Channel int

const (
	// Email is email type
	Email Channel = iota
	// Phone is phone number type
	Phone
)

// String implements fmt.Stringer.
func (c Channel) String() string {
	switch c {
	case Email:
		return "email"
	case Phone:
		return "phone"
	default:
		return "unknown"
	}
}

// compile time check for the driver.Valuer interface.
var _ driver.Valuer = Phone

// Value implements driver.Valuer interface.
func (c Channel) Value() (driver.Value, error) {
	return driver.Value(strings.ToUpper(c.String())), nil
}

// compile time check for the sql.Scanner interface.
var (
	tmpc             = Phone
	_    sql.Scanner = &tmpc
)

// Scan implements sql.Scanner interface.
func (c *Channel) Scan(value interface{}) error {
	bin, ok := value.([]byte)
	if !ok {
		return errdef.ErrInternalf("contact_channel_scan", "invalid type")
	}
	switch string(bin) {
	case "PHONE":
		*c = Phone
	case "EMAIL":
		*c = Email
	default:
		return errdef.ErrInternalf("contact_channel", "unknown channel value: "+string(bin))
	}
	return nil
}

// compile time check for the encoding.TextMarshaler interface.
var _ encoding.TextMarshaler = Email

// MarshalText implements encoding.TextMarshaler interface.
func (c Channel) MarshalText() ([]byte, error) {
	return []byte(c.String()), nil
}

// compile time check for the encoding.TextUnmarshaler interface.
var _ encoding.TextUnmarshaler = &tmpc

// UnmarshalText implement encoding.TextUnmarshaler interface.
func (c *Channel) UnmarshalText(text []byte) error {
	switch string(text) {
	case "phone":
		*c = Phone
	case "email":
		*c = Email
	default:
		return errdef.ErrInvalidArgument("invalid channel value: " + string(text))
	}
	return nil
}
