# Mod Object

## id (string)

The unique identifier for the mod.

## name (string)

The name of the mod.

## enabled (boolean)

Indicates whether the mod is enabled (`true`) or disabled (`false`).

## game_id (string)

The ID of the game this mod is associated with.

## Full Example

```lua
local mod = {
    id = "example_mod",
    name = "Example Mod",
    enabled = true,
    game_id = "game123"
}
```