package scripting

import (
	"TotalControl/backend/utils"
	log "github.com/sirupsen/logrus"
	lua "github.com/yuin/gopher-lua"
)

func luaCacheHasKey(L *lua.LState) int {
	key := L.ToString(1)
	if key == "" {
		L.RaiseError("cache.has_key: key cannot be empty")
		return 0
	}

	// Get LuaEngine from context
	engine := GetLuaEngine(L)
	if engine == nil {
		L.RaiseError("cache.has_key: LuaEngine not found in context")
		return 0
	}

	exists, err := engine.cache.HasKey(key)
	if err != nil {
		L.RaiseError("cache.has_key: %v", err)
		return 0
	}

	L.Push(lua.LBool(exists)) // Push true or false based on existence
	return 1
}

// luaCacheGet retrieves a value from the cache by key.
func luaCacheGet(L *lua.LState) int {
	key := L.ToString(1)
	if key == "" {
		L.RaiseError("cache.get: key cannot be empty")
		return 0
	}

	// Get LuaEngine from context
	engine := GetLuaEngine(L)
	if engine == nil {
		L.RaiseError("cache.get: LuaEngine not found in context")
		return 0
	}
	value, err := engine.cache.Get(key)
	if err != nil {
		L.RaiseError("cache.get: %v", err)
		return 0
	}

	if value == "" {
		L.Push(lua.LNil)
		return 1
	}

	L.Push(utils.ToLuaValue(L, value))
	return 1
}

// luaCacheSet sets a value in the cache with a key and optional expiration time.
func luaCacheSet(L *lua.LState) int {
	if L.GetTop() < 2 {
		L.RaiseError("cache.set: at least key and value are required")
		return 0
	}

	key := L.ToString(1)
	value := L.Get(2)
	var goValue interface{}

	switch v := value.(type) {
	case lua.LString:
		goValue = string(v)
	case lua.LNumber:
		goValue = int(v)
	case lua.LBool:
		goValue = bool(v)
	case *lua.LTable:
		goValue = utils.LuaTableToMap(L, v)
	default:
		L.RaiseError("cache.set: unsupported value type")
		return 0
	}
	log.Debugf("Setting cache key: %s, value: %v", key, goValue)

	engine := GetLuaEngine(L)
	if engine == nil {
		L.RaiseError("cache.set: LuaEngine not found in context")
		return 0
	}

	expiration := -1 // Default to no expiration
	if L.GetTop() > 2 {
		exp := L.ToInt(3)
		if exp < 0 {
			L.RaiseError("cache.set: expiration must be non-negative")
			return 0
		}
		expiration = exp
	}
	err := engine.cache.Set(key, goValue, expiration)
	if err != nil {
		L.RaiseError("cache.set: %v", err)
		return 0
	}
	L.Push(lua.LTrue) // Return true on success
	return 1
}

// luaCacheDelete removes a value from the cache by key.
func luaCacheDelete(L *lua.LState) int {
	key := L.ToString(1)
	if key == "" {
		L.RaiseError("cache.delete: key cannot be empty")
		return 0
	}

	engine := GetLuaEngine(L)
	if engine == nil {
		L.RaiseError("cache.delete: LuaEngine not found in context")
		return 0
	}

	err := engine.cache.Delete(key)
	if err != nil {
		L.RaiseError("cache.delete: %v", err)
		return 0
	}

	L.Push(lua.LTrue)
	return 1
}

// luaCacheClear clears the entire cache.
func luaCacheClear(L *lua.LState) int {
	engine := GetLuaEngine(L)
	if engine == nil {
		L.RaiseError("cache.clear: LuaEngine not found in context")
		return 0
	}

	err := engine.cache.Clear()
	if err != nil {
		L.RaiseError("cache.clear: %v", err)
		return 0
	}

	L.Push(lua.LTrue)
	return 1
}

func luaRegisterCacheObject(L *lua.LState) {
	cache := L.NewTable()
	L.SetGlobal("cache", cache)

	L.SetField(cache, "has", L.NewFunction(luaCacheHasKey))
	L.SetField(cache, "get", L.NewFunction(luaCacheGet))
	L.SetField(cache, "set", L.NewFunction(luaCacheSet))
	L.SetField(cache, "delete", L.NewFunction(luaCacheDelete))
	L.SetField(cache, "clear", L.NewFunction(luaCacheClear))
}
