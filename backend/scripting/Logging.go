package scripting

import (
	log "github.com/sirupsen/logrus"
	lua "github.com/yuin/gopher-lua"
)

func luaLog(L *lua.LState, stackTrace bool) (string, []interface{}, log.Fields) {
	if L.GetTop() < 1 {
		log.WithField("lua", true).Error("No message provided for debug log")
		return "", nil, nil
	}
	msg := L.ToString(1)

	var args []interface{}
	for i := 2; i <= L.GetTop(); i++ {
		args = append(args, L.Get(i))
	}

	fields := log.Fields{
		"lua": true,
	}

	if stackTrace {
		stack, ok := L.GetStack(1)
		if ok {
			_, err := L.GetInfo("nSl", stack, lua.LNil)
			if err == nil {
				if stack.Name != "" {
					fields["func"] = stack.Name
				}
				if stack.Source != "" {
					fields["source"] = stack.Source
				}
				if stack.CurrentLine > 0 {
					fields["line"] = stack.CurrentLine
				}
			}
		}
	}

	return msg, args, fields
}

func luaLogDebug(L *lua.LState) int {
	msg, args, fields := luaLog(L, true)
	if msg == "" {
		return 0
	}
	log.WithFields(fields).Debugf(msg, args...)
	return 0
}

func luaLogInfo(L *lua.LState) int {
	msg, args, fields := luaLog(L, false)
	if msg == "" {
		return 0
	}
	log.WithFields(fields).Infof(msg, args...)
	return 0
}

func luaLogWarn(L *lua.LState) int {
	msg, args, fields := luaLog(L, false)
	if msg == "" {
		return 0
	}
	log.WithFields(fields).Warnf(msg, args...)
	return 0
}

func luaLogError(L *lua.LState) int {
	msg, args, fields := luaLog(L, true)
	if msg == "" {
		return 0
	}
	log.WithFields(fields).Errorf(msg, args...)
	return 0
}

func luaLogFatal(L *lua.LState) int {
	msg, args, fields := luaLog(L, true)
	if msg == "" {
		return 0
	}
	log.WithFields(fields).Fatalf(msg, args...)
	return 0
}

func luaRegisterLogObject(L *lua.LState) {
	logTable := L.NewTable()
	logTable.RawSetString("debug", L.NewFunction(luaLogDebug))
	logTable.RawSetString("info", L.NewFunction(luaLogInfo))
	logTable.RawSetString("warn", L.NewFunction(luaLogWarn))
	logTable.RawSetString("error", L.NewFunction(luaLogError))
	logTable.RawSetString("fatal", L.NewFunction(luaLogFatal))
	L.SetGlobal("log", logTable)
}
