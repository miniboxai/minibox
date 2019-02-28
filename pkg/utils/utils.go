package utils

import "strings"

func Empty(str string) bool {
	return len(strings.TrimSpace(str)) == 0
}
