# readFileFromZip

## Syntax

```lua
string input_output.readFileFromZip(zipFile, fileName)
```

Reads a single file from a zip archive.

## Arguments

- **zipFile**: The zip archive to read from.
- **fileName**: The name of the file or a regex to read from the zip archive.

## Returns

Returns the content of the file as a string.

## Example

This example reads the first file it finds in the zip archive `myArchive.zip`
that matches the regex `.*?/myFile.txt`.

```lua
local content = input_output.readFileFromZip("myArchive.zip", ".*?/myFile.txt")
print(content)
```