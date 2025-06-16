package scripting

import (
	"github.com/yuin/gopher-lua"
	"os"
	"testing"
)

func SetupHttpTests() *lua.LState {
	L := lua.NewState()
	LuaRegisterHttpObject(L)
	return L
}

func TestLuaHttpGet(t *testing.T) {
	L := SetupHttpTests()
	defer L.Close()

	script := `
		local res, err = http.get("https://httpbin.org/get")
		if not res then error("http.get failed: " .. tostring(err)) end
		local body = res.body
		if not body or body == "" then
			error("http.get returned empty body")
		end
	`
	if err := L.DoString(script); err != nil {
		t.Fatalf("Lua http.get failed: %v", err)
	}
}

func TestLuaHttpPost(t *testing.T) {
	L := SetupHttpTests()
	defer L.Close()

	script := `
		local res, err = http.post("https://httpbin.org/post", { key = "value" })
		if not res then error("http.post failed: " .. tostring(err)) end
		local body = res.body
		if not body or body == "" then
			error("http.post returned empty body")
		end
	`
	if err := L.DoString(script); err != nil {
		t.Fatalf("Lua http.post failed: %v", err)
	}
}

func TestLuaHttpDownload(t *testing.T) {
	tmpDir := t.TempDir()
	L := SetupHttpTests()
	defer L.Close()

	script := `
		local ok, err = http.download("https://httpbin.org/image/png", "` + tmpDir + `/test_image.png")
		if not ok then error("http.download failed: " .. tostring(err)) end
	`
	if err := L.DoString(script); err != nil {
		t.Fatalf("Lua http.download failed: %v", err)
	}

	if _, err := os.Stat(tmpDir + "/test_image.png"); os.IsNotExist(err) {
		t.Fatalf("Downloaded file does not exist: %v", err)
	}
}
