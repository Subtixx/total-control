plugin = {
    GetMods = function()
        return {
            {
                id = "mod1",
                name = "Test Mod",
                enabled = true,
                game_id = "game1"
            }
        }
    end,
    GetModByID = function(id)
        -- Example: return a mod by id
    end,
    GetGameModDirectory = function()
        -- This is usually located at:
        -- - Linux: ~/.factorio/mods/
        -- - Windows: C:\Users\<Username>\AppData\Roaming\Factorio\mods\
        -- - macOS: ~/Library/Application Support/factorio/mods/
        if operating_system.is_windows then
            local appdata = os.getenv("APPDATA")
            if appdata then
                return appdata .. "\\Factorio\\mods\\"
            end
        elseif operating_system.is_linux then
            return os.getenv("HOME") .. "/.factorio/mods/"
        elseif operating_system.is_macos then
            return os.getenv("HOME") .. "/Library/Application Support/factorio/mods/"
        end
        return nil -- Unsupported OS
    end,
    AddMod = function(mod)
        -- Use write_file/read_file for I/O
    end,
    RemoveMod = function(id)
    end,
    UpdateMod = function(mod)
    end,
    ListGameMods = function()
        return {}
    end,
    GetGameID = function()
        return "game1"
    end
}