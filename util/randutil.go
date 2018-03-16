package util

import (
	"math/rand"
)

const (

	s = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
)

var (	
	letters []rune
	length int
)

func init() {
	letters = []rune(s)
	length = len(letters)
}

func RandAlphabetic(count int) string {
	s := ""
	for i := 0; i < count; i++ {
		ra := rand.Intn(length)
		s += string(letters[ra])
	}
	return s
}
