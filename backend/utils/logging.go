package utils

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"strings"
)

type CustomFormatter struct {
	log.Formatter
}

func (f *CustomFormatter) Format(entry *log.Entry) ([]byte, error) {
	prefixes := []string{
		fmt.Sprintf("[%s]", entry.Time.Format("15:04:05")),
		fmt.Sprintf("[%s]", strings.ToUpper(entry.Level.String())),
	}

	if entry.Data["lua"] != nil {
		source, hasSource := entry.Data["source"].(string)
		funcName, hasFunc := entry.Data["function"].(string)
		line, hasLine := entry.Data["line"].(int)
		if hasSource && hasFunc && hasLine {
			prefixes = append(prefixes, fmt.Sprintf("[%s:%d %s]", source, line, funcName))
		} else if hasSource && hasFunc {
			prefixes = append(prefixes, fmt.Sprintf("[%s %s]", source, funcName))
		} else if hasSource && hasLine {
			prefixes = append(prefixes, fmt.Sprintf("[%s:%d]", source, line))
		} else if hasFunc && hasLine {
			prefixes = append(prefixes, fmt.Sprintf("[%s:%d]", funcName, line))
		} else if hasSource {
			prefixes = append(prefixes, fmt.Sprintf("[%s]", source))
		} else if hasFunc {
			prefixes = append(prefixes, fmt.Sprintf("[%s]", funcName))
		} else if hasLine {
			prefixes = append(prefixes, fmt.Sprintf("[%d]", line))
		} else {
			prefixes = append(prefixes, "[LUA]")
		}

		if plugin, ok := entry.Data["plugin"].(string); ok {
			prefixes = append(prefixes, fmt.Sprintf("[Plugin: %s]", plugin))
		}
	}
	formattedText := strings.Join(prefixes, " ") + " " + entry.Message + "\n"

	return []byte(
		formattedText,
	), nil
}
