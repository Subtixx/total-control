//go:build !windows

package utils

import "path/filepath"

func IsHidden(path string) bool {
	if path == "" {
		return false
	}

	base := filepath.Base(path)
	return base[0] == '.'
}
