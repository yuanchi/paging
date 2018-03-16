package util

import (
	"fmt"
	"testing"
)

func TestCurrentTimef(t *testing.T) {
	f := CurrentTimef()
	fmt.Println(f)
}
