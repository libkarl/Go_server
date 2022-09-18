package contactdb

import (
	"context"
	"fmt"
	"testing"

	"github.com/kr/pretty"
	"github.com/satori/go.uuid"

	"github.com/investapp/backend/models/user/contact/contacttst"
	"github.com/stretchr/testify/require"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/investapp/backend/models/user/contact"
	"github.com/investapp/backend/pkg/db"
	"github.com/investapp/backend/pkg/errdef"
)

const processName = contact.ProcessName

// Contact ...
type Contact struct {
	//lint:ignore U1000 tableName is used by pg library to find model specific table
	tableName struct{} `sql:"user_contact"`
	contact.Contact
}

// BeforeInsert ...
func (c *Contact) BeforeInsert(context.Context, orm.DB) error {
	if c.CreatedAt.IsZero() {
		c.CreatedAt = db.Now()
	}
	if c.UpdatedAt.IsZero() {
		c.UpdatedAt = db.Now()
	}
	c.ID = 0
	return nil
}

// BeforeUpdate ...
func (c *Contact) BeforeUpdate(context.Context, orm.DB) error {
	c.UpdatedAt = db.Now()
	return nil
}

// Create will create contact
func Create(ctx context.Context, conn orm.DB, contact *contact.Contact) *errdef.Error {
	const operation = "failed to create contact"
	if err := db.NotNil(contact, operation); err != nil {
		return err
	}
	contact.Sanitize()
	if err := contact.Validate(); err != nil {
		return err
	}
	if err := db.CtxCheck(ctx, processName); err != nil {
		return err
	}

	model := Contact{Contact: *contact}
	_, err := conn.ModelContext(ctx, &model).Insert()
	if x, ok := err.(pg.Error); ok && x.IntegrityViolation() {
		msg := fmt.Sprintf("contact %s %s already exists", contact.Channel, contact.Contact)
		return errdef.Wrap(err, errdef.CodeAlreadyExists, msg)
	}
	if err != nil {
		return db.Wrap(err, processName)
	}
	*contact = model.Contact
	return nil
}

// DeleteByID will delete contact by user ID
func DeleteByID(ctx context.Context, conn orm.DB, id uint) *errdef.Error {
	const operation = "failed to delete contacts by user id"
	if err := db.CtxCheck(ctx, processName); err != nil {
		return err
	}
	res, err := conn.ModelContext(ctx, &Contact{}).
		Where("id = ?", id).
		Delete()
	if err != nil {
		return db.Wrap(err, operation)
	}
	if res.RowsAffected() != 1 {
		return errdef.ErrNotFound(processName, operation)
	}
	return nil
}

// EmailExists checks if given email exists
func EmailExists(ctx context.Context, conn orm.DB, email string) (bool, *errdef.Error) {
	if err := db.CtxCheck(ctx, processName); err != nil {
		return false, err
	}
	exists, err := conn.ModelContext(ctx, &Contact{}).
		Where("?TableAlias.contact = ?", email).
		Where("?TableAlias.channel = ?", contact.Email).
		Exists()
	if err != nil {
		return false, db.Wrap(err, processName)
	}
	return exists, nil
}

// GetByID will return contact by ID
func GetByID(ctx context.Context, conn orm.DB, id uint) (contact.Contact, *errdef.Error) {
	const operation = "failed to get contact"
	if err := db.CtxCheck(ctx, processName); err != nil {
		return contact.Contact{}, err
	}
	model := Contact{Contact: contact.Contact{ID: id}}
	err := conn.ModelContext(ctx, &model).Where("id = ?", id).Select()
	if err == pg.ErrNoRows {
		return model.Contact, errdef.ErrNotFound(processName, operation)
	}
	return model.Contact, db.Wrap(err, operation)
}

// GetExisting will populate existing
func GetExisting(ctx context.Context, conn orm.DB, contact *contact.Contact) *errdef.Error {
	const operation = "failed to get existing contact"
	if err := db.CtxCheck(ctx, processName); err != nil {
		
		return err
	}
	
	model := Contact{Contact: *contact}
	err := conn.ModelContext(ctx, &model).
		Where("contact = ?contact").
		Where("channel = ?channel").
		First()
	if err == pg.ErrNoRows {
		return errdef.ErrNotFound(processName, operation)
	}
	pretty.Println(err)
	if err != nil {
		return db.Wrap(err, processName)
	}
	*contact = model.Contact
	return nil
}

