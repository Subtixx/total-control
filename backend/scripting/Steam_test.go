package scripting

import (
	lua "github.com/yuin/gopher-lua"
	"testing"
)

func TestLuaGetAppDetails(t *testing.T) {
	L := lua.NewState()
	LuaRegisterSteamObject(L)
	err := LoadLibs(L)
	if err != nil {
		t.Fatalf("Error loading libraries: %v", err)
	}
	err = L.DoString(`return Steam.GetAppDetails("440")`)
	if err != nil {
		t.Errorf("Error executing Lua script: %v", err)
	}
	result := L.Get(-1)
	if result.Type() != lua.LTTable {
		t.Errorf("Expected a table, got %s, %v", result.Type(), result)
	} else {
		t.Logf("Lua script executed successfully, result: %s", result.String())
	}
}
