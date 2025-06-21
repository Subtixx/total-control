package scripting

import (
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
	httpClient := GetLuaHttpClient(L)
	if httpClient == nil {
		L.RaiseError("HTTP client not found in context")
		return 0
	}

	if L.GetTop() != 1 {
		L.RaiseError("Expected 1 argument: URL string/HttpRequest")
		return 0
	}

	var httpReq *http.Request
	var err error
	if L.Get(1).Type() == lua.LTString {
		httpReq, err = httpRequestFromUrl(L)
		if err != nil {
			L.RaiseError(err.Error())
			return 0
		}
	} else if L.Get(1).Type() == lua.LTUserData {
		httpReq, err = httpRequestFromStack(L)
		if err != nil {
			L.RaiseError(err.Error())
			return 0
		}
	} else {
		L.RaiseError("Expected HttpRequest userdata or URL string")
		return 0
	}

	httpReq.Method = "GET"
	httpResponse, err := httpClient.Do(httpReq)
	if err != nil {
		log.Errorf("HTTP GET request failed: %v", err)
		L.RaiseError(err.Error())
		return 0
	}

	L.Push(newHttpResponseUserData(L, httpResponse))
	return 1
}

func luaHttpPost(L *lua.LState) int {
	httpClient := GetLuaHttpClient(L)
	if httpClient == nil {
		L.RaiseError("HTTP client not found in context")
		return 0
	}

	if L.GetTop() < 1 {
		L.RaiseError("Expected 1 argument: URL string/HttpRequest")
		return 0
	}

	var httpReq *http.Request
	var err error
	if L.Get(1).Type() == lua.LTString {
		urlStr := L.CheckString(1)
		// Body table
		bodyTable := L.Get(2)
		if bodyTable.Type() != lua.LTTable {
			L.RaiseError("Expected body to be a table")
			return 0
		}
		bodyData := utils.LuaTableToMap(L, bodyTable.(*lua.LTable))
		bodyString, err := json.Marshal(bodyData)
		if err != nil {
			L.RaiseError(err.Error())
			return 0
		}
		reader := strings.NewReader(string(bodyString))

		httpReq, err = http.NewRequest("POST", urlStr, reader)
		if err != nil {
			L.RaiseError(err.Error())
			return 0
		}
	} else if L.Get(1).Type() == lua.LTUserData {
		httpReq, err = httpRequestFromStack(L)
		if err != nil {
			L.RaiseError(err.Error())
			return 0
		}
	} else {
		L.RaiseError("Expected HttpRequest userdata or URL string")
		return 0
	}

	httpReq.Method = "POST"
	httpResponse, err := httpClient.Do(httpReq)
	if err != nil {
		log.Errorf("HTTP POST request failed: %v", err)
		L.RaiseError(err.Error())
		return 0
	}

	L.Push(newHttpResponseUserData(L, httpResponse))
	err = httpResponse.Body.Close()
	if err != nil {
		log.Errorf("Failed to close response body: %v", err)
	}
	return 1
}

func luaHttpDownloadFile(L *lua.LState) int {
	httpClient := GetLuaHttpClient(L)
	if httpClient == nil {
		L.RaiseError("HTTP client not found in context")
		return 0
	}

	if L.GetTop() != 2 {
		L.RaiseError("Expected 2 arguments: URL string/HttpRequest, destination file path")
		return 0
	}

	var httpReq *http.Request
	if L.Get(1).Type() == lua.LTString {
		urlStr := L.CheckString(1)
		httpReq, _ = http.NewRequest("GET", urlStr, http.NoBody)
	} else if L.Get(1).Type() == lua.LTUserData {
		var err error
		httpReq, err = httpRequestFromStack(L)
		if err != nil {
			L.RaiseError(err.Error())
			return 0
		}
	} else {
		L.RaiseError("Expected HttpRequest userdata or URL string")
		return 0
	}

	destPath := L.CheckString(2)
	if destPath == "" {
		L.RaiseError("Destination path cannot be empty")
		return 0
	}

	resp, err := httpClient.Do(httpReq)
	if err != nil {
		log.Errorf("HTTP request failed: %v", err)
		L.RaiseError(err.Error())
		return 0
	}

	if resp.StatusCode != http.StatusOK {
		L.RaiseError(fmt.Sprintf("HTTP request failed with status code %d", resp.StatusCode))
		return 0
	}

	file, err := os.Create(destPath)
	if err != nil {
		L.RaiseError(fmt.Sprintf("Failed to create file: %v", err))
		return 0
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

		L.RaiseError(fmt.Sprintf("Failed to write file: %v", err))
		return 0
	}

	if err := resp.Body.Close(); err != nil {
		log.Errorf("Failed to close response body: %v", err)
	}

	L.Push(lua.LTrue)
	return 1
}
