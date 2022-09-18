package contacttst

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/investapp/backend/models/user/contact"
	"github.com/investapp/backend/models/user"
	"github.com/investapp/backend/pkg/random"
)

// NewRandomEmail will generate new contact with email
func NewRandomEmail(t *testing.T) contact.Contact {
	c := &contact.Contact{}
	c.CreatedAt = user.Now()
	c.UpdatedAt = user.Now()
	c.Contact = random.Email()
	c.Channel = contact.Email
	c.Sanitize()
	require.NoError(t, c.Validate())
	return *c
}

// NewRandomPhone will generate new contact with phone
func NewRandomPhone(t *testing.T) contact.Contact {
	c := &contact.Contact{}
	c.CreatedAt = user.Now()
	c.UpdatedAt = user.Now()
	c.Contact = random.Phone()
	c.Channel = contact.Phone
	c.Sanitize()
	require.NoError(t, c.Validate())
	return *c
}

// NewRandom give you random contacts phone and email
func NewRandom(t *testing.T) contact.Contacts {
	return contact.Contacts{
		NewRandomEmail(t),
		NewRandomPhone(t),
	}
}
