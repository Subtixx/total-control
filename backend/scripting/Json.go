package scripting

import (
	"TotalControl/backend/utils"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	lua "github.com/yuin/gopher-lua"
)

func LuaRegisterJsonObject(l *lua.LState) {
	jsonTable := l.NewTable()
	jsonTable.RawSetString("encode", l.NewFunction(LuaJsonEncode))
	jsonTable.RawSetString("decode", l.NewFunction(LuaJsonDecode))
	l.SetGlobal("json", jsonTable)
}

func LuaJsonDecode(L *lua.LState) int {
	if L.GetTop() != 1 {
		L.Push(lua.LNil)
		log.Warn("LuaJsonDecode: Expected 1 argument, got none")
		return 1
	}

	jsonStr := L.CheckString(1)
	if jsonStr == "" {
		log.Warn("LuaJsonDecode: Empty JSON string provided")
		L.Push(lua.LNil)
		return 1
	}

	var result map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &result)
	if err != nil {
		log.Errorf("Failed to decode JSON: %v", err)
		L.Push(lua.LNil)
		return 1
	}
	L.Push(utils.MapToLuaTable(L, result))
	return 1
}

func LuaJsonEncode(L *lua.LState) int {
	table := L.ToTable(1)
	if table == nil {
		L.Push(lua.LNil)
		return 1
	}

	result := make(map[string]interface{})
	table.ForEach(func(key lua.LValue, value lua.LValue) {
		switch key.Type() {
		case lua.LTString:
			switch value.Type() {
			case lua.LTString:
				result[key.String()] = value.String()
			case lua.LTNumber:
				result[key.String()] = float64(value.(lua.LNumber))
			case lua.LTBool:
				result[key.String()] = bool(value.(lua.LBool))
			default:
				result[key.String()] = nil
			}
		default:
			L.RaiseError("Unsupported key type: %s", key.Type().String())
		}
	})

	jsonData, err := json.Marshal(result)
	if err != nil {
		L.Push(lua.LNil)
		return 1
	}
	L.Push(lua.LString(jsonData))
	return 1
}
