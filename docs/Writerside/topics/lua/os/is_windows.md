# is_windows

## `operating_system.is_windows`

A boolean value indicating if the current operating system is Windows.

### Example

```lua
local isWindows = operating_system.is_windows
if isWindows then
    print("This script is running on Windows.")
else
    print("This script is not running on Windows.")
end
```

<include from="lib.topic" element-id="lua.os.footer" use-filter="empty,windows"></include>