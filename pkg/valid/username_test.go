package valid

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUsername(t *testing.T) {
	assert.True(t, Username("okd"))
	assert.True(t, Username("okd2"))
}
