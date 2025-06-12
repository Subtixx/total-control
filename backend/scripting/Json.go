package scripting

import (
	"TotalControl/backend/utils"
	lua "github.com/yuin/gopher-lua"
)

func luaRegisterJsonObject(l *lua.LState) {
	jsonTable := l.NewTable()
	jsonTable.RawSetString("encode", l.NewFunction(utils.LuaJsonEncode))
	jsonTable.RawSetString("decode", l.NewFunction(utils.LuaJsonDecode))
	l.SetGlobal("json", jsonTable)
}
