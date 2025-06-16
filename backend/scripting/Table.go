package scripting

// ExtendTableLib adds useful functions to the Lua table library.
import (
	"github.com/yuin/gopher-lua"
)

func LuaExtendTableObject(L *lua.LState) {
	table := GetOrCreateTable(L, "table")
	L.SetField(table, "keys", L.NewFunction(luaTableKeys))
	L.SetField(table, "values", L.NewFunction(luaTableValues))
	L.SetField(table, "length", L.NewFunction(luaTableLength))
}

// luaTableKeys returns a list of keys in the table. The order of keys is not guaranteed.
func luaTableKeys(L *lua.LState) int {
	tbl := L.CheckTable(1)
	keys := L.NewTable()
	tbl.ForEach(func(key, _ lua.LValue) {
		keys.Append(key)
	})
	L.Push(keys)
	return 1
}

// luaTableValues returns a list of values in the table. The order of values is not guaranteed.
func luaTableValues(L *lua.LState) int {
	tbl := L.CheckTable(1)
	values := L.NewTable()
	tbl.ForEach(func(_, value lua.LValue) {
		values.Append(value)
	})
	L.Push(values)
	return 1
}

// luaTableLength returns the number of elements in the table.
func luaTableLength(L *lua.LState) int {
	tbl := L.CheckTable(1)
	var count int
	tbl.ForEach(func(_, _ lua.LValue) {
		count++
	})
	L.Push(lua.LNumber(count))
	return 1
}

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
