package utils

import (
	"os"
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

func IsWindows() bool {
	return GetOperatingSystem() == WindowsOS
}

func IsMacOS() bool {
	return GetOperatingSystem() == MacOS
}

func IsLinux() bool {
	return GetOperatingSystem() == LinuxOS
}
