function readModListFile(plugin)
    local mod_list_file = io.getFileContent(
            plugin:GetGameModDirectory() .. "mod-list.json"
    )
    if mod_list_file == nil or mod_list_file == "" then
        log.warn("No mod-list.json found in the mods directory " ..
                plugin:GetGameModDirectory())
        return {}
    end
    local mod_list = json.decode(mod_list_file)
    if mod_list == nil or mod_list.mods == nil or #mod_list.mods == 0 then
        log.warn("No mods found in the mod-list.json file.")
        return {}
    end

    local mod_ids = {}
    for _, mod in ipairs(mod_list.mods) do
        if mod and mod.name then
            mod_ids[mod.name] = mod.enabled or false
        end
    end

    return mod_ids
end

function loadModsFromApi()

    -- Do Http requests to https://mods.factorio.com/api/mods
    local response = http.get("https://mods.factorio.com/api/mods")
    if response == nil then
        log.error("Failed to fetch mods from Factorio API: " .. tostring(response))
        return {}
    end
    if response.status_code ~= 200 then
        log.error("Failed to fetch mods from Factorio API: " .. response.status_code)
        return {}
    end

    local pagination = response.body.pagination
    local results = response.body.results
    print(serpent.block(response.body))
    cache.set("factorio_mods", {
        pagination = pagination,
        results = results,
    })
    return response.body
end

function loadMods()
    if cache.has("factorio_mods") then
        -- If the cache exists, return it
        local cached_mods = cache.get("factorio_mods")
        if cached_mods and cached_mods.results then
            print("Got " .. table_size(cached_mods.results) ..
                    " mods from cache.")
            return cached_mods
        end
    end
    print("Cache not found, loading mods from API...")
    return loadModsFromApi()
end

return {
    -- The first load of the plugins will be slower.
    mods = loadMods(),
    GetInstalledMods = function(self)
        print("GetInstalledMods called")
        if self.mods ~= nil then
            return self.mods
        end

        local mod_list = readModListFile(self)
        local mod_files = io.getFilesInDirectory(
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
            local info_file = io.readFileFromZip(mod_file, ".*?/info\\.json")
            local modInfo = json.decode(info_file)
            self.mods[#self.mods + 1] = {
                id = modInfo.name,
                game_id = self:GetGameID(),
                name = modInfo.title or modInfo.name,
                version = modInfo.version,
                description = modInfo.description or "",
                author = modInfo.author or "Unknown",
                file_path = mod_file,
                enabled = mod_list[modInfo.name] or false,
            }
        end
        if self.mods == nil or #self.mods == 0 then
            log.warn("No mods found in the mods directory " ..
                    self:GetGameModDirectory() .. tostring(#self.mods))
            return {}
        end
        return self.mods
    end,
    GetInstalledModByID = function(self, id)
        -- Example: return a mod by id
    end,
    GetGameModDirectory = function()
        -- This is usually located at:
        -- - Linux: ~/.factorio/mods/
        -- - Windows: C:\Users\<Username>\AppData\Roaming\Factorio\mods\
        -- - macOS: ~/Library/Application Support/factorio/mods/
        if os.is_windows then
            local appdata = os.getenv("APPDATA")
            if appdata then
                return appdata .. "\\Factorio\\mods\\"
            end
        elseif os.is_linux then
            return os.getenv("HOME") .. "/.factorio/mods/"
        elseif os.is_macos then
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
    end,
    GetMods = function(self)
        -- Convert self.mods.results to a table with fields id, name, version, enabled, game_id
        if self.mods == nil or self.mods.results == nil then
            return {}
        end
        local mods = {}
        for k, mod in pairs(self.mods.results) do
            mods[#mods + 1] = {
                id = mod.name,
                name = mod.name,
                -- latest_release.version
                version = mod.latest_release and mod.latest_release.version or "unknown",
                game_id = self:GetGameID(),
                author = mod.owner or "Unknown",
                game_versions = {
                    [mod.latest_release.info_json.factorio_version] = mod.latest_release.version,
                }
            }
        end
        return mods
    end,
    GetModByID = function(self, id)
    end,
}