package scripting

import (
	lua "github.com/yuin/gopher-lua"
	"testing"
)

func SetupTableTest() *lua.LState {
	L := lua.NewState()
	LuaExtendTableObject(L)
	return L
}

// Test keys
func TestTableKeys(t *testing.T) {
	L := SetupTableTest()
	defer L.Close()

	script := `
		local tbl = { a = 1, b = 2, c = 3 }
		local keys = table.keys(tbl)
		if #keys ~= 3 then
			error("Expected 3 keys, got " .. #keys)
		end
		local foundA, foundB, foundC = false, false, false
		for _, key in ipairs(keys) do
			if key == "a" then foundA = true end
			if key == "b" then foundB = true end
			if key == "c" then foundC = true end
		end
		if not (foundA and foundB and foundC) then
			error("Keys do not match expected keys")
		end
	`
	if err := L.DoString(script); err != nil {
		t.Fatalf("Lua table.keys failed: %v", err)
	}
}

// Test values
func TestTableValues(t *testing.T) {
	L := SetupTableTest()
	defer L.Close()

	script := `
		local tbl = { a = 1, b = 2, c = 3 }
		local values = table.values(tbl)
		if #values ~= 3 then
			error("Expected 3 values, got " .. #values)
		end
		local found1, found2, found3 = false, false, false
		for _, value in ipairs(values) do
			if value == 1 then found1 = true end
			if value == 2 then found2 = true end
			if value == 3 then found3 = true end
		end
		if not (found1 and found2 and found3) then
			error("Values do not match expected values")
		end
	`
	if err := L.DoString(script); err != nil {
		t.Fatalf("Lua table.values failed: %v", err)
	}
}

// Test length
func TestTableLength(t *testing.T) {
	L := SetupTableTest()
	defer L.Close()

	script := `
		local tbl = { a = 1, b = 2, c = 3 }
		local length = table.length(tbl)
		if length ~= 3 then
			error("Expected length 3, got " .. length)
		end
	`
	if err := L.DoString(script); err != nil {
		t.Fatalf("Lua table.length failed: %v", err)
	}
}
