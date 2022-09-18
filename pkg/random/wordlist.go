package random

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/investapp/backend/pkg/random/wordlist"
	"github.com/investapp/backend/pkg/valid"
)

var firstnameLst, lastnameLst, emaildomainsLst []string

func init() {
	b1, _ := wordlist.Asset("firstnames.txt")
	firstnameLst = strings.Split(string(b1), "\n")

	b2, _ := wordlist.Asset("lastnames.txt")
	lastnameLst = strings.Split(string(b2), "\n")

	b3, _ := wordlist.Asset("emaildomains.txt")
	emaildomainsLst = strings.Split(string(b3), "\n")
}

// Wordlist generates random data based on embeded
// wordlists
type Wordlist struct {
	rand *rand.Rand
}

// NewWordlistRand creates new instance
func NewWordlistRand(rand *rand.Rand) *Wordlist {
	return &Wordlist{rand: rand}
}

// FullNameLst ...
type FullNameLst []string

func (fn FullNameLst) String() string {
	var strs []string
	for _, s := range fn {
		if len(s) == 0 {
			continue
		}
		strs = append(strs, strings.Title(s))
	}
	return strings.Join(strs, " ")
}

// FirstName is a firstname
func (fn FullNameLst) FirstName() string {
	if len(fn) > 0 {
		return strings.Title(fn[0])
	}
	return ""
}

// LastName is a lastname
func (fn FullNameLst) LastName() string {
	if len(fn) > 0 {
		return strings.Title(fn[len(fn)-1])
	}
	return ""

}

// WordlistFullName generates fullname from wordlist
func (w *Wordlist) WordlistFullName() FullNameLst {
	i1 := w.rand.Intn(len(firstnameLst) - 1)
	i2 := w.rand.Intn(len(lastnameLst) - 1)
	name := FullNameLst{firstnameLst[i1], lastnameLst[i2]}
	return name
}

// WordlistEmailDomain will return email domain out of wordlist
func (w *Wordlist) WordlistEmailDomain() string {
	i1 := w.rand.Intn(len(emaildomainsLst) - 1)
	return emaildomainsLst[i1]
}

// WordlistEmail generates email from wordlist
func (w *Wordlist) WordlistEmail(name FullNameLst) (email string) {
	domain := w.WordlistEmailDomain()
	firstname := strings.ToLower(name.FirstName())
	lastname := strings.ToLower(name.LastName())

	defer func() {
		// This is just a protection to always met valid username
		for {
			if valid.Email(email) {
				break
			}
			email = w.WordlistEmail(name)
		}
	}()

	switch w.rand.Intn(2) {
	case 0:
		email = fmt.Sprintf("%s.%s@%s", firstname, lastname, domain)
		return
	case 1:
		email = fmt.Sprintf("%s%s@%s", firstname, lastname, domain)
		return
	case 2:
		email = fmt.Sprintf("%s%d@%s", lastname, w.rand.Intn(10), domain)
		return
	default:
		email = fmt.Sprintf("%s.%s@%s", firstname, lastname, domain)
		return
	}
}

// WordlistUsername generates random username
func (w *Wordlist) WordlistUsername(name FullNameLst) (username string) {
	firstname := strings.ToLower(name.FirstName())
	lastname := strings.ToLower(name.LastName())
	number := strconv.FormatInt(int64(w.rand.Intn(99)), 10)

	defer func() {
		// This is just a protection to always met valid username
		for {
			if valid.Username(username) {
				break
			}
			username = w.WordlistUsername(name)
		}
	}()

	switch w.rand.Intn(15) {
	case 0:
		username = firstname + "ky" + number
		return
	case 1:
		username = lastname + "ky" + number
		return
	case 2:
		username = firstname + number
		return
	case 3:
		username = lastname + number
		return lastname + number
	case 4:
		username = Username()
		return
	case 5:
		username = Username() + number
		return
	case 6:
		username = firstname + "." + lastname
		return
	case 7:
		username = firstname + "." + lastname
		return
	case 8:
		username = "x" + firstname
		return
	case 9:
		username = "x" + lastname
		return
	case 10:
		username = firstname + "_" + lastname
		return
	case 11:
		username = lastname + "_" + firstname
		return
	case 12:
		username = lastname + "_" + number
		return
	case 13:
		username = firstname + "_" + number
		return
	case 14:
		username = firstname + lastname
		return
	default:
		username = firstname + lastname
		return
	}
}

// WordlistFullName generates fullname from wordlist
func WordlistFullName() FullNameLst {
	return NewWordlistRand(random).WordlistFullName()
}

// WordlistEmailDomain will return email domain out of wordlist
func WordlistEmailDomain() string {
	return NewWordlistRand(random).WordlistEmailDomain()
}

// WordlistEmail generates email from wordlist
func WordlistEmail(name FullNameLst) (email string) {
	return NewWordlistRand(random).WordlistEmail(name)
}

// WordlistUsername generates random username
func WordlistUsername(name FullNameLst) (email string) {
	return NewWordlistRand(random).WordlistUsername(name)
}
