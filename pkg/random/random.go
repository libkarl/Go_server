package random

import (
	"math/rand"
	"time"

	randomdata "github.com/Pallinder/go-randomdata"
)

var (
	randSrc = rand.NewSource(time.Now().UnixNano())
	random  = rand.New(randSrc)
)

// Gender defines random gender enum.
var Gender = randomdata.RandomGender

func init() {
	randomdata.CustomRand(random)
}
