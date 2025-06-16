package scripting

import (
	lua "github.com/yuin/gopher-lua"
	"testing"
)

func SetupPluginManagerTests() *lua.LState {
	L := lua.NewState()
	RegisterPluginManagerObject(L)
	return L
}

func TestLuaIsPluginLoaded(t *testing.T) {
	L := SetupPluginManagerTests()
	defer L.Close()

	luaCode := `
loaded = plugin.isLoaded("testPlugin")
assert(loaded == false, "Expected plugin.isLoaded to return false")
`
	if err := L.DoString(luaCode); err != nil {
		t.Errorf("Lua code execution failed: %v", err)
	}
}

func TestLuaIsPluginEnabled(t *testing.T) {
	L := SetupPluginManagerTests()
	defer L.Close()

	luaCode := `
enabled = plugin.isEnabled("testPlugin")
assert(enabled == false, "Expected plugin.isEnabled to return false")
`
	if err := L.DoString(luaCode); err != nil {
		t.Errorf("Lua code execution failed: %v", err)
	}
}
