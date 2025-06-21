package utils

import (
	"regexp"
	"strconv"
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

func StringToInt(s string) int {
	if s == "" {
		return 0
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		return 0 // Return 0 if conversion fails
	}

	return i
}

func StringToInt64(s string) int64 {
	if s == "" {
		return 0
	}

	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0 // Return 0 if conversion fails
	}

	return i
}

func StringToFloat(s string) float64 {
	if s == "" {
		return 0.0
	}

	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0.0 // Return 0.0 if conversion fails
	}

	return f
}

func ArrayContains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}
