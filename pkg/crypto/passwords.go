package crypto

import "golang.org/x/crypto/bcrypt"

func clear(b []byte) {
	for i := 0; i < len(b); i++ {
		b[i] = 0
	}
}

// Crypt ...
func Crypt(password []byte) ([]byte, error) {
	defer clear(password)
	return bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
}

// CompareCrypts ...
func CompareCrypts(pwd1, pwd2 []byte) bool {
	err := bcrypt.CompareHashAndPassword(pwd1, pwd2)
	return err == nil
}
