package scripting

import (
	lua "github.com/yuin/gopher-lua"
	"testing"
)

func SetupJsonTests() *lua.LState {
	l := lua.NewState()
	LuaRegisterJsonObject(l)
	return l
}

func TestJsonEncode(t *testing.T) {
	l := SetupJsonTests()
	defer l.Close()

	// Test encoding a Lua table to JSON
	luaCode := `
		local t = {name = "John", age = 30, isEmployed = true}
		return json.encode(t)
	`
	if err := l.DoString(luaCode); err != nil {
		t.Fatalf("Failed to execute Lua code: %v", err)
	}

	result := l.Get(-1).String()
	expected := `{"age":30,"isEmployed":true,"name":"John"}`
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func TestJsonDecode(t *testing.T) {
	l := SetupJsonTests()
	defer l.Close()

	// Test decoding a JSON string to a Lua table
	luaCode := `
		local jsonStr = '{"name":"John","age":30,"isEmployed":true}'
		return json.decode(jsonStr)
	`
	if err := l.DoString(luaCode); err != nil {
		t.Fatalf("Failed to execute Lua code: %v", err)
	}

	result := l.Get(-1).(*lua.LTable)
	name := result.RawGetString("name").String()
	age := result.RawGetString("age").(lua.LNumber)
	isEmployed := result.RawGetString("isEmployed").(lua.LBool)

	if name != "John" || age != 30 || !isEmployed {
		t.Errorf("Decoded values do not match expected: %s, %d, %t", name, age, isEmployed)
	}
}
