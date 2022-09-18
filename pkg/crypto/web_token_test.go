package crypto

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWebToken(t *testing.T) {
	token, err := CreateToken("pass", uint(1), time.Second/2)
	assert.Empty(t, err)
	id, err := DecodeToken("pass", "bearer "+token)
	assert.Empty(t, err)
	assert.Equal(t, uint(1), id)
}

func TestWebTokenInvalid(t *testing.T) {
	token, err := CreateToken("test", uint(1), time.Second/2)
	assert.Empty(t, err)
	_, err = DecodeToken("pass", "bearer "+token)
	assert.Equal(t, "signature is invalid", err.Error())
}

func TestWebTokenExpiration(t *testing.T) {
	token, err := CreateToken("test", uint(1), -time.Second)
	assert.Empty(t, err)
	_, err = DecodeToken("test", "bearer "+token)
	assert.Equal(t, "Token is expired", err.Error())
}
