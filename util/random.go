package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const alphabets = "abcdefghijklmnopqrstuvwxyz"
var rnd = rand.New(rand.NewSource(time.Now().UnixNano()))

// Return random owner
func RandOwner() string {
	return RandStr(6)
}

// Return random balance
func RandBalance() int64 {
	return RandomInt(0, 1000)
}

func RandCurrency() string {
	curr := []string{
		"INR", "USD", "CAD", "YEN",
	}
	n := int64(len(curr))
	
	return curr[rnd.Int63n(n)]

}



// ---- BASE ----- //

// Random int generator
func RandomInt(min, max int64) int64 {
	return min + rnd.Int63n(max-min+1)
}

// Random str generator
func RandStr(len int) string {
	var randStr strings.Builder
	k := int64(26)

	for i:= 0; i < len; i++{
		char := alphabets[rnd.Int63n(k)]
		randStr.WriteByte(char)
	}
	return randStr.String() 
}

// function that returns random email
func RandEmail () string{
	return fmt.Sprintf("%s@email.com", RandStr(6))
}