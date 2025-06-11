# readFilesFromZip

## Syntax

```lua
table input_output.readFilesFromZip(zipFile)
```

Reads all files from a zip archive.

## Arguments

- **zipFile**: The zip archive to read from.

## Returns

Returns a table containing the paths of all files and their contents in the zip archive.

## Example

This example reads all files from the zip archive `myArchive.zip` and prints their contents.

```lua
local files = input_output.readFilesFromZip("myArchive.zip")
for filePath, content in pairs(files) do
    print("File: " .. filePath)
    print("Content: " .. content)
end
```
