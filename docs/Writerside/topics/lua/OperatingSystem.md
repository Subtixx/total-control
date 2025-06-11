# Operating System

The `operating_system` table contains useful functions and constants related to the os the application is running on.

## Constants

### Windows

`operating_system.Windows` A LUA Number constant representing the Windows operating system.

### Linux

`operating_system.Linux` A LUA Number constant representing the Linux operating system.

### MacOS

`operating_system.MacOS` A LUA Number constant representing the MacOS operating system.

### Unknown

`operating_system.Unknown` A LUA Number constant representing an unknown operating system.

## Methods

### getOperatingSystem

Returns the current operating system as one of the constants defined above.

```lua
local os = operating_system.getOperatingSystem()
```

## is_windows

Holds a boolean value indicating if the current operating system is Windows.

```lua
local isWindows = operating_system.is_windows
```

### is_linux

Holds a boolean value indicating if the current operating system is Linux.

```lua
local isLinux = operating_system.is_linux
```

### is_macos

Holds a boolean value indicating if the current operating system is MacOS.

```lua
local isMacOS = operating_system.is_macos
```

### is_unknown

Holds a boolean value indicating if the current operating system is unknown.

```lua
local isUnknown = operating_system.is_unknown
```

## getenv

Get the value of an environment variable.

```lua
local value = operating_system.getenv("MY_ENV_VAR")
```