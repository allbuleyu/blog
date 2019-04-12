package session

import (
	"fmt"
	"testing"
)

func Test_isCookieNameValid(t *testing.T) {
	tcs := []struct {
		input string
		ans bool
	}{
		{"@#$", true},
	}

	fmt.Println(isTokenTable)
	for _, tc := range tcs {
		fmt.Println(isCookieNameValid(tc.input))
	}
}
