package valid

import (
	"fmt"

	"github.com/investapp/backend/pkg/errdef"
)

// Flagger is same as fmt.Stringer, just custom type
type Flagger fmt.Stringer

// Flag is implementation of Flagger or fmt.Stringer with just simple string method.
// String() in this case is used for comparism.
type Flag string

// Name to satisfy Flagger or fmt.Stringer interface
func (s Flag) String() string {
	return string(s)
}

// Flags helps you to manage Flag type to put multiple flags
// into validation process. Adjusting validation vector
// is useful if you have same input, but based on endpoint you need
// fields to be validated differently
type Flags []Flagger

// Get will return either Flagger if found or nil
// if not found.
func (ff Flags) Get(flag string) Flagger {
	for _, f := range ff {
		if f.String() == flag {
			return f
		}
	}
	return nil
}

// Has return if list of flags has flag you are looking for
func (ff Flags) Has(flag Flag) bool {
	return ff.Get(flag.String()) != nil
}

// Validator is an interface used to validate structures.
// If failed it returns errdef.ErrSet error.
type Validator interface {
	Validate() errdef.ErrSet
}

type validatorWithFlags interface {
	Validate(Flags) errdef.ErrSet
}

// Validate will take input and run validations if input has any.
// You can also specify flags which can move validation vector, depends on your
// need of validation
func Validate(result interface{}, flags ...Flagger) errdef.ErrSet {
	v1, ok := result.(Validator)
	if ok {
		return v1.Validate()
	}
	v2, ok := result.(validatorWithFlags)
	if ok {
		return v2.Validate(Flags(flags))
	}
	return nil
}
