package random

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"strconv"

	"github.com/Pallinder/go-randomdata"
)

const (
	letterBytes         = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterAlphaNumBytes = "abcdefghijklmnopqrstuvwxyz0123456789"
	letterIdxBits       = 6                    // 6 bits to represent a letter index
	letterIdxMask       = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax        = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var (
	currencies = []string{"BTC", "LTC", "DOGE", "USD"}
)

// Base64 returns random value encoded in base64 of 'bytesLength' length in bytes of creation.
func Base64(bytesLength int) string {
	return base64.StdEncoding.EncodeToString(Bytes(bytesLength))
}

// City returns random city name.
func City() string {
	return randomdata.City()
}

// CountryAlpha2 returns iso code of country
func CountryAlpha2() string {
	return randomdata.Country(randomdata.TwoCharCountry)
}

// CountryFull returns new random country full name.
func CountryFull() (country string) {
	generateCountryIfEmpty(&country)
	return country
}

// Currency returns random currency name.
func Currency(cryptoOnly bool) string {
	i := 2
	if !cryptoOnly {
		i = 3
	}
	return currencies[random.Intn(i)]
}

// Email return valid random email
func Email() string {
	return randomdata.Email()
}

// Firstname returns firstname
func Firstname(gender int) string {
	return randomdata.FirstName(gender)
}

// FromChartset returns a random string of n characters.
// Src: https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
func FromChartset(n int, charset string) string {
	b := make([]byte, n)
	// A randSrc.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, randSrc.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = randSrc.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(charset) {
			b[i] = charset[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return string(b)
}

// Hash generates random hexadecimal hash value with provided 'length'.
func Hash(length int) string {
	b := make([]byte, (length/2)+1)
	rand.Read(b)

	return hex.EncodeToString(b)[:length]
}

// Lastname returns lastname
func Lastname() string {
	return randomdata.LastName()
}

// LowercaseAlphaNumeric returns a random string with alpha numeric
func LowercaseAlphaNumeric(n int) string {
	return FromChartset(n, letterAlphaNumBytes)
}

// String returns a random string of n characters.
// Src: https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
func String(n int) string {
	return FromChartset(n, letterBytes)
}

// StringInt string of random number
func StringInt(l int) string {
	str := ""
	for i := 0; i < l; i++ {
		str += strconv.FormatInt(int64(Int(0, 9)), 10)
	}
	return str
}

// Paragraph return random paragraph of text.
func Paragraph() string {
	return randomdata.Paragraph()
}

// Phone return valid random phone.
func Phone() string {
	// TODO
	// return randomdata.PhoneNumber()
	return "+420773" + StringInt(6)
}

// PostalCode returns random postal code for given 'countryCode'.
// If no country code provided random one is generated.
func PostalCode(countryCode string) (postalCode string) {
	generateCountryCodeIfEmpty(&countryCode)
	if postalCode = randomdata.PostalCode(countryCode); postalCode == "" {
		postalCode = randomdata.PostalCode("US")
	}
	return postalCode
}

// Region returns random region for given 'countryCode'.
// If no country code provided random one is generated.
// If country code is not recognized or there is no region for given country
// the function returns silly name.
func Region(countryCode string) (region string) {
	generateCountryCodeIfEmpty(&countryCode)
	region = randomdata.ProvinceForCountry(countryCode)
	if region == "" {
		region = randomdata.SillyName()
	}
	return region
}

// SillyName returns silly name
func SillyName() string {
	return randomdata.SillyName()
}

// Street returns random street name.
// If no country code provided random one is generated.
// If the country code is not supported then the street name is randomly generated.
func Street(countryCode string) (street string) {
	generateCountryCodeIfEmpty(&countryCode)
	street = randomdata.StreetForCountry(countryCode)
	if street == "" {
		street = randomdata.Street()
	}
	return street
}

// Username generates random valid username.
func Username() (username string) {
	for len(username) <= 2 || len(username) >= 20 {
		username = randomdata.SillyName()
	}
	return username
}

func generateCountryIfEmpty(country *string) {
	if *country == "" {
		*country = randomdata.Country(randomdata.FullCountry)
	}
}

func generateCountryCodeIfEmpty(countryCode *string) {
	if *countryCode == "" {
		*countryCode = randomdata.Country(randomdata.TwoCharCountry)
	}
}
