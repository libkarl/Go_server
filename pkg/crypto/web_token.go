package crypto

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Strips 'Bearer ' prefix from bearer token string
func stripBearer(tok string) string {
	// Should be a bearer token
	if len(tok) > 6 && strings.ToUpper(tok[0:7]) == "BEARER " {
		return tok[7:]
	}
	return tok
}

// CreateToken ...
func CreateToken(password string, id uint, expiration time.Duration) (string, error) {
	token := jwt.New(jwt.GetSigningMethod("HS256"))
	// Set some claims
	exp := time.Now().Add(expiration).Unix()
	m := jwt.MapClaims{
		"Id": strconv.FormatUint(uint64(id), 10),
	}
	if expiration != 0 {
		m["exp"] = exp
	}
	token.Claims = m
	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString([]byte(password))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// DecodeToken ...
func DecodeToken(password string, token string) (uint, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		b := ([]byte(password))
		return b, nil
	}
	token = stripBearer(token)
	tokenParsed, err := jwt.ParseWithClaims(token, jwt.MapClaims{}, keyFunc)
	if err != nil {
		return 0, err
	}
	if !tokenParsed.Valid {
		return 0, errors.New("session is no longer valid")
	}
	claims := tokenParsed.Claims.(jwt.MapClaims)
	id, err := strconv.ParseUint(claims["Id"].(string), 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}
