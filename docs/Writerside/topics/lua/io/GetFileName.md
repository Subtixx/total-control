# getFileName

## Syntax

```lua
string getFileName(string path)
```

## Description

Returns the name of the file from a given path. If the path is a directory, it returns an empty string.

## Parameters

- `path`: A string representing the file path from which to extract the file name.

## Returns

- A string containing the file name if the path is valid and points to a file; otherwise, it returns an empty string.

## Example

```lua
local fileName = getFileName("/path/to/file.txt")
if fileName ~= "" then
    print("File name is: " .. fileName)
else
    print("The path does not point to a file.")
end
```