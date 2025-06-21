package utils

import (
	"regexp"
	"strings"
)

func Slugify(input string) string {
	// Convert to lowercase
	slug := strings.ToLower(input)

	// Replace spaces with hyphens
	slug = strings.ReplaceAll(slug, " ", "-")

	// Remove special characters (keep alphanumeric and hyphens)
	slug = regexp.MustCompile(`[^a-z0-9-]`).ReplaceAllString(slug, "")

	// Trim leading and trailing hyphens
	slug = strings.Trim(slug, "-")

	return slug
}

func IsValidPackageName(name string) bool {
	// Java package name: dot-separated identifiers, each starting with a letter, followed by letters/digits/underscores, no Java keywords, no empty parts
	re := regexp.MustCompile(`^([a-zA-Z_][a-zA-Z0-9_]*)(\.[a-zA-Z_][a-zA-Z0-9_]*)*$`)
	return re.MatchString(name)
}

func BoolToString(val bool) string {
	if val {
		return "true"
	}
	return "false"
}

func ArrayContains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}
