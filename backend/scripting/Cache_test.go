package scripting

import (
	"TotalControl/backend/utils"
	"context"
	"github.com/google/uuid"
	lua "github.com/yuin/gopher-lua"
	"testing"
)

func SetupCacheTests(t *testing.T) *lua.LState {
	// Create tmp
	id, err := uuid.NewUUID()
	if err != nil {
		t.Fatalf("Failed to create UUID: %v", err)
	}

	cachePath := t.TempDir() + "/cache_test.json"

	l := lua.NewState()
	ctx := context.WithValue(context.Background(), "cache", utils.NewCache(cachePath, id))
	l.SetContext(ctx)
	LuaRegisterCacheObject(l)
	return l
}

func TestCacheHasKey(t *testing.T) {
	l := SetupCacheTests(t)
	defer l.Close()

	luaCode := `
local key = "test_key"
cache.set(key, "test_value")
assert(cache.has(key) == true, "Expected cache.has to return true for existing key")
`
	if err := l.DoString(luaCode); err != nil {
		t.Errorf("Failed to execute Lua code: %v", err)
	}
}

func TestCacheGet(t *testing.T) {
	l := SetupCacheTests(t)
	defer l.Close()

	luaCode := `
local key = "test_key"
cache.set(key, "test_value")
local value = cache.get(key)
assert(value == "test_value", "Expected cache.get to return 'test_value'")
`
	if err := l.DoString(luaCode); err != nil {
		t.Errorf("Failed to execute Lua code: %v", err)
	}
}

func TestCacheSet(t *testing.T) {
	l := SetupCacheTests(t)
	defer l.Close()

	luaCode := `
local key = "test_key"
cache.set(key, "test_value")
local value = cache.get(key)
assert(value == "test_value", "Expected cache.get to return 'test_value' after set")
`
	if err := l.DoString(luaCode); err != nil {
		t.Errorf("Failed to execute Lua code: %v", err)
	}
}

func TestCacheDelete(t *testing.T) {
	l := SetupCacheTests(t)
	defer l.Close()

	luaCode := `
local key = "test_key"
cache.set(key, "test_value")
cache.delete(key)
local exists = cache.has(key)
assert(exists == false, "Expected cache.has to return false after delete")
`
	if err := l.DoString(luaCode); err != nil {
		t.Errorf("Failed to execute Lua code: %v", err)
	}
}
