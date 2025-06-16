package utils

import (
	"os"
	"path/filepath"
	"runtime"
)

const (
	OperatingSystemWindows = "windows"
	OperatingSystemMac     = "darwin"
	OperatingSystemLinux   = "linux"
	OperatingSystemUnknown = "unknown"
)

func GetOperatingSystem() string {
	switch runtime.GOOS {
	case "windows":
		return OperatingSystemWindows
	case "darwin":
		return OperatingSystemMac
	case "linux":
		return OperatingSystemLinux
	default:
		return OperatingSystemUnknown
	}
}

// GetCommonUserDataPath returns the common user data path based on the operating system.
// Windows: %LOCALAPPDATA%
// macOS: ~/Library/Application Support
// Linux: ~/.local/share
func GetCommonUserDataPath() string {
	if GetOperatingSystem() == OperatingSystemWindows {
		return os.Getenv("LOCALAPPDATA")
	} else if GetOperatingSystem() == OperatingSystemMac {
		return os.Getenv("HOME") + "/Library/Application Support"
	} else if GetOperatingSystem() == OperatingSystemLinux {
		return os.Getenv("HOME") + "/.local/share"
	}

	panic("Unsupported OS: " + GetOperatingSystem())
}

// GetAppDataPath returns the application data path for TotalControl.
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
	return GetOperatingSystem() == OperatingSystemWindows
}

func IsMac() bool {
	return GetOperatingSystem() == OperatingSystemMac
}

func IsLinux() bool {
	return GetOperatingSystem() == OperatingSystemLinux
}
