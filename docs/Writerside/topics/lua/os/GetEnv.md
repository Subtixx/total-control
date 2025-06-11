# get_env

Returns the value of an environment variable.

## Syntax
```lua
string operating_system.get_env(string variable_name)
```

## Arguments
- **variable_name**: A key representing the environment variable you want to retrieve.

## Returns

The value of the environment variable if it exists, or `nil` if it does not.

## Example

```lua
local value = operating_system.get_env("VARIABLE_NAME")
if value then
    print("Value of VARIABLE_NAME: " .. value)
else
    print("VARIABLE_NAME is not set")
end
```