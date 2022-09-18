package contact

import (
	"testing"

	"github.com/investapp/backend/pkg/errdef"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSanitize tests the sanitize functions of contact
func TestSanitize(t *testing.T) {
	testCases := []struct {
		contact  Contact
		expected Contact
	}{
		{
			contact:  Contact{},
			expected: Contact{},
		},
		{
			contact: Contact{
				Channel: Phone,
				Contact: "",
			},
			expected: Contact{
				Channel: Phone,
			},
		},
		{
			contact: Contact{
				Channel: Phone,
				Contact: " ",
			},
			expected: Contact{
				Channel: Phone,
			},
		},
		{
			contact: Contact{
				Channel: Phone,
				Contact: "2025550158",
			},
			expected: Contact{
				Channel: Phone,
				Contact: "+12025550158",
			},
		},
		{
			contact: Contact{
				Channel: Phone,
				Contact: "202-555-0158",
			},
			expected: Contact{
				Channel: Phone,
				Contact: "+12025550158",
			},
		},
		{
			contact: Contact{
				Channel: Email,
				Contact: "test ",
			},
			expected: Contact{
				Channel: Email,
				Contact: "test",
			},
		},
		{
			contact: Contact{
				Channel: Email,
				Contact: "tesT ",
			},
			expected: Contact{
				Channel: Email,
				Contact: "test",
			},
		},
	}
	for _, tC := range testCases {
		h := &tC.contact
		h.Sanitize()
		assert.Equal(t, tC.expected, *h)
	}
}

func TestValidate(t *testing.T) {
	testCases := []struct {
		haserr  bool
		contact Contact
	}{
		{
			contact: Contact{},
			haserr:  true,
		},
		{
			contact: Contact{
				Channel: Phone,
				Contact: "202-555-0158",
			},
			haserr: false,
		},
		{
			contact: Contact{
				Channel: Email,
				Contact: "richard@hutta.com",
			},
			haserr: false,
		},
		{
			contact: Contact{
				Channel: Email,
				Contact: "richard@hutta",
			},
			haserr: true,
		},
		{
			contact: Contact{
				Channel: Phone,
				Contact: "+1",
			},
			haserr: true,
		},
		{
			contact: Contact{
				Channel: Phone,
				Contact: "",
			},
			haserr: true,
		},
		{
			contact: Contact{
				Channel: Email,
				Contact: "",
			},
			haserr: true,
		},
	}
	for _, tC := range testCases {
		err := tC.contact.Validate()
		assert.Equal(t, err != nil, tC.haserr)
	}
}

func TestContacts(t *testing.T) {
	p := Contact{Channel: Phone, Contact: "202-555-0158"}
	e := Contact{Channel: Email, Contact: "richard@hutta.com"}
	c := Contacts{p, e}
	x, ok := c.GetByType(Phone)
	assert.Equal(t, x, p)
	assert.True(t, ok)
	x, ok = c.GetByType(Email)
	assert.Equal(t, x, e)
	assert.True(t, ok)
	c = Contacts{}
	x, ok = c.GetByType(Email)
	assert.False(t, ok)
}

func TestContactsPreferred(t *testing.T) {
	p := Contact{Channel: Phone, Contact: "202-555-0158"}
	e := Contact{Channel: Email, Contact: "richard@hutta.com"}
	c := Contacts{p, e}
	x, ok := c.Preferred(Phone)
	assert.Equal(t, x, p)
	assert.True(t, ok)
	c = Contacts{e}
	x, ok = c.Preferred(Phone)
	assert.Equal(t, x, e)
	assert.True(t, ok)
	c = Contacts{}
	x, ok = c.Preferred(Phone)
	assert.False(t, ok)
}

func TestContactsGetByID(t *testing.T) {
	p := Contact{ID: 10, Channel: Phone, Contact: "202-555-0158"}
	e := Contact{Channel: Email, Contact: "richard@hutta.com"}
	c := Contacts{p, e}
	x, ok := c.GetByID(10)
	assert.Equal(t, x, p)
	assert.True(t, ok)
}

func TestContactsSanitize(t *testing.T) {
	p := Contact{ID: 1, Channel: Phone, Contact: "+1 202-555-0158"}
	e := Contact{ID: 2, Channel: Email, Contact: "richard@hutta.com "}
	c := &Contacts{p, e}
	c.Sanitize()
	x, ok := c.GetByID(1)
	assert.True(t, ok)
	x2, ok := c.GetByID(2)
	assert.True(t, ok)

	assert.Equal(t, x.Contact, "+12025550158")
	assert.Equal(t, x2.Contact, "richard@hutta.com")
}

func TestContactsValidate(t *testing.T) {
	p := Contact{ID: 1, Channel: Phone, Contact: "+1 2"}
	e := Contact{ID: 2, Channel: Email, Contact: "richard@hutta.com "}
	c := &Contacts{p, e}
	err := c.Validate()
	require.IsType(t, errdef.ErrInvalidArgument(), err)
}

func TestRemoveDup(t *testing.T) {
	cts := Contacts{
		{Contact: "hello", ID: 1},
		{Contact: "hello", ID: 2},
		{Contact: "World", ID: 3},
	}
	cts = cts.RemoveDup()
	require.Len(t, cts, 2)
	assert.Equal(t, cts[0].Contact, "hello")
	assert.Equal(t, cts[1].Contact, "World")
}
