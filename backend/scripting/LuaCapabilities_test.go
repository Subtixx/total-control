package scripting

import (
	lua "github.com/yuin/gopher-lua"
	"testing"
)

func SetupLuaCapabilitiesTest() *lua.LState {
	L := lua.NewState()
	LuaRegisterCapabilitiesObject(L)
	return L
}

func TestLuaCheckCan(t *testing.T) {
	L := SetupLuaCapabilitiesTest()
	defer L.Close()
	err := L.DoString(`
		local result = capabilities.can("test")
		return result
	`)
	if err != nil {
		t.Errorf("Lua code execution failed: %v", err)
	}
	result := L.Get(-1)
	if result.Type() != lua.LTBool {
		t.Errorf("Expected a boolean, got %s", result.Type())
	}
	if result.(lua.LBool) == false {
		t.Errorf("Expected result to be true")
	}
}

func TestLuaCheckCanAccessNetwork(t *testing.T) {
	L := SetupLuaCapabilitiesTest()
	defer L.Close()
	err := L.DoString(`
		local result = capabilities.canAccessNetwork()
		return result
	`)
	if err != nil {
		t.Errorf("Lua code execution failed: %v", err)
	}
	result := L.Get(-1)
	if result.Type() != lua.LTBool {
		t.Errorf("Expected a boolean, got %s", result.Type())
	}
	if result.(lua.LBool) == false {
		t.Errorf("Expected result to be true")
	}
}

func TestLuaCheckCanAccessFileSystem(t *testing.T) {
	L := SetupLuaCapabilitiesTest()
	defer L.Close()
	err := L.DoString(`
		local result = capabilities.canAccessFileSystem()
		return result
	`)
	if err != nil {
		t.Errorf("Lua code execution failed: %v", err)
	}
	result := L.Get(-1)
	if result.Type() != lua.LTBool {
		t.Errorf("Expected a boolean, got %s", result.Type())
	}
	if result.(lua.LBool) == false {
		t.Errorf("Expected result to be true")
	}
}
