package random

import (
	"testing"

	"github.com/investapp/backend/pkg/valid"
	"github.com/stretchr/testify/assert"
)

func TestWordlist(t *testing.T) {
	email := WordlistEmail(WordlistFullName())
	assert.True(t, len(email) > 0)
	assert.Contains(t, email, "@")
}

func TestUsername(t *testing.T) {
	for i := 0; i <= 10000; i++ {
		name := WordlistUsername(WordlistFullName())
		assert.True(t, valid.Username(name), name)
	}
}
