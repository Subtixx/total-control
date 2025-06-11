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
	"path/filepath"
	"strings"
	"time"
)

type HttpResponse struct {
	StatusCode int
	Headers    http.Header
	Body       string
	JsonBody   map[string]interface{}
}

type HttpRequest struct {
	URL     string
	Method  string
	Headers map[string]string
	Body    string
}

func newHttpResponseUserData(L *lua.LState, resp *HttpResponse) *lua.LUserData {
	ud := L.NewUserData()
	ud.Value = resp
	mt := L.GetTypeMetatable("HttpResponse")
	L.SetMetatable(ud, mt)
	L.SetField(mt, "__tostring", L.NewFunction(func(L *lua.LState) int {
		ud := L.CheckUserData(1)
		if v, ok := ud.Value.(*HttpResponse); ok {
			// Create a string representation of the HttpResponse
			var sb strings.Builder
			sb.WriteString(fmt.Sprintf("HttpResponse(StatusCode: %d, Headers: {", v.StatusCode))
			for k, values := range v.Headers {
				sb.WriteString(fmt.Sprintf("%s: %s, ", k, strings.Join(values, ", ")))
			}
			sb.WriteString("}, Body: ")
			sb.WriteString(v.Body)
			if v.JsonBody != nil {
				sb.WriteString(", JsonBody: {")
				first := true
				for k, val := range v.JsonBody {
					if !first {
						sb.WriteString(", ")
					}
					first = false
					sb.WriteString(fmt.Sprintf("%s: %v", k, val))
				}
				sb.WriteString("}")
			}
			sb.WriteString(")")
			L.Push(lua.LString(sb.String()))
		} else {
			L.Push(lua.LNil) // Return nil if not a HttpResponse
		}
		return 1
	}))
	return ud
}

func newHttpRequestUserData(L *lua.LState, req *HttpRequest) *lua.LUserData {
	ud := L.NewUserData()
	ud.Value = req
	mt := L.GetTypeMetatable("HttpRequest")
	L.SetMetatable(ud, mt)
	L.SetField(mt, "__tostring", L.NewFunction(func(L *lua.LState) int {
		ud := L.CheckUserData(1)
		if v, ok := ud.Value.(*HttpRequest); ok {
			L.Push(lua.LString(fmt.Sprintf("HttpRequest(URL: %s, Method: %s, Headers: %v, Body: %s)", v.URL, v.Method, v.Headers, v.Body)))
		} else {
			L.Push(lua.LNil) // Return nil if not a HttpRequest
		}
		return 1
	}))
	return ud
}

func luaRegisterHttpObject(L *lua.LState) {
	httpTable := L.NewTable()
	httpTable.RawSetString("get", L.NewFunction(luaHttpGet))
	httpTable.RawSetString("post", L.NewFunction(luaHttpPost))
	httpTable.RawSetString("downloadFile", L.NewFunction(luaHttpDownloadFile))
	L.SetGlobal("http", httpTable)

	registerHttpRequestType(L)
	registerHttpResponseType(L)
}

func registerHttpRequestType(l *lua.LState) {
	mt := l.NewTypeMetatable("HttpRequest")
	l.SetGlobal("HttpRequest", mt)
	l.SetField(mt, "__index", l.NewFunction(func(L *lua.LState) int {
		ud := L.CheckUserData(1)
		key := L.CheckString(2)
		if v, ok := ud.Value.(*HttpRequest); ok {
			switch key {
			case "url":
				L.Push(lua.LString(v.URL))
			case "method":
				L.Push(lua.LString(v.Method))
			case "headers":
				headersTable := L.NewTable()
				for k, v := range v.Headers {
					headersTable.RawSetString(k, lua.LString(v))
				}
				L.Push(headersTable)
			case "body":
				L.Push(lua.LString(v.Body))
			default:
				L.Push(lua.LNil) // Return nil for unknown fields
			}
			return 1
		}
		L.Push(lua.LNil) // Return nil if userdata is invalid
		return 1
	}))
}

func registerHttpResponseType(l *lua.LState) {
	mt := l.NewTypeMetatable("HttpResponse")
	l.SetGlobal("HttpResponse", mt)
	l.SetField(mt, "__index", l.NewFunction(func(L *lua.LState) int {
		ud := L.CheckUserData(1)
		key := L.CheckString(2)
		if v, ok := ud.Value.(*HttpResponse); ok {
			switch key {
			case "status_code":
				L.Push(lua.LNumber(v.StatusCode))
			case "headers":
				headersTable := L.NewTable()
				for k, values := range v.Headers {
					headersTable.RawSetString(k, lua.LString(strings.Join(values, ", ")))
				}
				L.Push(headersTable)
			case "body":
				if v.JsonBody != nil {
					L.Push(utils.MapToLuaTable(L, v.JsonBody)) // Convert JsonBody to Lua table
				} else {
					L.Push(lua.LString(v.Body))
				}
			default:
				L.Push(lua.LNil) // Return nil for unknown fields
			}
			return 1
		}
		L.Push(lua.LNil) // Return nil if userdata is invalid
		return 1
	}))
	/*l.SetField(mt, "__index", l.SetFuncs(l.NewTable(), map[string]lua.LGFunction{
		"getStatusCode": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			if v, ok := ud.Value.(*HttpResponse); ok {
				L.Push(lua.LNumber(v.StatusCode))
				return 1
			}
			return 0 // Return nil if not a HttpResponse
		},
		"getHeaders": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			if v, ok := ud.Value.(*HttpResponse); ok {
				headersTable := l.NewTable()
				for k, v := range v.Headers {
					// Join header values with commas
					headersTable.RawSetString(k, lua.LString(strings.Join(v, ", ")))
				}
				return 1
			}
			return 0 // Return nil if not a HttpResponse
		},
		"getBody": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			if v, ok := ud.Value.(*HttpResponse); ok {
				L.Push(lua.LString(v.Body))
				return 1
			}
			return 0 // Return nil if not a HttpResponse
		},
		"getJsonBody": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			if v, ok := ud.Value.(*HttpResponse); ok {
				if v.JsonBody != nil {
					L.Push(v.JsonBody) // Push the Lua table if it exists
					return 1
				}
				L.Push(lua.LNil) // Return nil if JsonBody is not set
				return 1
			}
			return 0 // Return nil if not a HttpResponse
		},
	}))*/
}

