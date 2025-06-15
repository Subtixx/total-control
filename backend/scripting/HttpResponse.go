package scripting

import (
	"TotalControl/backend/utils"
	lua "github.com/yuin/gopher-lua"
	"io"
	"net/http"
	"strconv"
)

func registerHttpResponseType(l *lua.LState) {
	mt := l.NewTypeMetatable("HttpResponse")
	l.SetGlobal("HttpResponse", mt)
	l.SetField(mt, "__index", l.NewFunction(func(L *lua.LState) int {
		ud := L.CheckUserData(1)
		key := L.CheckString(2)
		if v, ok := ud.Value.(*http.Response); ok {
			L.Push(luaValueFromResponse(L, key, v))
			return 1
		}
		L.Push(lua.LNil)
		return 1
	}))
}

func httpResponseToString(resp *http.Response) string {
	if resp == nil {
		return "HttpResponse(nil)"
	}
	return "HttpResponse(Status: " + resp.Status + ", StatusCode: " + strconv.Itoa(resp.StatusCode) + ")"
}

func newHttpResponseUserData(L *lua.LState, resp *http.Response) *lua.LUserData {
	ud := L.NewUserData()
	ud.Value = resp
	mt := L.GetTypeMetatable("HttpResponse")
	L.SetMetatable(ud, mt)
	L.SetField(mt, "__tostring", L.NewFunction(func(L *lua.LState) int {
		ud := L.CheckUserData(1)
		luaResp, ok := ud.Value.(*http.Response)
		if !ok {
			L.Push(lua.LNil)
			return 1
		}
		L.Push(utils.ToLuaValue(L, httpResponseToString(luaResp)))
		return 1
	}))
	return ud
}

func luaValueFromResponse(L *lua.LState, key string, resp *http.Response) lua.LValue {
	switch key {
	case "status":
		return utils.ToLuaValue(L, resp.Status)
	case "status_code":
		return utils.ToLuaValue(L, resp.StatusCode)
	case "proto":
		return utils.ToLuaValue(L, resp.Proto)
	case "proto_major":
		return utils.ToLuaValue(L, resp.ProtoMajor)
	case "proto_minor":
		return utils.ToLuaValue(L, resp.ProtoMinor)
	case "headers":
		return luaTableFromHeaders(L, resp.Header)
	case "body":
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return utils.ToLuaValue(L, err.Error())
		}
		return utils.ToLuaValue(L, string(bodyBytes))
	case "content_length":
		return utils.ToLuaValue(L, resp.ContentLength)
	case "uncompressed":
		return utils.ToLuaValue(L, resp.Uncompressed)
	}
	return lua.LNil
}
