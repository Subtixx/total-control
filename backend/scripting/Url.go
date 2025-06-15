package scripting

import (
	"TotalControl/backend/utils"
	lua "github.com/yuin/gopher-lua"
	"net/url"
	"strings"
)

func UrlToLuaUserData(L *lua.LState, u *url.URL) *lua.LUserData {
	if u == nil {
		return nil
	}

	ud := L.NewUserData()
	ud.Value = u
	mt := L.GetTypeMetatable("Url")
	L.SetMetatable(ud, mt)
	L.SetField(mt, "__tostring", L.NewFunction(func(L *lua.LState) int {
		ud := L.CheckUserData(1)
		if v, ok := ud.Value.(*url.URL); ok {
			L.Push(lua.LString(v.String()))
		} else {
			L.Push(lua.LNil)
		}
		return 1
	}))
	return ud
}

func registerUrlType(l *lua.LState) {
	mt := l.NewTypeMetatable("Url")
	l.SetGlobal("Url", mt)
	l.SetField(mt, "__index", l.NewFunction(func(L *lua.LState) int {
		ud := L.CheckUserData(1)
		udUrl, ok := ud.Value.(*url.URL)
		if !ok {
			L.Push(lua.LNil)
			return 1
		}
		key := L.CheckString(2)
		pushLuaValueFromUrl(L, key, udUrl)
		return 1
	}))
}

func luaTableFromUserInfo(L *lua.LState, userInfo *url.Userinfo) *lua.LTable {
	userTable := L.NewTable()
	if userInfo != nil {
		userTable.RawSetString("username", lua.LString(userInfo.Username()))
		if password, ok := userInfo.Password(); ok {
			userTable.RawSetString("password", lua.LString(password))
		} else {
			userTable.RawSetString("password", lua.LNil) // Handle no password
		}
	} else {
		userTable.RawSetString("username", lua.LNil) // Handle no user
		userTable.RawSetString("password", lua.LNil)
	}
	return userTable
}

func luaTableFromQueryValues(L *lua.LState, queryValues url.Values) *lua.LTable {
	queryTable := L.NewTable()
	for k, v := range queryValues {
		if len(v) > 0 {
			queryTable.RawSetString(k, lua.LString(strings.Join(v, ",")))
		} else {
			queryTable.RawSetString(k, lua.LNil)
		}
	}
	return queryTable
}

func pushLuaValueFromUrl(L *lua.LState, key string, udUrl *url.URL) {
	switch key {
	case "scheme":
		L.Push(utils.ToLuaValue(L, udUrl.Scheme))
	case "opaque":
		L.Push(utils.ToLuaValue(L, udUrl.Opaque))
	case "user":
		L.Push(luaTableFromUserInfo(L, udUrl.User))
	case "host":
		L.Push(utils.ToLuaValue(L, udUrl.Host))
	case "path":
		L.Push(utils.ToLuaValue(L, udUrl.Path))
	case "query":
		L.Push(luaTableFromQueryValues(L, udUrl.Query()))
	case "fragment":
		L.Push(utils.ToLuaValue(L, udUrl.Fragment))
	case "rawQuery":
		L.Push(utils.ToLuaValue(L, udUrl.RawQuery))
	case "rawPath":
		L.Push(utils.ToLuaValue(L, udUrl.RawPath))
	default:
		L.Push(lua.LNil)
	}
}
