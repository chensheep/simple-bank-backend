package util

import (
	"math/rand"
	"strings"
	"time"

	"github.com/bxcodec/faker/v3"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		sb.WriteByte(alphabet[rand.Intn(k)])
	}
	return sb.String()
}

func RandomOwner() string {
	return RandomString(8)
}

func RandomBalance() int64 {
	return RandomInt(0, 2000)
}

func RandomMoney() int64 {
	return RandomInt(0, 2000)
}

func RandomCurrency() string {
	currencies := []string{USD, EUR, TWD}
	return currencies[rand.Intn(len(currencies))]
}

func RandomEmail() string {
	return faker.Email()
}
