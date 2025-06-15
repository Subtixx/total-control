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

func BoolToString(val bool) string {
	if val {
		return "true"
	}
	return "false"
}
