package scripting

import (
	"TotalControl/backend/steam"
	"TotalControl/backend/utils"
	lua "github.com/yuin/gopher-lua"
)

func LuaGetAppDetails(L *lua.LState, appID string) int {
	httpClient := GetLuaHttpClient(L)
	if httpClient == nil {
		L.RaiseError("No HTTP client available")
		return 0
	}

	appDetails, err := steam.GetAppDetails(httpClient, appID)
	if err != nil {
		L.RaiseError("Error fetching app details: %v", err)
		return 0
	}
	appDetailsTable := utils.StructToLuaTable(L, appDetails)
	L.Push(appDetailsTable)
	return 1
}

func LuaGetLibraries(L *lua.LState) int {
	luaSteam := GetLuaSteam(L)
	if luaSteam == nil {
		L.RaiseError("No Steam instance available")
		return 0
	}

	librariesTable := L.NewTable()
	for _, library := range luaSteam.LibraryFolders {
		L.SetField(librariesTable, library.ContentId, library.ToLuaTable(L))
	}

	L.Push(librariesTable)
	return 1
}

func LuaGetAppSchema(L *lua.LState) int {
	appId := utils.GetString(L, 1)

	luaSteam := GetLuaSteam(L)
	if luaSteam == nil {
		L.RaiseError("No Steam instance available")
		return 0
	}

	appSchema, exists := luaSteam.AppSchemas[appId]
	if !exists {
		L.RaiseError("App schema not found for app ID: %s", appId)
		return 0
	}

	appSchemaTable := appSchema.ToLuaTable(L)
	L.Push(appSchemaTable)
	return 1
}

func LuaRegisterSteamObject(L *lua.LState) {
	steamTable := L.NewTable()
	L.SetField(steamTable, "GetAppDetails", L.NewFunction(func(L *lua.LState) int {
		appID := L.ToString(1)
		return LuaGetAppDetails(L, appID)
	}))
	L.SetField(steamTable, "GetAppSchema", L.NewFunction(LuaGetAppSchema))
	L.SetField(steamTable, "GetLibraries", L.NewFunction(LuaGetLibraries))
	L.SetGlobal("steam", steamTable)
}
