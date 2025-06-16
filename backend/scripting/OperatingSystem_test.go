package scripting

import (
	lua "github.com/yuin/gopher-lua"
	"testing"
)

func SetupOperatingSystemTests() *lua.LState {
	l := lua.NewState()
	LuaExtendOsTable(l)
	return l
}

func TestLuaGetOperatingSystem(t *testing.T) {
	l := SetupOperatingSystemTests()
	defer l.Close()
	err := l.DoString(`
		local result = os.getOperatingSystem()
		return result
	`)
	if err != nil {
		t.Errorf("Lua code execution failed: %v", err)
	}
	result := l.Get(-1)
	if result.Type() != lua.LTNumber {
		t.Errorf("Expected a number, got %s", result.Type())
	}
	if result.(lua.LNumber) == 0 {
		t.Errorf("Expected result to be non-zero")
	}
}

func TestLuaIsWindows(t *testing.T) {
	l := SetupOperatingSystemTests()
	defer l.Close()
	err := l.DoString(`
		local result = os.isWindows()
		return result
	`)
	if err != nil {
		t.Errorf("Lua code execution failed: %v", err)
	}
	result := l.Get(-1)
	if result.Type() != lua.LTBool {
		t.Errorf("Expected a boolean, got %s", result.Type())
	}
}

func TestLuaIsWindowsConst(t *testing.T) {
	l := SetupOperatingSystemTests()
	defer l.Close()
	err := l.DoString(`
		local result = os.is_windows
		return result
	`)
	if err != nil {
		t.Errorf("Lua code execution failed: %v", err)
	}
	result := l.Get(-1)
	if result.Type() != lua.LTBool {
		t.Errorf("Expected a boolean, got %s", result.Type())
	}
}

func TestLuaIsLinux(t *testing.T) {
	l := SetupOperatingSystemTests()
	defer l.Close()
	err := l.DoString(`
		local result = os.isLinux()
		return result
	`)
	if err != nil {
		t.Errorf("Lua code execution failed: %v", err)
	}
	result := l.Get(-1)
	if result.Type() != lua.LTBool {
		t.Errorf("Expected a boolean, got %s", result.Type())
	}
}

func TestLuaIsLinuxConst(t *testing.T) {
	l := SetupOperatingSystemTests()
	defer l.Close()
	err := l.DoString(`
		local result = os.is_linux
		return result
	`)
	if err != nil {
		t.Errorf("Lua code execution failed: %v", err)
	}
	result := l.Get(-1)
	if result.Type() != lua.LTBool {
		t.Errorf("Expected a boolean, got %s", result.Type())
	}
}

func TestLuaIsMac(t *testing.T) {
	l := SetupOperatingSystemTests()
	defer l.Close()
	err := l.DoString(`
		local result = os.isMac()
		return result
	`)
	if err != nil {
		t.Errorf("Lua code execution failed: %v", err)
	}
	result := l.Get(-1)
	if result.Type() != lua.LTBool {
		t.Errorf("Expected a boolean, got %s", result.Type())
	}
}

func TestLuaIsMacConst(t *testing.T) {
	l := SetupOperatingSystemTests()
	defer l.Close()
	err := l.DoString(`
		local result = os.is_mac
		return result
	`)
	if err != nil {
		t.Errorf("Lua code execution failed: %v", err)
	}
	result := l.Get(-1)
	if result.Type() != lua.LTBool {
		t.Errorf("Expected a boolean, got %s", result.Type())
	}
}
