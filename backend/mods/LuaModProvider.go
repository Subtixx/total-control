package mods

import (
	"fmt"
	"github.com/yuin/gopher-lua"
	"os"
)

type LuaModProvider struct {
	L      *lua.LState
	plugin lua.LValue
}

func NewLuaModProvider(luaFile string) (*LuaModProvider, error) {
	L := lua.NewState()
	// Expose I/O helpers
	L.SetGlobal("read_file", L.NewFunction(func(L *lua.LState) int {
		path := L.ToString(1)
		data, err := os.ReadFile(path)
		if err != nil {
			L.Push(lua.LNil)
			L.Push(lua.LString(err.Error()))
			return 2
		}
		L.Push(lua.LString(data))
		return 1
	}))
	L.SetGlobal("write_file", L.NewFunction(func(L *lua.LState) int {
		path := L.ToString(1)
		content := L.ToString(2)
		err := os.WriteFile(path, []byte(content), 0644)
		if err != nil {
			L.Push(lua.LString(err.Error()))
			return 1
		}
		return 0
	}))
	L.SetGlobal("log_error", L.NewFunction(func(L *lua.LState) int {
		errMsg := L.ToString(1)
		fmt.Printf("Lua Error: %s\n", errMsg)
		return 0
	}))
	if err := L.DoFile(luaFile); err != nil {
		return nil, err
	}
	plugin := L.GetGlobal("plugin")
	return &LuaModProvider{L: L, plugin: plugin}, nil
}

func (p *LuaModProvider) call(method string, args ...lua.LValue) (lua.LValue, error) {
	fn := p.L.GetField(p.plugin, method)
	if fn.Type() != lua.LTFunction {
		return lua.LNil, fmt.Errorf("method '%s' not found in Lua plugin", method)
	}

	if err := p.L.CallByParam(lua.P{
		Fn:      fn,
		NRet:    1,
		Protect: true,
	}, args...); err != nil {
		return lua.LNil, fmt.Errorf("error calling Lua method '%s': %v", method, err)
	}

	ret := p.L.Get(-1)
	p.L.Pop(1)
	return ret, nil
}

// GetMods Implement ModProvider methods by calling Lua
func (p *LuaModProvider) GetMods() ([]Mod, error) {
	val, err := p.call("GetMods")
	if err != nil {
		return nil, fmt.Errorf("failed to get mods: %v", err)
	}

	if val.Type() != lua.LTTable {
		return nil, fmt.Errorf("expected Lua table, got %s", val.Type().String())
	}

	var mods []Mod
	luaTable := val.(*lua.LTable)
	luaTable.ForEach(func(_, value lua.LValue) {
		if value.Type() == lua.LTTable {
			modTable := value.(*lua.LTable)
			mod := Mod{
				ID:      modTable.RawGetString("id").String(),
				Name:    modTable.RawGetString("name").String(),
				Enabled: modTable.RawGetString("enabled").String() == "true",
				GameID:  modTable.RawGetString("game_id").String(),
			}
			mods = append(mods, mod)
		}
	})

	return mods, nil
}

func (p *LuaModProvider) GetModByID(id string) (*Mod, error) {
	val, err := p.call("GetModByID", lua.LString(id))
	if err != nil {
		return nil, err
	}

	if val.Type() != lua.LTTable {
		return nil, nil
	}

	modTable := val.(*lua.LTable)
	mod := &Mod{
		ID:      modTable.RawGetString("id").String(),
		Name:    modTable.RawGetString("name").String(),
		Enabled: modTable.RawGetString("enabled").String() == "true",
		GameID:  modTable.RawGetString("game_id").String(),
	}
	return mod, nil
}

func (p *LuaModProvider) GetGameModDirectory() (string, error) {
	val, err := p.call("GetGameModDirectory")
	if err != nil {
		return "", err
	}

	if val.Type() != lua.LTString {
		return "", nil
	}

	return val.String(), nil
}

func (p *LuaModProvider) AddMod(mod Mod) error {
	modTable := p.L.NewTable()
	modTable.RawSetString("id", lua.LString(mod.ID))
	modTable.RawSetString("name", lua.LString(mod.Name))
	modTable.RawSetString("enabled", lua.LString(fmt.Sprintf("%v", mod.Enabled)))
	modTable.RawSetString("game_id", lua.LString(mod.GameID))

	_, err := p.call("AddMod", modTable)
	return err
}

func (p *LuaModProvider) RemoveMod(id string) error {
	_, err := p.call("RemoveMod", lua.LString(id))
	return err
}

func (p *LuaModProvider) UpdateMod(mod Mod) error {
	modTable := p.L.NewTable()
	modTable.RawSetString("id", lua.LString(mod.ID))
	modTable.RawSetString("name", lua.LString(mod.Name))
	modTable.RawSetString("enabled", lua.LString(fmt.Sprintf("%v", mod.Enabled)))
	modTable.RawSetString("game_id", lua.LString(mod.GameID))

	_, err := p.call("UpdateMod", modTable)
	return err
}

func (p *LuaModProvider) ListGameMods() ([]Mod, error) {
	val, err := p.call("ListGameMods")
	if err != nil {
		return nil, err
	}

	if val.Type() != lua.LTTable {
		return nil, nil
	}

	var mods []Mod
	luaTable := val.(*lua.LTable)
	luaTable.ForEach(func(_, value lua.LValue) {
		if value.Type() == lua.LTTable {
			modTable := value.(*lua.LTable)
			mod := Mod{
				ID:      modTable.RawGetString("id").String(),
				Name:    modTable.RawGetString("name").String(),
				Enabled: modTable.RawGetString("enabled").String() == "true",
				GameID:  modTable.RawGetString("game_id").String(),
			}
			mods = append(mods, mod)
		}
	})

	return mods, nil
}

func (p *LuaModProvider) GetGameID() (string, error) {
	val, err := p.call("GetGameID")
	if err != nil {
		return "", err
	}

	if val.Type() != lua.LTString {
		return "", nil
	}

	return val.String(), nil
}

// Close Remember to close the Lua state when done
func (p *LuaModProvider) Close() {
	p.L.Close()
}
