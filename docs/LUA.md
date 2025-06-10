# Lua API Reference

All functions of the plugin need to be implemented in a global `plugin` table.

## Overview

This document provides an overview of the Lua API available for managing mods in the plugin system.
Each function is described with its parameters, return values, and examples.

---

## Mod Object Structure

Each mod object has the following fields:

- `id` (string)
- `name` (string)
- `enabled` (boolean)
- `game_id` (string)

## plugin.GetMods()

**Returns:**  
A table (array) of mod objects.

**Example:**

```lua
function GetMods()
    return {
        {id = "mod1", name = "Mod One", enabled = true, game_id = "game1"},
        {id = "mod2", name = "Mod Two", enabled = false, game_id = "game2"},
    }
end
```

---

## plugin.GetModByID(id)

**Parameters:**

- `id` (string): The mod ID.

**Returns:**  
A mod object or `nil` if not found.

**Example:**

```lua
function GetModByID(id)
    local mods = plugin.GetMods()
    for _, mod in ipairs(mods) do
        if mod.id == id then
            return mod
        end
    end
    return nil
end
```

---

## plugin.AddMod(mod)

**Parameters:**

- `mod` (table): A mod object with fields `id`, `name`, `enabled`, `game_id`.

**Returns:**  
`true` on success.

**Example:**

```lua
function AddMod(mod)
    if not mod.id or not mod.name or not mod.game_id then
        error("Invalid mod object")
    end
    -- Add the mod to the system (implementation depends on your system)
    print("Adding mod:", mod.name)
    return true
end
```

---

## plugin.RemoveMod(id)

**Parameters:**

- `id` (string): The mod ID.

**Returns:**  
`true` on success.

**Example:**

```lua
function RemoveMod(id)
    local mods = plugin.GetMods()
    for i, mod in ipairs(mods) do
        if mod.id == id then
            table.remove(mods, i)
            print("Removed mod:", mod.name)
            return true
        end
    end
    print("Mod not found:", id)
    return false
end
```

---

## plugin.UpdateMod(mod)

**Parameters:**

- `mod` (table): A mod object with updated fields.

**Returns:**  
`true` on success.

**Example:**

```lua
function UpdateMod(mod)
    if not mod.id or not mod.name or not mod.game_id then
        error("Invalid mod object")
    end
    local mods = plugin.GetMods()
    for i, existingMod in ipairs(mods) do
        if existingMod.id == mod.id then
            mods[i] = mod
            print("Updated mod:", mod.name)
            return true
        end
    end
    print("Mod not found for update:", mod.id)
    return false
end
```

---

## plugin.GetGameModDirectory()

**Returns:**  
A string with the path to the mods directory.

**Example:**

```lua
function GetGameModDirectory()
    return "/path/to/game/mods"
end
```

---

## plugin.GetGameID()

**Returns:**  
A string with the game ID.

**Example:**

```lua
function GetGameID()
    return "game123"
end
```

---

## Logging

## log.debug(message)

**Parameters:**

- `message` (string): The debug message to log.

## log.info(message)

**Parameters:**

- `message` (string): The message to log.

## log.warn(message)

**Parameters:**

- `message` (string): The warning message to log.

## log.error(message)

**Parameters:**

- `message` (string): The error message to log.

## Helper Functions

### os_getenv(name)

**Parameters:**

- `name` (string): The name of the environment variable.

**Returns:**

- The value of the environment variable or `nil` if not set.

```lua
local value = os.getenv("MY_ENV_VAR")
if value then
    print("Environment variable MY_ENV_VAR:", value)
else
    print("MY_ENV_VAR is not set.")
end
```

---

## Full example

```lua
local mods = {
    {ID = "mod1", Name = "Mod One", Author = "Author A", Version = "1.0", Enabled = true, GameVersions = { {Version = "1.0"} }},
    {ID = "mod2", Name = "Mod Two", Author = "Author B", Version = "1.1", Enabled = false, GameVersions = { {Version = "1.1"} }},
}
plugin = {
    GetMods = function()
        return mods
    end,
    GetModByID = function(id)
        for _, mod in ipairs(mods) do
            if mod.ID == id then
                return mod
            end
        end
    end,
    AddMod = function(mod)
        print("Adding mod:", mod.Name)
        return true
    end,
    RemoveMod = function(id)
        for i, mod in ipairs(mods) do
            if mod.ID == id then
                table.remove(mods, i)
                print("Removed mod:", mod.Name)
                return true
            end
        end
        print("Mod not found:", id)
        return false
    end,
    UpdateMod = function(mod)
        for i, existingMod in ipairs(mods) do
            if existingMod.ID == mod.ID then
                mods[i] = mod
                print("Updated mod:", mod.Name)
                return true
            end
        end
        print("Mod not found for update:", mod.ID)
        return false
    end,
    GetGameModDirectory = function()
        return "/path/to/game/mods"
    end,
    GetGameID = function()
        return "game123"
    end,
}
```

**Note:**  
All API calls are synchronous and may return errors if called incorrectly.