package utils

import (
	lua "github.com/yuin/gopher-lua"
	"strings"
)

func LuaArgsToString(args []lua.LValue) string {
	var sb strings.Builder
	for i, arg := range args {
		if i > 0 {
			sb.WriteString(", ")
		}
		switch arg.Type() {
		case lua.LTNil:
			sb.WriteString("nil")
		case lua.LTBool:
			if arg.(lua.LBool) {
				sb.WriteString("true")
			} else {
				sb.WriteString("false")
			}
		case lua.LTNumber:
			sb.WriteString(arg.String())
		case lua.LTString:
			sb.WriteString(`"` + arg.String() + `"`)
		default:
			sb.WriteString(arg.String())
		}
	}
	return sb.String()
}

func FromLuaValue(L *lua.LState, value lua.LValue) interface{} {
	switch v := value.(type) {
	case lua.LString:
		return v.String()
	case lua.LNumber:
		return float64(v)
	case lua.LBool:
		return bool(v)
	case *lua.LTable:
		return LuaTableToMap(L, v)
	default:
		L.RaiseError("Unsupported Lua value type: %s", v.Type().String())
		return nil
	}
}

func ToLuaValue(L *lua.LState, value interface{}) lua.LValue {
	switch v := value.(type) {
	case nil:
		return lua.LNil
	case string:
		return lua.LString(v)
	case float64:
		return lua.LNumber(v)
	case float32:
		return lua.LNumber(v)
	case int:
		return lua.LNumber(v)
	case int64:
		return lua.LNumber(v)
	case int32:
		return lua.LNumber(v)
	case int16:
		return lua.LNumber(v)
	case int8:
		return lua.LNumber(v)
	case uint:
		return lua.LNumber(v)
	case uint64:
		return lua.LNumber(v)
	case uint32:
		return lua.LNumber(v)
	case uint16:
		return lua.LNumber(v)
	case uint8:
		return lua.LNumber(v)
	case bool:
		return lua.LBool(v)
	case map[string]interface{}:
		tbl := L.CreateTable(len(v), 0)
		for key, val := range v {
			tbl.RawSetString(key, ToLuaValue(L, val))
		}
		return tbl
	case []interface{}:
		tbl := L.CreateTable(len(v), 0)
		for i, val := range v {
			tbl.RawSetInt(i+1, ToLuaValue(L, val))
		}
		return tbl
	default:
		return lua.LNil
	}
}

func MapToLuaTable(L *lua.LState, m map[string]interface{}) *lua.LTable {
	tbl := L.CreateTable(len(m), 0)
	for key, value := range m {
		switch v := value.(type) {
		case string:
			tbl.RawSetString(key, lua.LString(v))
		case float64:
			tbl.RawSetString(key, lua.LNumber(v))
		case bool:
			tbl.RawSetString(key, lua.LBool(v))
		case nil:
			tbl.RawSetString(key, lua.LNil)
		default:
			tbl.RawSetString(key, ToLuaValue(L, v))
		}
	}
	return tbl
}

func MultiMapToLuaTable(L *lua.LState, m map[string][]interface{}) *lua.LTable {
	tbl := L.CreateTable(len(m), 0)
	for key, values := range m {
		arr := L.CreateTable(len(values), 0)
		for i, value := range values {
			arr.RawSetInt(i+1, ToLuaValue(L, value))
		}
		tbl.RawSetString(key, arr)
	}
	return tbl
}

func LuaTableToMap(L *lua.LState, tbl *lua.LTable) map[string]interface{} {
	result := make(map[string]interface{})
	// Use FromLuaValue to convert Lua values to Go values
	tbl.ForEach(func(key lua.LValue, value lua.LValue) {
		result[key.String()] = FromLuaValue(L, value)
	})
	return result
}

func LuaTableToMultiMap(L *lua.LState, tbl *lua.LTable) map[string][]interface{} {
	result := make(map[string][]interface{})
	// Use FromLuaValue to convert Lua values to Go values
	tbl.ForEach(func(key lua.LValue, value lua.LValue) {
		if arr, ok := value.(*lua.LTable); ok {
			var values []interface{}
			arr.ForEach(func(i lua.LValue, v lua.LValue) {
				values = append(values, FromLuaValue(L, v))
			})
			result[key.String()] = values
		} else {
			result[key.String()] = []interface{}{FromLuaValue(L, value)}
		}
	})
	return result
}
