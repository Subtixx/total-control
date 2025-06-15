package utils

import (
	"os"
	"path/filepath"
	"runtime"
)

const (
	WindowsOS = "windows"
	MacOS     = "darwin"
	LinuxOS   = "linux"
	UnknownOS = "unknown"
)

func GetOperatingSystem() string {
	switch runtime.GOOS {
	case "windows":
		return WindowsOS
	case "darwin":
		return MacOS
	case "linux":
		return LinuxOS
	default:
		return UnknownOS
	}
}

// GetCommonUserDataPath returns the common user data path based on the operating system.
// Windows: %LOCALAPPDATA%
// macOS: ~/Library/Application Support
// Linux: ~/.local/share
func GetCommonUserDataPath() string {
	if GetOperatingSystem() == WindowsOS {
		return os.Getenv("LOCALAPPDATA")
	} else if GetOperatingSystem() == MacOS {
		return os.Getenv("HOME") + "/Library/Application Support"
	} else if GetOperatingSystem() == LinuxOS {
		return os.Getenv("HOME") + "/.local/share"
	}

	panic("Unsupported OS: " + GetOperatingSystem())
}

func GetAppDataPath() string {
	// Get the user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err) // Handle error appropriately in production code
	}

	// Construct the app data path
	appDataPath := filepath.Join(homeDir, ".totalcontrol")

	// Create the directory if it doesn't exist
	if err := CreateDirectoryIfNotExists(appDataPath); err != nil {
		panic(err) // Handle error appropriately in production code
	}

	return appDataPath
}

func IsWindows() bool {
	return GetOperatingSystem() == WindowsOS
}

func IsMacOS() bool {
	return GetOperatingSystem() == MacOS
}

func IsLinux() bool {
	return GetOperatingSystem() == LinuxOS
}
