package valid

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/investapp/backend/pkg/errdef"
)

func TestFlags(t *testing.T) {
	f1 := Flag("hello")
	f2 := Flag("world")
	set := Flags([]Flagger{f1, f2})
	assert.True(t, set.Has(f1))
	assert.True(t, set.Has(f2))
	assert.False(t, set.Has(Flag("nope")))
}

type input struct {
	err      *errdef.Error
	callback func() errdef.Error
}

var _ Validator = &input{}

func (i *input) Validate() errdef.Error {
	if i.callback != nil {
		return i.callback()
	}
	return i.err
}

type inputT2 struct {
	err      errdef.Error
	callback func(Flags) errdef.Error
}

var _ validatorWithFlags = &inputT2{}

func (i *inputT2) Validate(f Flags) errdef.Error {
	if i.callback != nil {
		return i.callback(f)
	}
	return i.err
}

func TestValidate(t *testing.T) {
	i1 := &input{}
	assert.Nil(t, Validate(i1))
	err := errdef.ErrInternalf("validation", "very bad bad validation")
	i2 := &input{err: err}
	assert.Equal(t, err, Validate(i2))
	f1 := Flag("hello")
	i3 := &inputT2{
		callback: func(ff Flags) errdef.Error {
			assert.True(t, ff.Has(f1))
			return *err
		},
	}
	assert.Equal(t, err, Validate(i3, f1))
	i4 := &inputT2{
		callback: func(ff Flags) errdef.Error {
			assert.False(t, ff.Has(f1))
			return nil
		},
	}
	assert.Equal(t, nil, Validate(i4))
}
