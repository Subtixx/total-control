package utils

import (
	"encoding/json"
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

func LuaJsonDecode(L *lua.LState) int {
	jsonStr := L.ToString(1)
	if jsonStr == "" {
		L.Push(lua.LNil)
		return 1
	}

	var result map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &result)
	if err != nil {
		L.Push(lua.LNil)
		return 1
	}
	L.Push(MapToLuaTable(L, result))
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
