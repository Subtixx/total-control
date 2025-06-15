package scripting

import (
	"TotalControl/backend/steam"
	"TotalControl/backend/utils"
	lua "github.com/yuin/gopher-lua"
)

func LuaGetAppDetails(L *lua.LState, appID string) int {
	appDetails, err := steam.GetAppDetails(appID)
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
