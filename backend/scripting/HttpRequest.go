package scripting

import (
	"TotalControl/backend/utils"
	"fmt"
	lua "github.com/yuin/gopher-lua"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func registerHttpRequestType(l *lua.LState) {
	mt := l.NewTypeMetatable("HttpRequest")
	l.SetGlobal("HttpRequest", mt)
	l.SetField(mt, "__index", l.NewFunction(func(L *lua.LState) int {
		ud := L.CheckUserData(1)
		key := L.CheckString(2)
		if v, ok := ud.Value.(*http.Request); ok {
			L.Push(luaValueFromRequest(L, key, v))
			return 1
		}
		L.Push(lua.LNil)
		return 1
	}))
}

func httpRequestToString(req *http.Request) string {
	if req == nil {
		return "HttpRequest(nil)"
	}
	return fmt.Sprintf("HttpRequest(Method: %s, URL: %s, Headers: %v)", req.Method, req.URL.String(), req.Header)
}

func newHttpRequestUserData(L *lua.LState, req *http.Request) *lua.LUserData {
	ud := L.NewUserData()
	ud.Value = req
	mt := L.GetTypeMetatable("HttpRequest")
	L.SetMetatable(ud, mt)
	L.SetField(mt, "__tostring", L.NewFunction(func(L *lua.LState) int {
		ud := L.CheckUserData(1)
		luaReq, ok := ud.Value.(*http.Request)
		if !ok {
			L.Push(lua.LNil)
			return 1
		}
		L.Push(utils.ToLuaValue(L, httpRequestToString(luaReq)))
		return 1
	}))
	return ud
}

func httpRequestFromUrl(L *lua.LState) (*http.Request, error) {
	urlStr := L.CheckString(1)
	req, err := http.NewRequest("GET", urlStr, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	// Set the User-Agent header to mimic a browser
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; TotalControl/1.0; +https://github.com/subtixx/total-control)")
	return req, nil
}

func httpRequestFromStack(L *lua.LState) (*http.Request, error) {
	var httpReq *http.Request
	valueType := L.Get(1).Type()
	if valueType != lua.LTUserData {
		return nil, fmt.Errorf("expected HttpRequest userdata, got %s", valueType.String())
	}

	var ok bool
	httpReq, ok = L.CheckUserData(1).Value.(*http.Request)
	if !ok {
		return nil, fmt.Errorf("expected HttpRequest userdata, got %s", valueType.String())
	}
	return httpReq, nil
}

func luaValueFromBody(L *lua.LState, reader io.Reader) lua.LValue {
	if reader == nil {
		return lua.LNil
	}
	bodyBytes, err := io.ReadAll(reader)
	if err != nil {
		return lua.LNil
	}
	return lua.LString(bodyBytes)
}

func luaTableFromHeaders(L *lua.LState, headers http.Header) *lua.LTable {
	headersTable := L.NewTable()
	for k, values := range headers {
		if len(values) > 0 {
			headersTable.RawSetString(k, lua.LString(strings.Join(values, ", ")))
		} else {
			headersTable.RawSetString(k, lua.LNil)
		}
	}
	return headersTable
}

func luaTableFromUrlValues(L *lua.LState, values *url.Values) *lua.LTable {
	valuesTable := L.NewTable()
	if values == nil {
		return valuesTable
	}

	for k, v := range *values {
		if len(v) > 0 {
			valuesTable.RawSetString(k, lua.LString(strings.Join(v, ", ")))
		} else {
			valuesTable.RawSetString(k, lua.LNil)
		}
	}
	return valuesTable
}

func luaValueFromRequest(L *lua.LState, key string, v *http.Request) lua.LValue {
	switch key {
	case "url":
		return UrlToLuaUserData(L, v.URL)
	case "method":
		return lua.LString(v.Method)
	case "headers":
		return luaTableFromHeaders(L, v.Header)
	case "body":
		return luaValueFromBody(L, v.Body)
	case "proto":
		return utils.ToLuaValue(L, v.Proto)
	case "proto_major":
		return utils.ToLuaValue(L, v.ProtoMajor)
	case "proto_minor":
		return utils.ToLuaValue(L, v.ProtoMinor)
	case "content_length":
		return utils.ToLuaValue(L, v.ContentLength)
	case "host":
		return utils.ToLuaValue(L, v.Host)
	case "remote_addr":
		return utils.ToLuaValue(L, v.RemoteAddr)
	case "request_uri":
		return utils.ToLuaValue(L, v.RequestURI)
	case "form":
		return luaTableFromUrlValues(L, &v.Form)
	case "post_form":
		return luaTableFromUrlValues(L, &v.PostForm)
	case "multipart_form":
		if v.MultipartForm != nil {
			multipartTable := L.NewTable()
			for key, values := range v.MultipartForm.Value {
				if len(values) > 0 {
					multipartTable.RawSetString(key, lua.LString(strings.Join(values, ", ")))
				} else {
					multipartTable.RawSetString(key, lua.LNil)
				}
			}
			return multipartTable
		}
	}
	return lua.LNil
}
