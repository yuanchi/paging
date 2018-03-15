package randutil

import (
	"testing"
	"fmt"
)

func TestRandAlphabetic(t *testing.T) {
	for i := 0; i < 10; i++ {
		s := RandAlphabetic(5)
		fmt.Println(s)
	}
	
}
