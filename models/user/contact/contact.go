package contact

import (
	"strings"
	"time"

	"github.com/huttarichard/phone"

	"github.com/investapp/backend/pkg/errdef"
	"github.com/investapp/backend/pkg/null"
	"github.com/investapp/backend/pkg/valid"
)

// ProcessName is the constant used to store the errdef key value.
const ProcessName = "user_contact"

// Contact representing contact information
// for user.
type Contact struct {
	ID                   uint      `json:"id" sql:",pk"`
	CreatedAt            time.Time `json:"created_at" sql:",notnull"`
	UpdatedAt            time.Time `json:"updated_at" sql:",notnull"`
	Channel              Channel   `json:"type" sql:",notnull"`
	Contact              string    `json:"contact" sql:",notnull"`
	Verified             bool      `json:"verified" sql:",notnull"`
	VerifyID             null.UUID `json:"verify_id"`
	UserID               uint      `json:"user_id" sql:",notnull"`
	ConfirmationRequests uint      `json:"confirmation_requests" sql:",notnull"`
}

// Sanitize will sanitize contact
func (c *Contact) Sanitize() {
	c.Contact = strings.Trim(c.Contact, " ")
	c.Contact = strings.ToLower(c.Contact)
	switch c.Channel {
	case Email:
	case Phone:
		normalized, err := phone.Normalize(c.Contact, "")
		if err != nil {
			return
		}
		c.Contact = normalized.PhoneNumber
	}
}

// Validate will validate contact.
func (c *Contact) Validate() *errdef.Error {
	errSet := errdef.New("",ProcessName, errdef.CodeInvalidArgument)
	if len(c.Contact) == 0 || len(c.Contact) > 250 {
		return errSet
	}
	switch c.Channel {
	case Email:
		if !valid.Email(c.Contact) {
			return errSet
		}
	case Phone:
		_, err := phone.Normalize(c.Contact, "")
		if err != nil {
			return errSet
		}
	}
	return nil
}

// Contacts is list of contacts
type Contacts []Contact

// Len is part of sort.Interface.
func (cc Contacts) Len() int {
	return len(cc)
}

// Swap is part of sort.Interface.
func (cc Contacts) Swap(i, j int) {
	cc[i], cc[j] = cc[j], cc[i]
}

// Less is part of sort.Interface.
func (cc Contacts) Less(i, j int) bool {
	return cc[i].Channel < cc[j].Channel
}

// GetByID will find contact by id in list
func (cc Contacts) GetByID(id uint) (contact Contact, found bool) {
	for _, c := range cc {
		if c.ID == id {
			found = true
			contact = c
			return
		}
	}
	return
}

// Sanitize will sanitize contacts
func (cc *Contacts) Sanitize() {
	for i := range *cc {
		contact := (*cc)[i]
		contact.Sanitize()
		(*cc)[i] = contact
	}
}

// Validate will validate contacts
func (cc Contacts) Validate() *errdef.Error {
	for _, c := range cc {
		err := c.Validate()
		if err != nil {
			return err
		}
	}
	return nil
}

// GetByType will search slice for specific contact
func (cc Contacts) GetByType(t Channel) (contact Contact, found bool) {
	for _, c := range cc {
		if c.Channel != t {
			continue
		}
		contact = c
		found = true
		return
	}
	return
}

// Preferred return first preferred contact
func (cc Contacts) Preferred(t Channel) (contact Contact, found bool) {
	if c, ok := cc.GetByType(t); ok {
		contact = c
		found = true
		return
	}
	if len(cc) > 0 {
		contact = cc[0]
		found = true
		return
	}
	return
}

// RemoveDup will remove duplicates
func (cc Contacts) RemoveDup() Contacts {
	keys := make(map[string]bool)
	list := Contacts{}
	for _, entry := range cc {
		if _, value := keys[entry.Contact]; !value {
			keys[entry.Contact] = true
			list = append(list, entry)
		}
	}
	return list
}
