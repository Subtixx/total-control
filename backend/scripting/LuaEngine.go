package scripting

// TODO: measure time and kill the lua file if it takes too long to execute
// TODO: Create threads for each lua script execution

import (
	"TotalControl/backend/plugins"
	"TotalControl/backend/utils"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	lua "github.com/yuin/gopher-lua"
	"path"
	"strings"
	"time"
)

type LuaEngine struct {
	uuid  uuid.UUID
	L     *lua.LState `json:"-"`
	Cache *utils.Cache
}

type LuaEngineSetupOptions struct {
	Context       context.Context
	ContextValues map[string]interface{}
}

func NewLuaEngine(luaEngineId uuid.UUID) (*LuaEngine, error) {
	L := lua.NewState()
	engine := &LuaEngine{
		L:    L,
		uuid: luaEngineId,
	}
	err := engine.Setup(LuaEngineSetupOptions{
		Context:       context.Background(),
		ContextValues: nil,
	})
	if err != nil {
		return nil, err
	}
	return engine, nil
}

func GetLuaEngine(L *lua.LState) *LuaEngine {
	ctx := L.Context()
	if ctx == nil {
		log.Error("LuaEngine context is nil")
		return nil
	}
	val := ctx.Value("luaengine")
	if val == nil {
		log.Error("LuaEngine not found in context")
		return nil
	}

	engine, ok := val.(*LuaEngine)
	if !ok {
		log.Error("Value in context is not of type *LuaEngine")
		return nil
	}

	return engine
}

func GetLogger(L *lua.LState) *log.Entry {
	luaEngine := GetLuaEngine(L)
	if luaEngine != nil {
		return luaEngine.Logger()
	}
	return log.WithFields(log.Fields{
		"lua": true,
	})
}

func (l *LuaEngine) Logger() *log.Entry {
	return log.WithFields(log.Fields{
		"lua":    true,
		"plugin": l.uuid.String(),
	})
}

func (l *LuaEngine) Setup(options LuaEngineSetupOptions) error {
	if err := l.setupContext(options); err != nil {
		return err
	}
	if err := l.setupCache(); err != nil {
		return err
	}
	//l.L.OpenLibs()
	libsToLoad := []string{ // This could be really bad if we allow all libraries!
		lua.LoadLibName,
		lua.BaseLibName,
		lua.TabLibName,
		lua.IoLibName,
		lua.OsLibName,
		lua.StringLibName,
		lua.MathLibName,
		lua.DebugLibName,
		lua.ChannelLibName,
		lua.CoroutineLibName,
	}
	err := LoadBuiltinLibs(l.L, libsToLoad)
	if err != nil {
		return err
	}

	l.registerGlobals()
	if err := LoadLibs(l.L); err != nil {
		l.Logger().Errorf("Failed to load Lua libraries: %v", err)
		l.L.Close()
		return err
	}
	l.registerObjects()
	return nil
}

// setupContext sets the LuaEngine context
func (l *LuaEngine) setupContext(options LuaEngineSetupOptions) error {
	if l.uuid == uuid.Nil {
		return errors.New("LuaEngine UUID cannot be nil")
	}
	ctx := context.WithValue(options.Context, "luaengine", l)
	if options.ContextValues != nil {
		for key, value := range options.ContextValues {
			ctx = context.WithValue(ctx, key, value)
		}
	}
	l.L.SetContext(ctx)
	return nil
}

// setupCache initializes the cache for the LuaEngine
func (l *LuaEngine) setupCache() error {
	err := utils.CreateDirectoryIfNotExists("plugins/.cache")
	if err != nil {
		log.Errorf("Failed to create cache directory: %v", err)
	}
	pluginAppDataPath := plugins.GetPluginAppDataPath(l.uuid)
	l.Cache = utils.NewCache(path.Join(pluginAppDataPath, "cache.json"), l.uuid)
	return nil
}

// registerGlobals registers global Lua functions
func (l *LuaEngine) registerGlobals() {
	l.L.SetGlobal("print", l.L.NewFunction(luaPrint))
	l.L.SetGlobal("error_handler", l.L.NewFunction(luaErrorHandler))
	l.L.SetGlobal("table_size", l.L.NewFunction(luaTableSize))
}

// registerObjects registers custom objects and extensions
func (l *LuaEngine) registerObjects() {
	LuaRegisterCapabilitiesObject(l.L)
	LuaRegisterLogObject(l.L)
	LuaRegisterJsonObject(l.L)
	LuaRegisterRegExObject(l.L)
	LuaRegisterHttpObject(l.L)
	LuaRegisterCacheObject(l.L)
	LuaRegisterZipObject(l.L)

	LuaExtendTableObject(l.L)
	LuaExtendIoTable(l.L)
	LuaExtendOsTable(l.L)
}

func (l *LuaEngine) Shutdown() error {
	var err error = nil
	if l.Cache != nil {
		err = l.Cache.Save()
		if err != nil {
			l.Logger().Errorf("Failed to save cache: %v", err)
		}
	}
	l.L.Close()
	l.Logger().Debugf("LuaEngine %s has been shut down", l.uuid.String())
	return err
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

func (l *LuaEngine) DoStringWithTimeout(script string, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	errCh := make(chan error, 1)
	go func() {
		errCh <- l.L.DoString(script)
	}()

	select {
	case <-ctx.Done():
		l.L.Close() // forcibly close the Lua state
		return errors.New("lua script execution timed out")
	case err := <-errCh:
		return err
	}
}

func (l *LuaEngine) CallFunc(obj *lua.LTable, method *lua.LFunction, args ...lua.LValue) (lua.LValue, error) {
	if obj == nil || method == nil {
		return lua.LNil, errors.New("object or method is nil")
	}
	if obj.Type() != lua.LTTable {
		return lua.LNil, fmt.Errorf("object is not a Lua table, got %s", obj.Type().String())
	}
	if method.Type() != lua.LTFunction {
		return lua.LNil, fmt.Errorf("method is not a Lua function, got %s", method.Type().String())
	}

	l.L.Push(method)
	l.L.Push(obj)
	for _, arg := range args {
		l.L.Push(arg)
	}

	// +1 for the object itself
	if err := l.L.PCall(len(args)+1, lua.MultRet, nil); err != nil {
		return lua.LNil, fmt.Errorf("error calling method '%s': %v", method.String(), err)
	}

	return l.L.Get(-1), nil
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

func (l *LuaEngine) String() string {
	return "LuaEngine{" +
		"\n\tUUID: " + l.uuid.String() +
		"\n\tCache: " + l.Cache.String() +
		"\n\tLuaState: " + fmt.Sprintf("%p", l.L) +
		"\n}"
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
