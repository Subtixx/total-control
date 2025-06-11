# getOperatingSystem

Returns the current operating system as one of the constants defined above.

## Syntax
```lua
number operating_system.getOperatingSystem()
```

## Arguments

## Returns

A number representing the current operating system. The possible values are:
  - [`operating_system.Windows`](Windows.md)
  - [`operating_system.Linux`](Linux.md)
  - [`operating_system.MacOS`](MacOS.md)
  - [`operating_system.Unknown`](Unknown.md)

## Example

```lua
local os = operating_system.getOperatingSystem()
if os == operating_system.Windows then
    print("Running on Windows")
elseif os == operating_system.Linux then
    print("Running on Linux")
elseif os == operating_system.MacOS then
    print("Running on MacOS")
else
    print("Unknown operating system")
end 
```

<include from="lib.topic" element-id="lua.os.footer" use-filter="empty,all"></include>