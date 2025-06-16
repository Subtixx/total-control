package scripting

import (
	lua "github.com/yuin/gopher-lua"
	"testing"
)

func SetupLuaLibTests() *lua.LState {
	l := lua.NewState()
	err := LoadLibs(l)
	if err != nil {
		return nil
	}
	err = LoadBuiltinLibs(l, []string{
		lua.MathLibName,
	})
	if err != nil {
		return nil
	}
	return l
}

func TestLuaMathLoaded(t *testing.T) {
	l := SetupLuaLibTests()
	defer l.Close()
	err := l.DoString(`
		local result = math.sin(math.pi / 2)
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
		t.Errorf("Expected result to be 1")
	}
}
