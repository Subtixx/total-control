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
	if entry.Data["stack"] != nil {
		if stack, ok := entry.Data["stack"].(string); ok {
			entry.Message += "\n" + stack
		}
	}

	formattedText := ""
	if entry.Data["lua"] != nil {
		formattedText = fmt.Sprintf("[%s] [%s] [LUA]: %s\n",
			entry.Time.Format("15:04:05"),
			strings.ToUpper(entry.Level.String()),
			entry.Message,
		)
	} else {
		formattedText = fmt.Sprintf("[%s] [%s]: %s\n",
			entry.Time.Format("15:04:05"),
			strings.ToUpper(entry.Level.String()),
			entry.Message,
		)
	}

	return []byte(
		formattedText,
	), nil
}
