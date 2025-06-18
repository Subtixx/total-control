package scripting

import (
	"TotalControl/backend/utils"
	log "github.com/sirupsen/logrus"
	lua "github.com/yuin/gopher-lua"
)

func LuaRegisterCacheObject(L *lua.LState) {
	cache := L.NewTable()
	L.SetGlobal("cache", cache)

	L.SetField(cache, "has", L.NewFunction(luaCacheHasKey))
	L.SetField(cache, "get", L.NewFunction(luaCacheGet))
	L.SetField(cache, "set", L.NewFunction(luaCacheSet))
	L.SetField(cache, "delete", L.NewFunction(luaCacheDelete))
	L.SetField(cache, "clear", L.NewFunction(luaCacheClear))
}

func GetCache(L *lua.LState) *utils.Cache {
	ctx := L.Context()
	if ctx == nil {
		log.Error("LuaEngine context is nil")
		return nil
	}

	cache := ctx.Value("cache")
	if cache != nil {
		if c, ok := cache.(*utils.Cache); ok {
			return c
		}
		log.Error("Value in context is not of type *utils.Cache")
		return nil
	}

	engine := GetLuaEngine(L)
	if engine != nil {
		return engine.Cache
	}

	L.RaiseError("GetCache: LuaEngine not found in context")
	return nil
}

func luaCacheHasKey(L *lua.LState) int {
	if L.GetTop() < 1 {
		L.RaiseError("cache.has_key: key is required")
		return 0
	}

	cache := GetCache(L)
	if cache == nil {
		L.RaiseError("cache.has_key: Cache not found in context")
		return 0
	}

	key := L.ToString(1)
	if key == "" {
		L.RaiseError("cache.has_key: key cannot be empty")
		return 0
	}

	exists, err := cache.HasKey(key)
	if err != nil {
		L.RaiseError("cache.has_key: %v", err)
		return 0
	}

	L.Push(lua.LBool(exists))
	return 1
}

// luaCacheGet retrieves a value from the cache by key.
func luaCacheGet(L *lua.LState) int {
	if L.GetTop() < 1 {
		L.RaiseError("cache.get: key is required")
		return 0
	}

	cache := GetCache(L)
	if cache == nil {
		L.RaiseError("cache.get: Cache not found in context")
		return 0
	}

	key := L.ToString(1)
	if key == "" {
		L.RaiseError("cache.get: key cannot be empty")
		return 0
	}

	value, err := cache.Get(key)
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

	cache := GetCache(L)
	if cache == nil {
		L.RaiseError("cache.set: Cache not found in context")
		return 0
	}

	key := L.ToString(1)
	value := L.Get(2)

	goValue := utils.FromLuaValue(L, value)

	expiration := -1 // Default to no expiration
	if L.GetTop() > 2 {
		exp := L.ToInt(3)
		if exp < 0 {
			L.RaiseError("cache.set: expiration must be non-negative")
			return 0
		}
		expiration = exp
	}
	err := cache.Set(key, goValue, expiration)
	if err != nil {
		L.RaiseError("cache.set: %v", err)
		return 0
	}
	L.Push(lua.LTrue) // Return true on success
	return 1
}

// luaCacheDelete removes a value from the cache by key.
func luaCacheDelete(L *lua.LState) int {
	if L.GetTop() < 1 {
		L.RaiseError("cache.delete: key is required")
		return 0
	}

	cache := GetCache(L)
	if cache == nil {
		L.RaiseError("cache.delete: Cache not found in context")
		return 0
	}

	key := L.ToString(1)
	if key == "" {
		L.RaiseError("cache.delete: key cannot be empty")
		return 0
	}

	err := cache.Delete(key)
	if err != nil {
		L.RaiseError("cache.delete: %v", err)
		return 0
	}

	L.Push(lua.LTrue)
	return 1
}

// luaCacheClear clears the entire cache.
func luaCacheClear(L *lua.LState) int {
	cache := GetCache(L)
	if cache == nil {
		L.RaiseError("cache.clear: Cache not found in context")
		return 0
	}

	err := cache.Clear()
	if err != nil {
		L.RaiseError("cache.clear: %v", err)
		return 0
	}

	L.Push(lua.LTrue)
	return 1
}
