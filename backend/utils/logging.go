package utils

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"runtime/debug"
	"strings"
)

type CustomFormatter struct {
	log.Formatter

	showCaller bool
}

func FormatPrefixString(prefix string, max int) string {
	if len(prefix) > max {
		prefix = prefix[:max]
	}

	var format = "[%-" + fmt.Sprintf("%d", max) + "s]"
	return fmt.Sprintf(format, strings.ToUpper(prefix))
}

func LoggingLevelString(level log.Level) string {
	switch level {
	case log.TraceLevel:
		return "TRACE"
	case log.DebugLevel:
		return "DEBUG"
	case log.InfoLevel:
		return "INFO"
	case log.WarnLevel:
		return "WARN"
	case log.ErrorLevel:
		return "ERROR"
	case log.FatalLevel:
		return "FATAL"
	case log.PanicLevel:
		return "PANIC"
	default:
		panic(fmt.Sprintf("Unknown log level: %v", level))
	}
}

func (f *CustomFormatter) FormatCaller(entry *log.Entry) string {
	if !entry.HasCaller() || !f.showCaller {
		return ""
	}
	file := entry.Caller.File
	line := entry.Caller.Line
	funcName := entry.Caller.Function
	if file != "" {
		return fmt.Sprintf("%s:%d %s", file, line, funcName)
	}
	return fmt.Sprintf("%s:%d", funcName, line)
}

func (f *CustomFormatter) FormatLua(entry *log.Entry) string {
	if entry.Data["lua"] == nil {
		return ""
	}

	source, hasSource := entry.Data["source"].(string)
	funcName, hasFunc := entry.Data["function"].(string)
	line, hasLine := entry.Data["line"].(int)

	formatted := ""
	switch {
	case hasSource && hasFunc && hasLine:
		formatted = fmt.Sprintf("[%s:%d %s]", source, line, funcName)
	case hasSource && hasFunc:
		formatted = fmt.Sprintf("[%s %s]", source, funcName)
	case hasSource && hasLine:
		formatted = fmt.Sprintf("[%s:%d]", source, line)
	case hasFunc && hasLine:
		formatted = fmt.Sprintf("[%s:%d]", funcName, line)
	case hasSource:
		formatted = fmt.Sprintf("[%s]", source)
	case hasFunc:
		formatted = fmt.Sprintf("[%s]", funcName)
	case hasLine:
		formatted = fmt.Sprintf("[%d]", line)
	}

	if formatted == "" {
		formatted = FormatPrefixString("LUA", 5)
		if plugin, ok := entry.Data["plugin"].(string); ok {
			formatted = formatted + FormatPrefixString(plugin, 36)
		}
	}

	return formatted
}

func (f *CustomFormatter) FormatPrefix(entry *log.Entry) string {
	if entry.Data["prefix"] == nil {
		return ""
	}

	prefix, ok := entry.Data["prefix"].(string)
	if !ok {
		log.Warn("Log entry prefix is not a string, skipping")
		return ""
	}

	return FormatPrefixString(prefix, 5)
}

func (f *CustomFormatter) Format(entry *log.Entry) ([]byte, error) {
	prefixes := []string{
		fmt.Sprintf("[%s]", entry.Time.Format("15:04:05")),
		FormatPrefixString(LoggingLevelString(entry.Level), 5),
	}

	caller := f.FormatCaller(entry)
	formattedLua := f.FormatLua(entry)
	if formattedLua != "" {
		prefixes = append(prefixes, formattedLua)
	}
	formattedPrefix := f.FormatPrefix(entry)
	if formattedPrefix != "" {
		if plugin, ok := entry.Data["plugin"].(string); ok {
			prefixes = append(prefixes, formattedPrefix+FormatPrefixString(plugin, 36))
		} else {
			prefixes = append(prefixes, formattedPrefix)
		}
	}

	formattedText := ""
	if caller != "" {
		formattedText = strings.Join(prefixes, " ") + " " + entry.Message + "\n" + fmt.Sprintf("  Caller: %s\n", caller)
	} else {
		formattedText = strings.Join(prefixes, " ") + " " + entry.Message + "\n"
	}

	// if error always append stacktrace
	if entry.Level == log.ErrorLevel || entry.Level == log.FatalLevel || entry.Level == log.PanicLevel {
		stackTrace := debug.Stack()
		formattedText += fmt.Sprintf("  Stack trace:\n%s\n", stackTrace)
	}

	return []byte(
		formattedText,
	), nil
}
