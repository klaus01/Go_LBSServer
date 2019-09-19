package utils

import (
	"math/rand"
	"time"
)

// RandomInt 返回随机整数 [0, max]
func RandomInt(max int) int {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	return r.Intn(max)
}
