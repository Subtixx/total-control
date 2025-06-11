package scripting

import (
	"TotalControl/backend/utils"
	"fmt"
	log "github.com/sirupsen/logrus"
	lua "github.com/yuin/gopher-lua"
	"os"
	"strings"
)

type LuaEngine struct {
	L *lua.LState
}

func luaCreateJsonTable(l *lua.LState) *lua.LTable {
	jsonTable := l.NewTable()
	jsonTable.RawSetString("encode", l.NewFunction(utils.LuaJsonEncode))
	jsonTable.RawSetString("decode", l.NewFunction(utils.LuaJsonDecode))
	return jsonTable
}

func luaCreateLogTable(L *lua.LState) *lua.LTable {
	// Log table with (log.info, log.warn, log.error, etc.) functions
	logTable := L.NewTable()
	logTable.RawSetString("debug", L.NewFunction(func(L *lua.LState) int {
		msg := L.ToString(1)
		log.WithField("lua", true).Debug(msg)
		return 0
	}))
	logTable.RawSetString("info", L.NewFunction(func(L *lua.LState) int {
		msg := L.ToString(1)
		log.WithField("lua", true).Info(msg)
		return 0
	}))
	logTable.RawSetString("warn", L.NewFunction(func(L *lua.LState) int {
		msg := L.ToString(1)
		log.WithField("lua", true).Warn(msg)
		return 0
	}))
	logTable.RawSetString("error", L.NewFunction(func(L *lua.LState) int {
		msg := L.ToString(1)
		log.WithField("lua", true).Error(msg)
		return 0
	}))
	return logTable
}

func luaCreateOperatingSystemTable(L *lua.LState) *lua.LTable {
	osTable := L.NewTable()
	L.SetField(osTable, "Windows", lua.LNumber(1))
	L.SetField(osTable, "Linux", lua.LNumber(2))
	L.SetField(osTable, "MacOS", lua.LNumber(3))
	L.SetField(osTable, "Unknown", lua.LNumber(0))

	L.SetField(osTable, "getOperatingSystem", L.NewFunction(luaGetOperatingSystem))
	//L.SetField(osTable, "isWindows", L.NewFunction(luaIsWindows))
	//L.SetField(osTable, "isLinux", L.NewFunction(luaIsLinux))
	//L.SetField(osTable, "isMacOS", L.NewFunction(luaIsMacOS))
	L.SetField(osTable, "is_windows", lua.LFalse)
	L.SetField(osTable, "is_linux", lua.LFalse)
	L.SetField(osTable, "is_macos", lua.LFalse)
	L.SetField(osTable, "is_unknown", lua.LFalse)
	switch utils.GetOperatingSystem() {
	case utils.WindowsOS:
		L.SetField(osTable, "is_windows", lua.LTrue)
	case utils.LinuxOS:
		L.SetField(osTable, "is_linux", lua.LTrue)
	case utils.MacOS:
		L.SetField(osTable, "is_macos", lua.LTrue)
	default:
		L.SetField(osTable, "is_unknown", lua.LTrue)
	}
	L.SetField(osTable, "getenv", L.NewFunction(luaOsGetenv))
	return osTable
}

func luaPrint(L *lua.LState) int {
	// Print function that captures Lua print calls
	for i := 1; i <= L.GetTop(); i++ {
		if str, ok := L.Get(i).(lua.LString); ok {
			// Log with context lua
			log.WithField("lua", true).Info(str.String())
		} else {
			log.WithField("lua", true).Info(L.Get(i).Type().String(), ": ", L.Get(i))
		}
	}
	return 0
}

func luaErrorHandler(L *lua.LState) int {
	// Custom error handler for Lua
	if err := L.ToString(1); err != "" {
		println("Lua Error:", L.ToString(1))
	}
	return 0
}

func luaOsGetenv(L *lua.LState) int {
	key := L.ToString(1)
	if key == "" {
		L.Push(lua.LNil)
		return 1
	}

	value, exists := os.LookupEnv(key)
	if !exists {
		L.Push(lua.LNil)
		return 1
	}

	L.Push(lua.LString(value))
	return 1
}

func luaGetOperatingSystem(L *lua.LState) int {
	// 1 - Windows, 2 - Linux, 3 - MacOS, 0 - Unknown
	switch utils.GetOperatingSystem() {
	case utils.WindowsOS:
		L.Push(lua.LNumber(1))
		return 1
	case utils.LinuxOS:
		L.Push(lua.LNumber(2))
		return 1
	case utils.MacOS:
		L.Push(lua.LNumber(3))
		return 1
	default:
		L.Push(lua.LNumber(0))
		return 1
	}
}

func luaIsWindows(L *lua.LState) int {
	// Check if the current operating system is Windows
	if utils.GetOperatingSystem() == utils.WindowsOS {
		L.Push(lua.LTrue)
	} else {
		L.Push(lua.LFalse)
	}
	return 1
}

func luaIsLinux(L *lua.LState) int {
	// Check if the current operating system is Linux
	if utils.GetOperatingSystem() == utils.LinuxOS {
		L.Push(lua.LTrue)
	} else {
		L.Push(lua.LFalse)
	}
	return 1
}

