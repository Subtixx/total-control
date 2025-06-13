package scripting

import (
	_ "embed"
	lua "github.com/yuin/gopher-lua"
)

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
