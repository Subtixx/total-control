package scripting

import (
	log "github.com/sirupsen/logrus"
	lua "github.com/yuin/gopher-lua"
)

// luaTableSize returns the size of a Lua table
func luaTableSize(L *lua.LState) int {
	table := L.CheckTable(1)
	size := 0
	table.ForEach(func(_ lua.LValue, _ lua.LValue) {
		size++
	})
	L.Push(lua.LNumber(size))
	return 1
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
