package scripting

import (
	log "github.com/sirupsen/logrus"
	lua "github.com/yuin/gopher-lua"
)

func GetOrCreateTable(L *lua.LState, name string) *lua.LTable {
	if L == nil {
		log.Error("Lua state is nil, cannot extend table object")
		return nil
	}
	var table *lua.LTable
	if v := L.GetGlobal(name); v.Type() != lua.LTTable {
		log.Warnf("Global '%s' is not a table, creating a new one", name)
		table = L.NewTable()
	} else {
		table = v.(*lua.LTable)
	}

	return table
}

func luaPrint(L *lua.LState) int {
	logger := GetLogger(L)

	for i := 1; i <= L.GetTop(); i++ {
		if str, ok := L.Get(i).(lua.LString); ok {
			// Log with context lua
			logger.Info(str.String())
		} else {
			logger.Info(L.Get(i).Type().String(), ": ", L.Get(i))
		}
	}
	return 0
}

func luaErrorHandler(L *lua.LState) int {
	logger := GetLogger(L)
	if L.GetTop() >= 1 {
		luaErr := L.ToString(1)
		if luaErr != "" {
			logger.Error(luaErr)
		}
		return 0
	}
	logger.Error("Undefined error occurred")

	return 0
}
