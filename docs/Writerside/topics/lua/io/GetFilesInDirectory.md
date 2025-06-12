# getFilesInDirectory

## Syntax

```lua
table input_output.getFilesInDirectory(string directory, table|string pattern)
```
Gets a list of files in a directory that match a given pattern.

## Arguments

- **directory**: The directory to search in.
- **pattern**: A string or a table of strings that specifies the file names to match. If a table is provided, it will match any of the patterns in the table. If a string is provided, it will be treated as a regex pattern.

## Returns

Returns a table containing the names of the files that match the pattern. If no files match, an empty table is returned.

## Example
This example retrieves all `.txt` files in the directory `myDirectory`.

```lua
local files = input_output.getFilesInDirectory("myDirectory", "%.txt$")
for _, file in ipairs(files) do
    print(file)
end
```