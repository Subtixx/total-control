# Plugin

This page describes how to create a plugin.
  
> All API calls are synchronous and may return errors if called incorrectly.
{style="warning"}

## Full Example

```lua
local mods = {
    {
        ID = "mod1",
        Name = "Mod One",
        Author = "Author A",
        Version = "1.0",
        Enabled = true,
        GameVersions = { { Version = "1.0" } }
    },

    {
        ID = "mod2",
        Name = "Mod Two",
        Author = "Author B",
        Version = "1.1",
        Enabled = false,
        GameVersions = { { Version = "1.1" } }
    },
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
