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
	return []byte(
		fmt.Sprintf("[%s] [%s]: %s\n",
			entry.Time.Format("15:04:05"),
			strings.ToUpper(entry.Level.String()),
			entry.Message,
		),
	), nil
}
