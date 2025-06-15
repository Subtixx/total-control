package scripting

import (
	"github.com/yuin/gopher-lua"
	"os"
	"testing"
)

func TestLuaHttpGet(t *testing.T) {
	L := lua.NewState()
	LuaRegisterHttpObject(L)
	defer L.Close()
	// Register your http.get function here if needed

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
	L := lua.NewState()
	LuaRegisterHttpObject(L)
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
	L := lua.NewState()
	LuaRegisterHttpObject(L)
	defer L.Close()

	script := `
		local ok, err = http.download("https://httpbin.org/image/png", "test_image.png")
		if not ok then error("http.download failed: " .. tostring(err)) end
	`
	if err := L.DoString(script); err != nil {
		t.Fatalf("Lua http.download failed: %v", err)
	}
	// Check if the file was created
	if _, err := os.Stat("test_image.png"); os.IsNotExist(err) {
		t.Fatalf("Downloaded file does not exist: %v", err)
	}
	err := os.Remove("test_image.png")
	if err != nil {
		t.Fatalf("Failed to remove test_image.png: %v", err)
	}
}
