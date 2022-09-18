package user

import (
	"fmt"
	"strings"
	"testing"
	"time"
	"github.com/stretchr/testify/require"
	

	"github.com/investapp/backend/models/user/contact"
	"github.com/investapp/backend/pkg/crypto"
	"github.com/investapp/backend/pkg/errdef"
	"github.com/investapp/backend/pkg/null"
	"github.com/investapp/backend/pkg/ptrto"
	"github.com/investapp/backend/pkg/random"
	"github.com/investapp/backend/pkg/valid"
)

// ProcessName is the constant used to store the errdef key value.
const ProcessName = "user"

// User is a database model for the user basic information as well as password and role.
type User struct {
	ID              uint             `json:"id" sql:",pk"`
	CreatedAt       time.Time        `json:"created_at" sql:",notnull"`
	UpdatedAt       time.Time        `json:"updated_at" sql:",notnull"`
	Dob             *string          `json:"dob" sql:"-"`
	LastSignedAt    *time.Time       `json:"last_signed_at,omitempty"`
	Firstname       string           `json:"firstname,omitempty" sql:",notnull"`
	Lastname        string           `json:"lastname,omitempty" sql:",notnull"`
	Username        string           `json:"username,omitempty" sql:",notnull"`
	Hash            *string          `json:"-" sql:"password,notnull"`
	CreatorID       *uint            `json:"creator_id,omitempty"`
	Contacts        contact.Contacts `json:"contacts,omitempty" sql:"-"`
	PicturePath     *string          `json:"picture_path,omitempty"`
	CryptoAddressID *uint            `json:"crypto_address_id" sql:",notnull"`
	// 2FA definitions
	TwoFactorAuthID       *uint     `json:"2fa_id" sql:"2fa_id"`
	TwoFactorAuthVerifyID null.UUID `json:"-" sql:"2fa_verify_id"`
}

// HasPwd will tell you if user set password
func (u User) HasPwd() bool {
	return u.Hash != nil
}

// ComparePwd will compare existing password with given
// and return bool
func (u User) ComparePwd(pwd string) bool {
	if !u.HasPwd() {
		return false
	}
	return crypto.CompareCrypts([]byte(*u.Hash), []byte(pwd))
}

// GetEmailContact gets the email contact for user.
// If not email is found - (SHOULD NOT OCCUR) - the contact is nil.
func (u *User) GetEmailContact() *contact.Contact {
	email, ok := u.Contacts.GetByType(contact.Email)
	if ok {
		return &email
	}
	return nil
}

// SetPwd will set pwd and hash
func (u *User) SetPwd(pwd string) *errdef.Error {
	if len(pwd) < 8 {
		return errdef.ErrInvalidArgumentf("password", "too weak")
	}
	hash, err := crypto.Crypt([]byte(pwd))
	if err != nil {
		return errdef.ErrInvalidArgumentf(err.Error(), "failed to encrypt pwd")
	}
	hashStr := string(hash)
	u.Hash = &hashStr
	return nil
}

// Sanitize will sanitize existing user
func (u *User) Sanitize() {
	u.sanUsername()
	u.sanFirstname()
	u.sanLastname()
}

// Name will return full name of user
func (u User) Name() string {
	return strings.Trim(fmt.Sprintf("%s %s", u.Firstname, u.Lastname), " ")
}

func Now() time.Time {
	// todo we should enfore UTC
	return time.Now().UTC().Truncate(time.Second)
}



// Validate validates struct content.
func (u User) Validate() *errdef.Error {
	errSet := errdef.ErrInvalidArgumentf(ProcessName, "validate")
	switch {
	case len(u.Firstname) > 200:
		msg := fmt.Sprintf("out of range 1-200 characters: '%s'", u.Firstname)
		errSet.Detail= msg
		return errSet
	case len(u.Lastname) > 200:
		msg := fmt.Sprintf("out of range 1-200 characters: '%s'", u.Lastname)
		errSet.Detail = msg
		return errSet
	case !valid.Username(u.Username):
		errSet.Detail=  "invalid username"
		return errSet
	case u.CreatorID != nil && u.ID == *u.CreatorID:
		errSet.Detail= "creator_id - is self referencing"
		return errSet
	}
	return nil
}

func (u *User) sanUsername() {
	u.Username = strings.ToLower(strings.Trim(u.Username, " "))
}

func (u *User) sanFirstname() {
	firstname := strings.Trim(u.Firstname, " ")
	firstname = strings.Title(firstname)
	u.Firstname = firstname
}

func (u *User) sanLastname() {
	lastname := strings.Trim(u.Lastname, " ")
	lastname = strings.Title(lastname)
	u.Lastname = lastname
}

// Users is list of users
type Users []User

// Validate validates struct content.
func (uu Users) Validate() *errdef.Error {
	for _, u := range uu {
		errSet := u.Validate()
		if errSet != nil {
			return errSet
		}
	}
	return nil
}

// GetByID get user by id
func (uu Users) GetByID(id uint) (User, bool) {
	for _, u := range uu {
		if u.ID == id {
			return u, true
		}
	}
	return User{}, false
}

// IDs return list of ids
func (uu Users) IDs() (results []uint) {
	for _, u := range uu {
		results = append(results, u.ID)
	}
	return
}

var minDate = time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)

// DOB parses provided date of birth string into time.Time.
// returns error if the format or value is not valid.
func DOB(s string) (*time.Time, *errdef.Error) {
	d, err := time.Parse("02/01/2006", s)
	if err != nil {
		return nil, errdef.ErrInternalf(err.Error(), "date_of_birth")
	}
	if d.IsZero() {
		return nil, errdef.ErrInternalf("date_of_birth", "zero value provided")
	}
	if d.Before(minDate) {
		return nil, errdef.ErrInternalf("date_of_birth", "value is not valid")
	}
	return &d, nil
}

// TstGenRandom will generate random user
func TstGenRandom(t testing.TB) User {
	u := User{}
	u.CreatedAt = Now()
	u.UpdatedAt = Now()
	firstname := random.Firstname(0)
	lastname := random.Lastname()
	u.Firstname = firstname
	u.Lastname = lastname
	u.Username = random.String(10)
	u.Contacts = contact.Contacts{}
	require.NoError(t, u.SetPwd("coinfinity2019"))
	if err := u.Validate(); err != nil {
		t.Fatalf("TstGenRandom: validation error: %s", err)
	}
	return u
}


// TstGenRandomFast will generate random user
func TstGenRandomFast() User {
	u := User{}
	u.CreatedAt = Now()
	u.UpdatedAt = Now()
	firstname := random.Firstname(0)
	lastname := random.Lastname()
	u.Firstname = firstname
	u.Lastname = lastname
	u.Username = random.String(10)
	u.Contacts = contact.Contacts{}

	u.Hash = ptrto.String("some")

	return u
}