// GetByVerifyID gets the contact by it's related verify ID.
func GetByVerifyID(ctx context.Context, conn orm.DB, verifyID uuid.UUID) (contact.Contact, *errdef.Error) {
	if err := db.CtxCheck(ctx, processName); err != nil {
		return contact.Contact{}, err
	}
	model := Contact{}
	err := conn.ModelContext(ctx, &model).
		Where("?TableAlias.verify_id = ?", verifyID).
		First()
	if err == pg.ErrNoRows {
		return contact.Contact{}, errdef.ErrNotFound(processName, "getting contact by verifyID")
	}
	if err != nil {
		return contact.Contact{}, db.Wrap(err, processName)
	}
	return model.Contact, nil
}

// FindByUserID will return all constacts with specific UserID
func FindByUserID(ctx context.Context, conn orm.DB, id uint) (contact.Contacts, *errdef.Error) {
	var cc []contact.Contact
	if err := db.CtxCheck(ctx, processName); err != nil {
		return cc, err
	}
	models := []Contact{}
	err := conn.ModelContext(ctx, &models).
		Where("user_id = ?", id).
		Select()
	for _, c := range models {
		cc = append(cc, c.Contact)
	}
	return cc, db.Wrap(err, processName)
}

// BatchCreate will create contact
func BatchCreate(ctx context.Context, conn orm.DB, contacts *contact.Contacts) *errdef.Error {
	operation := "failed to create contacts"
	if err := db.NotNil(contacts, operation); err != nil {
		return err
	}
	contacts.Sanitize()
	errSet := contacts.Validate()
	if errSet != nil {
		return errSet
	}
	if err := db.CtxCheck(ctx, processName); err != nil {
		return err
	}
	var cc []Contact
	for _, c := range *contacts {
		cc = append(cc, Contact{Contact: c})
	}
	_, err := conn.ModelContext(ctx, &cc).Insert()
	if x, ok := err.(pg.Error); ok && x.IntegrityViolation() {
		return errdef.ErrInvalidArgument("contacts", "contact already exists")
	}
	if err != nil {
		return db.Wrap(err, processName)
	}
	var models contact.Contacts
	for _, c := range cc {
		models = append(models, c.Contact)
	}
	*contacts = models
	return nil
}

// Update will update user
func Update(ctx context.Context, conn orm.DB, contact *contact.Contact) (int, *errdef.Error) {
	const operation = "failed to update contact"
	if err := db.NotNil(contact, operation); err != nil {
		return 0, err
	}
	contact.Sanitize()
	if err := contact.Validate(); err != nil {
		return 0, err
	}
	if err := db.CtxCheck(ctx, processName); err != nil {
		return 0, err
	}
	model := Contact{Contact: *contact}
	res, err := conn.ModelContext(ctx, &model).
		Set("updated_at = ?updated_at").
		Set("verified = ?verified").
		Set("verify_id = ?verify_id").
		Set("contact = ?contact").
		Set("confirmation_requests = ?confirmation_requests").
		Where("id = ?id").
		Update()
	if err != nil {
		return 0, db.Wrap(err, processName)
	}
	*contact = model.Contact
	return res.RowsAffected(), nil
}

// TstCreate will create record and check for errors.
// If there are any it will stop test execution with t.Fail(...).
func TstCreate(t *testing.T, conn orm.DB, model *contact.Contact) {
	err := Create(context.Background(), conn, model)
	require.NoError(t, err)
}

// TstCreateRandomEmail will create random record in database.
func TstCreateRandomEmail(t *testing.T, conn orm.DB, userID uint) contact.Contact {
	model := contacttst.NewRandomEmail(t)
	model.UserID = userID
	TstCreate(t, conn, &model)
	return model
}

// TstCreateRandomPhone creates random phone contact.
func TstCreateRandomPhone(t *testing.T, conn orm.DB, userID uint) contact.Contact {
	model := contacttst.NewRandomPhone(t)
	model.UserID = userID
	TstCreate(t, conn, &model)
	return model
}
