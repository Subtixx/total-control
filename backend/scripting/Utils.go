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
	// Print function that captures Lua print calls
	for i := 1; i <= L.GetTop(); i++ {
		if str, ok := L.Get(i).(lua.LString); ok {
			// Log with context lua
			log.WithField("lua", true).Info(str.String())
		} else {
			log.WithField("lua", true).Info(L.Get(i).Type().String(), ": ", L.Get(i))
		}
	}
	return 0
}

func luaErrorHandler(L *lua.LState) int {
	// Custom error handler for Lua
	if err := L.ToString(1); err != "" {
		println("Lua Error:", L.ToString(1))
	}
	return 0
}
