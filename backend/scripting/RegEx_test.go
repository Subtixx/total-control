package scripting

import (
	lua "github.com/yuin/gopher-lua"
	"testing"
)

func SetupRegExTests() *lua.LState {
	L := lua.NewState()
	LuaRegisterRegExObject(L)
	return L
}

func TestLuaRegexFindAll(t *testing.T) {
	L := SetupRegExTests()
	defer L.Close()
	err := L.DoString(`
		local result = regexp.findAll("a", "abcabc")
		return result
	`)
	if err != nil {
		t.Fatalf("Error executing Lua script: %v", err)
	}

	result := L.Get(-1)
	if result.Type() != lua.LTTable {
		t.Fatalf("Expected a table, got %s", result.Type())
	}

	expected := []string{"a", "a"}
	table := result.(*lua.LTable)
	for i, v := range expected {
		if table.RawGetInt(i+1).String() != v {
			t.Errorf("Expected %s at index %d, got %s", v, i+1, table.RawGetInt(i+1).String())
		}
	}
}

func TestLuaRegexMatch(t *testing.T) {
	L := SetupRegExTests()
	defer L.Close()
	err := L.DoString(`
		local result = regexp.match("abc", "abc")
		return result
	`)
	if err != nil {
		t.Fatalf("Error executing Lua script: %v", err)
	}

	result := L.Get(-1)
	if result.Type() != lua.LTBool {
		t.Fatalf("Expected a boolean, got %s", result.Type())
	}

	if !result.(lua.LBool) {
		t.Error("Expected match to be true, but got false")
	}
}

func TestLuaRegexReplace(t *testing.T) {
	L := SetupRegExTests()
	defer L.Close()
	err := L.DoString(`
		local result = regexp.replace("a", "x", "abcabc")
		return result
	`)
	if err != nil {
		t.Fatalf("Error executing Lua script: %v", err)
	}

	result := L.Get(-1)
	if result.Type() != lua.LTString {
		t.Fatalf("Expected a string, got %s", result.Type())
	}

	expected := "xbcxbc"
	if result.String() != expected {
		t.Errorf("Expected %s, got %s", expected, result.String())
	}
}
