// Wrapper over the slog package

package slogger

import (
	"io"
	"log"
	"strconv"

	"golang.org/x/exp/slog"
)

// level is either Level type or string
// Also can handle error for logging errors.
type Logger struct {
	Logger *slog.Logger
}

// Generic logger wrapper to log events to a file || stdout as JSON
func NewLogger(f io.Writer) *Logger {
	return &Logger{Logger: slog.New(slog.NewJSONHandler(f))}
}

// Generates an argument map to be passed into the log even
// Must be even as all values will be mapped as key -> Value pairs
// Any extraneous fields will be placed in misc catagory
func Arguments(args ...any) map[string]any {
	argMap := make(map[string]any)

	// Handle alternate types passed in as keys
	for i := 0; i < len(args)-1; i++ {
		switch val := args[i].(type) {
		case string:
			argMap[val] = args[i+1]
		case int:
			argMap[strconv.Itoa(val)] = args[i+1]
		case bool:
			argMap[strconv.FormatBool(val)] = args[i+1]
		}
	}

	return argMap
}

// LogEvent will log the values out after parsings args
func (p *Logger) LogEvent(logType string, msg string, args ...any) {
	switch logType {
	case "warn":
		p.Logger.Warn(msg, attributeBuilder(args)...)
	case "debug":
		p.Logger.Debug(msg, attributeBuilder(args)...)
	case "info":
		p.Logger.Info(msg, attributeBuilder(args)...)
	}
}

// Specifically logs error messages
func (p *Logger) LogError(msg string, err error, args ...any) {
	p.Logger.Error(msg, err, attributeBuilder(args)...)
}

// Converts incoming attributes values into slog Atttributes to be further mapped within record.setAttrsFromArgs()
func attributeBuilder(attributes []any) []any {
	if len(attributes) == 0 {
		return nil
	}

	log.Println(len(attributes))

	attrs := make([]any, len(attributes)-1)

	for _, attr := range attributes {
		if m, ok := attr.(map[string]any); ok {
			for k, v := range m {
				attrs = append(attrs, slog.Attr{Key: k, Value: slog.AnyValue(v)})
			}
		}
	}

	return attrs
}