func getHttpRequestFromStack(L *lua.LState) (*HttpRequest, error) {
	var httpReq *HttpRequest
	if L.Get(1).Type() == lua.LTUserData {
		var ok bool
		httpReq, ok = L.CheckUserData(1).Value.(*HttpRequest)
		if !ok {
			return nil, fmt.Errorf("expected HttpRequest userdata, got %s", L.Get(1).Type().String())
		}
		println("HTTP GET request:", httpReq.URL)
	} else {
		httpReq = &HttpRequest{
			URL: L.ToString(1),
		}
	}
	return httpReq, nil
}

func httpRequest(request *HttpRequest) (*HttpResponse, error) {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	var req *http.Request
	var err error

	if request.Method == "GET" {
		req, err = http.NewRequest("GET", request.URL, http.NoBody)
	} else if request.Method == "POST" {
		req, err = http.NewRequest("POST", request.URL, strings.NewReader(request.Body))
		if err == nil {
			for key, value := range request.Headers {
				req.Header.Set(key, value)
			}
		}
	} else if request.Method == "PUT" {
		req, err = http.NewRequest("PUT", request.URL, strings.NewReader(request.Body))
		if err == nil {
			for key, value := range request.Headers {
				req.Header.Set(key, value)
			}
		}
	} else if request.Method == "DELETE" {
		req, err = http.NewRequest("DELETE", request.URL, http.NoBody)
		if err == nil {
			for key, value := range request.Headers {
				req.Header.Set(key, value)
			}
		}
	} else if request.Method == "PATCH" {
		req, err = http.NewRequest("PATCH", request.URL, strings.NewReader(request.Body))
		if err == nil {
			for key, value := range request.Headers {
				req.Header.Set(key, value)
			}
		}
	} else if request.Method == "HEAD" {
		req, err = http.NewRequest("HEAD", request.URL, http.NoBody)
		if err == nil {
			for key, value := range request.Headers {
				req.Header.Set(key, value)
			}
		}
	} else if request.Method == "OPTIONS" {
		req, err = http.NewRequest("OPTIONS", request.URL, http.NoBody)
		if err == nil {
			for key, value := range request.Headers {
				req.Header.Set(key, value)
			}
		}
	} else {
		return nil, fmt.Errorf("unsupported HTTP method: %s", request.Method)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Errorf("Failed to close response body: %v", err)
		}
	}(resp.Body)

	bodyStr, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	response := &HttpResponse{
		StatusCode: resp.StatusCode,
		Headers:    resp.Header,
		Body:       string(bodyStr),
		JsonBody:   nil,
	}

	if resp.Header.Get("Content-Type") == "application/json" {
		var jsonData map[string]interface{}
		if err := json.Unmarshal(bodyStr, &jsonData); err != nil {
			return nil, fmt.Errorf("failed to parse JSON response: %v", err)
		}
		response.JsonBody = jsonData
	}

	return response, nil
}

func luaHttpGet(L *lua.LState) int {
	if L.GetTop() != 1 {
		L.Push(lua.LNil)
		L.Push(lua.LString("Expected 1 argument: HttpRequest userdata or URL string"))
		return 2
	}
	httpReq, err := getHttpRequestFromStack(L)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	httpReq.Method = "GET"

	resp, err := httpRequest(httpReq)
	if err != nil {
		log.Errorf("HTTP GET request failed: %v", err)
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	ud := newHttpResponseUserData(L, resp)
	L.Push(ud)
	L.Push(lua.LNil) // No error

	return 2
}

func luaHttpPost(L *lua.LState) int {
	if L.GetTop() != 1 {
		L.Push(lua.LNil)
		L.Push(lua.LString("Expected 1 argument: HttpRequest userdata"))
		return 2
	}
	httpReq, err := getHttpRequestFromStack(L)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	httpReq.Method = "POST"

	resp, err := httpRequest(httpReq)
	if err != nil {
		log.Errorf("HTTP POST request failed: %v", err)
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	ud := newHttpResponseUserData(L, resp)
	L.Push(ud)
	L.Push(lua.LNil) // No error
	return 2
}

func luaHttpDownloadFile(L *lua.LState) int {
	if L.GetTop() != 2 {
		L.Push(lua.LBool(false))
		L.Push(lua.LString("Expected 2 arguments: url and filePath"))
		return 2
	}
	url := L.ToString(1)
	filePath := L.ToString(2)

	fileDir := filepath.Dir(filePath)
	if _, err := os.Stat(fileDir); os.IsNotExist(err) {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(fmt.Sprintf("Directory does not exist: %s", fileDir)))
		return 2
	}

	resp, err := http.Get(url)
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(err.Error()))
		return 2
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Errorf("Failed to close response body: %v", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		L.Push(lua.LBool(false))
		L.Push(lua.LString("Failed to download file: " + resp.Status))
		return 2
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(err.Error()))
		return 2
	}

	err = os.WriteFile(filePath, body, 0644)
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(lua.LBool(true))
	L.Push(lua.LNil)
	return 2
}
