package utils

import (
	"strings"
)

// IsBlank returns true if the given string empty or only whitespace, false otherwise
func IsBlank(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}
