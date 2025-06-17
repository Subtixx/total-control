package scripting

import (
	_ "embed"
	lua "github.com/yuin/gopher-lua"
)

type luaLib struct {
	libName string
	libFunc lua.LGFunction
}

//go:embed lib/serpent.lua
var serpentLua string

func LoadSerpentLib(L *lua.LState) error {
	if err := L.DoString(serpentLua); err != nil {
		return err
	}
	L.SetGlobal("serpent", L.Get(-1))
	return nil
}

func LoadLibs(L *lua.LState) error {
	if err := LoadSerpentLib(L); err != nil {
		return err
	}
	return nil
}

func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func LoadBuiltinLibs(L *lua.LState, libNames []string) error {
	luaLibs := []luaLib{
		//{lua.LoadLibName, lua.OpenPackage},
		//{lua.BaseLibName, lua.OpenBase},
		{lua.TabLibName, lua.OpenTable},
		//{lua.IoLibName, lua.OpenIo},
		//{lua.OsLibName, lua.OpenOs},
		{lua.StringLibName, lua.OpenString},
		{lua.MathLibName, lua.OpenMath},
		//{lua.DebugLibName, lua.OpenDebug},
		//{lua.ChannelLibName, lua.OpenChannel},
		{lua.CoroutineLibName, lua.OpenCoroutine},
	}

	for _, lib := range luaLibs {
		if !contains(libNames, lib.libName) {
			continue
		}

		if err := loadBuiltinLib(L, lib.libName, lib.libFunc); err != nil {
			return err
		}
	}

	err := loadBuiltinLib(L, lua.OsLibName, func(L *lua.LState) int {
		return lua.OpenOsBlacklist(L, "setlocale", "setenv", "remove", "rename")
	})
	if err != nil {
		return err
	}

	return nil
}

func loadBuiltinLib(L *lua.LState, libName string, libFunc lua.LGFunction) error {
	L.Push(L.NewFunction(libFunc))
	L.Push(lua.LString(libName))
	if err := L.PCall(1, 0, nil); err != nil {
		return err
	}
	return nil
}
