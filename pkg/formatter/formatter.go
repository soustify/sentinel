package logger

import (
	"fmt"
	"github.com/mgutz/ansi"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

type (
	customFormatter struct {
		githubRef             string
		applicationProperties *ApplicationProperties
	}
	ApplicationProperties struct {
		ApplicationName  string
		BranchName       string
		OrganizationName string
	}
)

func NewFormatter(properties *ApplicationProperties) logrus.Formatter {

	if properties != nil && properties.ApplicationName != "" && properties.OrganizationName != "" && properties.BranchName != "" {
		return &customFormatter{
			applicationProperties: properties,
			githubRef:             fmt.Sprintf("https://github.com/%s/%s/blob/%s", properties.OrganizationName, properties.ApplicationName, properties.BranchName),
		}
	}
	return &customFormatter{}
}

func (f *customFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var file string
	var line int

	if entry.HasCaller() {
		if f.applicationProperties != nil && f.applicationProperties.ApplicationName != "" && f.applicationProperties.OrganizationName != "" && f.applicationProperties.BranchName != "" {
			file = fmt.Sprintf("%s%s", f.githubRef, strings.Split(entry.Caller.File, f.applicationProperties.ApplicationName)[1])
		} else {
			file = entry.Caller.File
		}
		line = entry.Caller.Line
	}

	// Define colors for different log levels
	var levelColor string
	switch entry.Level {
	case logrus.InfoLevel:
		levelColor = ansi.ColorCode("blue")
	case logrus.WarnLevel:
		levelColor = ansi.ColorCode("yellow")
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		levelColor = ansi.ColorCode("red")
	default:
		levelColor = ansi.ColorCode("white")
	}

	resetColor := ansi.ColorCode("reset")
	pathColor := ansi.ColorCode("green")

	var pattern string

	if f.githubRef != "" {
		pattern = "[%s%s%s] [%s] %s (%s%s#L%d%s)\n"
	} else {
		pattern = "[%s%s%s] [%s] %s (%s%s:%d%s)\n"
	}

	logMessage := fmt.Sprintf(pattern,
		levelColor, strings.ToUpper(entry.Level.String()), resetColor,
		entry.Time.Format(time.DateTime),
		entry.Message,
		pathColor, file, line, resetColor)

	return []byte(logMessage), nil
}