func luaIsMacOS(L *lua.LState) int {
	// Check if the current operating system is MacOS
	if utils.GetOperatingSystem() == utils.MacOS {
		L.Push(lua.LTrue)
	} else {
		L.Push(lua.LFalse)
	}
	return 1
}

func NewLuaEngine() *LuaEngine {
	L := lua.NewState()
	engine := &LuaEngine{
		L: L,
	}
	engine.Setup()
	return engine
}

func (l *LuaEngine) Setup() {
	l.L.OpenLibs() // This could be really bad if we allow all libraries!
	l.L.SetGlobal("print", l.L.NewFunction(luaPrint))
	l.L.SetGlobal("error_handler", l.L.NewFunction(luaErrorHandler))
	l.L.SetGlobal("os_getenv", l.L.NewFunction(luaOsGetenv))
	l.L.SetGlobal("operating_system", luaCreateOperatingSystemTable(l.L))
	l.L.SetGlobal("json", luaCreateJsonTable(l.L)) // Assuming you have a function to create a JSON table

	l.L.SetGlobal("log", luaCreateLogTable(l.L))
}

func (l *LuaEngine) CheckLuaError(err error) {
	if err == nil {
		return
	}
	// Log the error with stack trace
	log.WithFields(log.Fields{
		"lua":   true,
		"stack": err.Error(),
	}).Error("Lua error occurred")
	l.L.RaiseError("Lua error: %s", err.Error())
}

func (l *LuaEngine) LoadScript(script string) error {
	if err := l.L.DoString(script); err != nil {
		return err
	}
	return nil
}

func (l *LuaEngine) LoadFile(filename string) error {
	if err := l.L.DoFile(filename); err != nil {
		return err
	}
	return nil
}

func (l *LuaEngine) CallGlobal(method string, args ...lua.LValue) (lua.LValue, error) {
	fn := l.L.GetGlobal(method)
	if fn.Type() != lua.LTFunction {
		return lua.LNil, fmt.Errorf("method '%s' not found in Lua engine", method)
	}

	l.L.Push(fn)
	for _, arg := range args {
		l.L.Push(arg)
	}
	if err := l.L.PCall(len(args), lua.MultRet, nil); err != nil {
		return lua.LNil, fmt.Errorf("error calling method '%s': %v", method, err)
	}

	return l.L.Get(-1), nil
}

func (l *LuaEngine) Call(obj lua.LValue, method string, args ...lua.LValue) (lua.LValue, error) {
	if obj.Type() != lua.LTTable {
		return lua.LNil, fmt.Errorf("object is not a Lua table")
	}

	fn := l.L.GetField(obj, method)
	if fn.Type() != lua.LTFunction {
		return lua.LNil, fmt.Errorf("method '%s' not found in Lua object", method)
	}

	l.L.Push(fn)
	l.L.Push(obj)
	for _, arg := range args {
		l.L.Push(arg)
	}

	// +1 for the object itself
	if err := l.L.PCall(len(args)+1, lua.MultRet, nil); err != nil {
		return lua.LNil, fmt.Errorf("error calling method '%s': %v", method, err)
	}

	return l.L.Get(-1), nil
}

func (l *LuaEngine) HasFunction(method string) bool {
	fn := l.L.GetGlobal(method)
	if fn.Type() == lua.LTFunction {
		return true
	}
	return false
}

func (l *LuaEngine) HasMethod(obj lua.LValue, method string) bool {
	if obj.Type() != lua.LTTable {
		return false
	}

	fn := l.L.GetField(obj, method)
	if fn.Type() == lua.LTFunction {
		return true
	}
	return false
}

func (l *LuaEngine) Close() {
	l.L.Close()
}

func (l *LuaEngine) debugPrintLuaState() {
	log.Debug("Available Lua functions/methods/objects:")
	var funcs []string
	l.L.G.Global.ForEach(func(key, value lua.LValue) {
		if value.Type() == lua.LTFunction {
			funcs = append(funcs, fmt.Sprintf("[FUNC]\t%s", key.String()))
		} else if value.Type() == lua.LTTable {
			// Iterate through the table to find methods
			funcs = append(funcs, fmt.Sprintf("[TAB ]\t%s", key.String()))
			value.(*lua.LTable).ForEach(func(subKey, subValue lua.LValue) {
				if subValue.Type() == lua.LTFunction {
					funcs = append(funcs, fmt.Sprintf("\t[METH]\t%s.%s", key.String(), subKey.String()))
				} else {
					funcs = append(funcs, fmt.Sprintf("\t[OBJ ]\t%s.%s (%s)", key.String(), subKey.String(), subValue.Type().String()))
				}
			})
		} else {
			funcs = append(funcs, fmt.Sprintf("[UNK ]\t%s (%s)", key.String(), value.Type().String()))
		}
	})
	// Join by \n\t
	log.Debugf("Available Lua functions/methods/objects:\n\t%s", strings.Join(funcs, "\n\t"))
}
