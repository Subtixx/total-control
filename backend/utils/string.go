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

func StringToInt(s interface{}) int {
	if s == nil {
		return 0
	}

	str, ok := s.(string)
	if !ok {
		return 0
	}

	i, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}

	return i
}

func StringToInt64(s interface{}) int64 {
	if s == nil {
		return 0
	}

	str, ok := s.(string)
	if !ok {
		return 0
	}

	if str == "" {
		return 0
	}

	i, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0
	}

	return i
}

func StringToFloat(s interface{}) float64 {
	if s == nil {
		return 0.0
	}

	str, ok := s.(string)
	if !ok {
		return 0.0
	}

	f, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0.0
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
