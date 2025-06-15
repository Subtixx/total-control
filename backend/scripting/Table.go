package scripting

// ExtendTableLib adds useful functions to the Lua table library.
import (
	"github.com/yuin/gopher-lua"
)

func LuaExtendTableObject(L *lua.LState) {
	table := GetOrCreateTable(L, "table")
	L.SetField(table, "keys", L.NewFunction(tableKeys))
	L.SetField(table, "values", L.NewFunction(tableValues))
	L.SetField(table, "length", L.NewFunction(tableLength))
}

// tableKeys returns a list of keys in the table.
func tableKeys(L *lua.LState) int {
	tbl := L.CheckTable(1)
	keys := L.NewTable()
	tbl.ForEach(func(key, _ lua.LValue) {
		keys.Append(key)
	})
	L.Push(keys)
	return 1
}

// tableValues returns a list of values in the table.
func tableValues(L *lua.LState) int {
	tbl := L.CheckTable(1)
	values := L.NewTable()
	tbl.ForEach(func(_, value lua.LValue) {
		values.Append(value)
	})
	L.Push(values)
	return 1
}

// tableLength returns the number of elements in the table.
func tableLength(L *lua.LState) int {
	tbl := L.CheckTable(1)
	var count int
	tbl.ForEach(func(_, _ lua.LValue) {
		count++
	})
	L.Push(lua.LNumber(count))
	return 1
}
