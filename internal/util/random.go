package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const lowercase = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// randomString generates a random string of length n.
func randomString(n int) string {
	var sb strings.Builder
	k := len(lowercase)

	for i := 0; i < n; i++ {
		c := lowercase[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// RandomUsername generates a random username.
func RandomUsername() string {
	return randomString(6)
}

// RandomFullName generates a random full name.
func RandomFullName() string {
	return fmt.Sprintf("%s %s", randomString(3), randomString(3))
}

// RandomEmail generates a random email.
func RandomEmail() string {
	return fmt.Sprintf("%s@email.com", randomString(6))
}
