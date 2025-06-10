package scripting

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	lua "github.com/yuin/gopher-lua"
)

func TestNewLuaEngine_Setup(t *testing.T) {
	engine := NewLuaEngine()
	defer engine.Close()

	assert.NotNil(t, engine.L)
	assert.Equal(t, lua.LTFunction, engine.L.GetGlobal("print").Type())
	assert.Equal(t, lua.LTTable, engine.L.GetGlobal("operating_system").Type())
	assert.Equal(t, lua.LTTable, engine.L.GetGlobal("log").Type())
}

func TestLuaEngine_LoadScript(t *testing.T) {
	engine := NewLuaEngine()
	defer engine.Close()

	err := engine.LoadScript(`a = 123`)
	assert.NoError(t, err)
	val := engine.L.GetGlobal("a")
	assert.Equal(t, lua.LTNumber, val.Type())
	assert.Equal(t, lua.LNumber(123), val)
}

func TestLuaEngine_LoadFile(t *testing.T) {
	engine := NewLuaEngine()
	defer engine.Close()

	f, err := os.CreateTemp("", "test.lua")
	assert.NoError(t, err)
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			t.Errorf("Failed to remove temp file %s: %v", name, err)
		}
	}(f.Name())
	_, err = f.WriteString("b = 456")
	assert.NoError(t, err)
	err = f.Close()
	assert.NoError(t, err)

	err = engine.LoadFile(f.Name())
	assert.NoError(t, err)
	val := engine.L.GetGlobal("b")
	assert.Equal(t, lua.LTNumber, val.Type())
	assert.Equal(t, lua.LNumber(456), val)
}

func TestLuaEngine_CallGlobal(t *testing.T) {
	engine := NewLuaEngine()
	defer engine.Close()

	err := engine.LoadScript(`function add(a, b) return a + b end`)
	assert.NoError(t, err)
	result, err := engine.CallGlobal("add", lua.LNumber(2), lua.LNumber(3))
	assert.NoError(t, err)
	assert.Equal(t, lua.LNumber(5), result)
}

func TestLuaEngine_Call(t *testing.T) {
	engine := NewLuaEngine()
	defer engine.Close()

	err := engine.LoadScript(`
		obj = {
			mul = function(self, a, b) print(tostring(a)) print(tostring(b)) return a * b end
		}
	`)
	assert.NoError(t, err)
	obj := engine.L.GetGlobal("obj")
	result, err := engine.Call(obj, "mul", lua.LNumber(2), lua.LNumber(4))
	assert.NoError(t, err)
	assert.Equal(t, lua.LNumber(8), result)
}

func TestLuaEngine_CheckLuaError(t *testing.T) {
	engine := NewLuaEngine()
	defer engine.Close()

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic from RaiseError")
		}
	}()
	engine.CheckLuaError(assert.AnError)
}

func TestLuaEngine_debugPrintLuaState(t *testing.T) {
	engine := NewLuaEngine()
	defer engine.Close()
	engine.debugPrintLuaState()
}

func TestLuaOsGetenv(t *testing.T) {
	engine := NewLuaEngine()
	defer engine.Close()

	err := os.Setenv("LUA_TEST_ENV", "hello")
	assert.NoError(t, err)
	defer func() {
		err := os.Unsetenv("LUA_TEST_ENV")
		if err != nil {
			t.Errorf("Failed to unset environment variable: %v", err)
		}
	}()

	engine.L.SetGlobal("os_getenv", engine.L.NewFunction(luaOsGetenv))
	err = engine.L.DoString(`val = os_getenv("LUA_TEST_ENV")`)
	assert.NoError(t, err)
	val := engine.L.GetGlobal("val")
	assert.Equal(t, lua.LString("hello"), val)
}
