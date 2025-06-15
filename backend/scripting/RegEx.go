package scripting

import (
	"regexp"

	lua "github.com/yuin/gopher-lua"
)

func LuaRegisterRegExObject(L *lua.LState) {
	regexTable := L.NewTable()
	L.SetGlobal("regexp", regexTable)
	regexTable.RawSetString("findAll", L.NewFunction(luaRegexFindAll))
	regexTable.RawSetString("match", L.NewFunction(luaRegexMatch))
	regexTable.RawSetString("replace", L.NewFunction(luaRegexReplace))
}

// FindAll returns all matches of the pattern in the input string.
func luaRegexFindAll(L *lua.LState) int {
	pattern := L.CheckString(1)
	input := L.CheckString(2)
	re, err := regexp.Compile(pattern)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	matches := re.FindAllString(input, -1)
	lTable := L.NewTable()
	for _, match := range matches {
		lTable.Append(lua.LString(match))
	}
	L.Push(lTable)
	return 1
}

// Match returns true if the pattern matches the input string.
func luaRegexMatch(L *lua.LState) int {
	pattern := L.CheckString(1)
	input := L.CheckString(2)
	re, err := regexp.Compile(pattern)
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(lua.LBool(re.MatchString(input)))
	return 1
}

func luaRegexReplace(L *lua.LState) int {
	pattern := L.CheckString(1)
	replacement := L.CheckString(2)
	input := L.CheckString(3)
	re, err := regexp.Compile(pattern)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	result := re.ReplaceAllString(input, replacement)
	L.Push(lua.LString(result))
	return 1
}
