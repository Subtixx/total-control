plugin = {
    mods = nil,
    GetMods = function(self)
        if self.mods ~= nil then
            return self.mods
        end
        local mod_files = input_output.getFilesInDirectory(
                self:GetGameModDirectory(),
                { "*.zip" }
        )
        if mod_files == nil or #mod_files == 0 then
            log.warn("No mods found in the mods directory " ..
                    self:GetGameModDirectory())
            return {}
        end
        log.info("Found " .. #mod_files .. " mod files in the mods directory.")
        self.mods = {}
        for _, mod_file in ipairs(mod_files) do
            local info_file = input_output.readFileFromZip(mod_file, ".*?/info\\.json")
            local modInfo = json.decode(info_file)
            self.mods[#self.mods + 1] = {
                id = modInfo.name,
                game_id = self:GetGameID(),
                name = modInfo.title or modInfo.name,
                version = modInfo.version,
                description = modInfo.description or "",
                author = modInfo.author or "Unknown",
                file_path = mod_file,
                enabled = true,
            }
        end
        if self.mods == nil or #self.mods == 0 then
            log.warn("No mods found in the mods directory " ..
                    self:GetGameModDirectory() .. tostring(#self.mods))
            return {}
        end
        return self.mods
    end,
    GetModByID = function(self, id)
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
    AddMod = function(self, mod)
        -- Use write_file/read_file for I/O
    end,
    RemoveMod = function(self, id)
    end,
    UpdateMod = function(self, mod)
    end,
    ListGameMods = function(self)
        return {}
    end,
    GetGameID = function(self)
        return "factorio"
    end
}