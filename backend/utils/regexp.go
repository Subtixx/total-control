package utils

import (
	"regexp"
	"strings"
)

func WildcardToRegex(input string) string {
	escapedInput := regexp.QuoteMeta(input)
	// Now we need to revert escaping *
	return strings.ReplaceAll(escapedInput, "\\*", ".*")
}
