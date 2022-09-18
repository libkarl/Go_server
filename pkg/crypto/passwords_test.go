package crypto

import (
	"testing"

	"golang.org/x/crypto/bcrypt"

	"github.com/stretchr/testify/assert"
)

func TestCryptFunctionMatch(t *testing.T) {
	pwd1, _ := Crypt([]byte("testtest"))
	pwd2 := []byte("testtest")
	assert.Equal(t, true, CompareCrypts(pwd1, pwd2))
}

func TestCryptFunctionNotMatch(t *testing.T) {
	pwd1, _ := Crypt([]byte("testtest"))
	pwd2 := []byte("testtest1")
	assert.Equal(t, false, CompareCrypts(pwd1, pwd2))
}

func TestBcryptLibrary(t *testing.T) {
	password := []byte("MyDarkSecret")
	// Hashing the password with the default cost of 10
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	// Comparing the password with the hash
	err = bcrypt.CompareHashAndPassword(hashedPassword, password)
	assert.Equal(t, nil, err)
}
