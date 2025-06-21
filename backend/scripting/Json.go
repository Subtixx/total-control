package scripting

import (
	"TotalControl/backend/utils"
	"encoding/json"
	lua "github.com/yuin/gopher-lua"
)

func LuaRegisterJsonObject(l *lua.LState) {
	jsonTable := l.NewTable()
	jsonTable.RawSetString("encode", l.NewFunction(LuaJsonEncode))
	jsonTable.RawSetString("decode", l.NewFunction(LuaJsonDecode))
	l.SetGlobal("json", jsonTable)
}

func LuaJsonDecode(L *lua.LState) int {
	jsonStr := utils.GetString(L, 1)
	if jsonStr == "" {
		L.RaiseError("LuaJsonDecode: JSON string cannot be empty")
		return 0
	}

	var result map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &result)
	if err != nil {
		L.RaiseError(err.Error())
		return 0
	}

	L.Push(utils.MapToLuaTable(L, result))
	return 1
}

func LuaJsonEncode(L *lua.LState) int {
	table := utils.GetTable(L, 1)
	if table == nil {
		L.RaiseError("LuaJsonEncode: Expected table argument")
		return 0
	}

	result := make(map[string]interface{})
	table.ForEach(func(key lua.LValue, value lua.LValue) {
		result[key.String()] = utils.FromLuaValue(L, value)
	})

	jsonData, err := json.Marshal(result)
	if err != nil {
		L.RaiseError(err.Error())
		return 0
	}

	L.Push(lua.LString(jsonData))
	return 1
}
