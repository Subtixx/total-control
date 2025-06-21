return {
    GetInstalledMods = function(self)
        return {}
    end,
    GetInstalledModByID = function(self, modID)
        return nil
    end,
    DetectGameInstallation = function(self, path)
        if path == nil or path == "" then
            log.error("Path is nil or empty")
            return false
        end

        if not io.fileExists(path) then
            log.error(string.format("Path '%s' does not exist", path))
            return false
        end

        local executablePath = io.pathJoin(path, "Stardew Valley.dll")
        if not io.fileExists(executablePath) then
            return false
        end
        return true
    end,
    GetGameModDirectory = function(self)
        return ""
    end,
    AddMod = function(self, mod)

    end,
    RemoveMod = function(self, mod)

    end,
    UpdateMod = function(self, mod)
    end,
    GetGameID = function(self)
        return "stardew-valley"
    end,
    GetMods = function(self)
        return {}
    end,
    GetModByID = function(self, modID)
        return nil
    end
}