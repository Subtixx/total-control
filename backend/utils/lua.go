package utils

import (
	lua "github.com/yuin/gopher-lua"
	"strings"
)

func LuaArgsToString(args []lua.LValue) string {
	var sb strings.Builder
	for i, arg := range args {
		if i > 0 {
			sb.WriteString(", ")
		}
		switch arg.Type() {
		case lua.LTNil:
			sb.WriteString("nil")
		case lua.LTBool:
			if arg.(lua.LBool) {
				sb.WriteString("true")
			} else {
				sb.WriteString("false")
			}
		case lua.LTNumber:
			sb.WriteString(arg.String())
		case lua.LTString:
			sb.WriteString(`"` + arg.String() + `"`)
		default:
			sb.WriteString(arg.String())
		}
	}
	return sb.String()
}
