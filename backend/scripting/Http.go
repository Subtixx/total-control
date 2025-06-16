package scripting

import (
	"TotalControl/backend/plugins"
	"TotalControl/backend/utils"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/yuin/gopher-lua"
	"io"
	"net/http"
	"os"
	"strings"
)

func LuaRegisterHttpObject(L *lua.LState) {
	httpTable := L.NewTable()
	httpTable.RawSetString("get", L.NewFunction(luaHttpGet))
	httpTable.RawSetString("post", L.NewFunction(luaHttpPost))
	httpTable.RawSetString("download", L.NewFunction(luaHttpDownloadFile))
	L.SetGlobal("http", httpTable)

	registerUrlType(L)
	registerHttpRequestType(L)
	registerHttpResponseType(L)
}

func luaHttpGet(L *lua.LState) int {
	if !CheckCan(L, plugins.CapabilityNetwork) {
		L.Push(lua.LFalse)
		L.Push(lua.LString(plugins.ErrNetworkAccessDenied))
		return 2
	}
	httpClient := GetLuaHttpClient(L)

	if L.GetTop() != 1 {
		L.Push(lua.LFalse)
		L.Push(lua.LString("Expected 1 argument: HttpRequest userdata or URL string"))
		return 2
	}

	var httpReq *http.Request
	var err error
	if L.Get(1).Type() == lua.LTString {
		httpReq, err = httpRequestFromUrl(L)
		if err != nil {
			L.Push(lua.LFalse)
			L.Push(lua.LString(err.Error()))
			return 2
		}
	} else if L.Get(1).Type() == lua.LTUserData {
		httpReq, err = httpRequestFromStack(L)
		if err != nil {
			L.Push(lua.LFalse)
			L.Push(lua.LString(err.Error()))
			return 2
		}
	} else {
		L.Push(lua.LFalse)
		L.Push(lua.LString("Expected HttpRequest userdata or URL string"))
		return 2
	}

	httpReq.Method = "GET"
	httpResponse, err := httpClient.Do(httpReq)
	if err != nil {
		log.Errorf("HTTP GET request failed: %v", err)
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(newHttpResponseUserData(L, httpResponse))
	return 1
}

func luaHttpPost(L *lua.LState) int {
	if !CheckCan(L, plugins.CapabilityNetwork) {
		L.Push(lua.LFalse)
		L.Push(lua.LString(plugins.ErrNetworkAccessDenied))
		return 2
	}
	httpClient := GetLuaHttpClient(L)

	if L.GetTop() < 1 {
		L.Push(lua.LFalse)
		L.Push(lua.LString("Expected at least 1 argument: URL string/HttpRequest and optional body table"))
		return 2
	}
	var httpReq *http.Request
	var err error
	if L.Get(1).Type() == lua.LTString {
		urlStr := L.CheckString(1)
		// Body table
		bodyTable := L.Get(2)
		if bodyTable.Type() != lua.LTTable {
			L.Push(lua.LFalse)
			L.Push(lua.LString("Expected second argument to be a table for POST body"))
			return 2
		}
		bodyData := utils.LuaTableToMap(L, bodyTable.(*lua.LTable))
		bodyString, err := json.Marshal(bodyData)
		if err != nil {
			L.Push(lua.LFalse)
			L.Push(lua.LString(fmt.Sprintf("Failed to marshal body data: %v", err)))
			return 2
		}
		reader := strings.NewReader(string(bodyString))

		httpReq, err = http.NewRequest("POST", urlStr, reader)
		if err != nil {
			L.Push(lua.LFalse)
			L.Push(lua.LString(err.Error()))
			return 2
		}
	} else if L.Get(1).Type() == lua.LTUserData {
		httpReq, err = httpRequestFromStack(L)
		if err != nil {
			L.Push(lua.LFalse)
			L.Push(lua.LString(err.Error()))
			return 2
		}
	} else {
		L.Push(lua.LFalse)
		L.Push(lua.LString("Expected HttpRequest userdata"))
		return 2
	}

	httpReq.Method = "POST"
	httpResponse, err := httpClient.Do(httpReq)
	if err != nil {
		log.Errorf("HTTP POST request failed: %v", err)
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(newHttpResponseUserData(L, httpResponse))
	err = httpResponse.Body.Close()
	if err != nil {
		log.Errorf("Failed to close response body: %v", err)
	}
	return 1
}

func luaHttpDownloadFile(L *lua.LState) int {
	if !CheckCan(L, plugins.CapabilityNetwork) {
		L.Push(lua.LFalse)
		L.Push(lua.LString(plugins.ErrNetworkAccessDenied))
		return 2
	}
	httpClient := GetLuaHttpClient(L)

	if L.GetTop() != 2 {
		L.Push(lua.LFalse)
		L.Push(lua.LString("Expected 2 arguments: URL string/HttpRequest and destination file path"))
		return 2
	}
	var httpReq *http.Request
	if L.Get(1).Type() == lua.LTString {
		urlStr := L.CheckString(1)
		httpReq, _ = http.NewRequest("GET", urlStr, http.NoBody)
	} else if L.Get(1).Type() == lua.LTUserData {
		var err error
		httpReq, err = httpRequestFromStack(L)
		if err != nil {
			L.Push(lua.LFalse)
			L.Push(lua.LString(err.Error()))
			return 2
		}
	} else {
		L.Push(lua.LFalse)
		L.Push(lua.LString("Expected URL string or HttpRequest userdata"))
		return 2
	}

	destPath := L.CheckString(2)
	if destPath == "" {
		L.Push(lua.LString("Destination file path cannot be empty"))
		return 2
	}

	resp, err := httpClient.Do(httpReq)
	if err != nil {
		log.Errorf("HTTP request failed: %v", err)
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	if resp.StatusCode != http.StatusOK {
		L.Push(lua.LFalse)
		L.Push(lua.LString(fmt.Sprintf("HTTP request failed with status code: %d", resp.StatusCode)))
		return 2
	}
	file, err := os.Create(destPath)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(fmt.Sprintf("Failed to create file: %v", err)))
		return 2
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Errorf("Failed to close file: %v", err)
		}
	}(file)
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		err := resp.Body.Close()
		if err != nil {
			log.Errorf("Failed to close response body: %v", err)
		}

		L.Push(lua.LFalse)
		L.Push(lua.LString(fmt.Sprintf("Failed to write to file: %v", err)))
		return 2
	}

	if err := resp.Body.Close(); err != nil {
		log.Errorf("Failed to close response body: %v", err)
	}

	L.Push(lua.LTrue)
	L.Push(lua.LNil)
	return 2
}
