package helpers

import "strings"

func StandardizeSpaces(s string) string {
	return strings.Join(strings.Fields(s), " ")
}
