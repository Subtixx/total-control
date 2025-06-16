package scripting

import (
	"TotalControl/backend/plugins"
	"TotalControl/backend/steam"
	"TotalControl/backend/utils"
	lua "github.com/yuin/gopher-lua"
	"net/http"
)

func LuaGetAppDetails(L *lua.LState, appID string) int {
	var httpClient *http.Client
	luaPlugin := GetLuaPlugin(L)
	if luaPlugin != nil {
		if !luaPlugin.CanAccessNetwork() {
			L.Push(lua.LNil)
			L.Push(lua.LString(plugins.ErrNetworkAccessDenied))
			return 2
		}
		httpClient = luaPlugin.GetHttpClient()
	} else {
		httpClient = http.DefaultClient
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

func LuaRegisterSteamObject(L *lua.LState) {
	steamTable := L.NewTable()
	L.SetGlobal("Steam", steamTable)

	// Register the GetAppDetails function
	L.SetField(steamTable, "GetAppDetails", L.NewFunction(func(L *lua.LState) int {
		appID := L.ToString(1)
		return LuaGetAppDetails(L, appID)
	}))
}
