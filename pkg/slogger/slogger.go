// Wrapper over the slog package

package slogger

import (
	"io"
	"reflect"
	"strconv"

	"golang.org/x/exp/slog"
)

// level is either Level type or string
// Also can handle error for logging errors.
type Logger struct {
	Logger *slog.Logger
}

type Arguments struct {
	args any
}

// Generic logger wrapper to log events to a file || stdout as JSON
func NewLogger(f io.Writer) *Logger {
	return &Logger{Logger: slog.New(slog.NewJSONHandler(f))}
}

// LogEvent will log the values out after parsings args
func (p *Logger) LogEvent(logType string, msg string, args ...any) {
	arguments := NewArgumentMapper(args).mapArguments().attributeBuilder()

	switch logType {
	case "warn":
		p.Logger.Warn(msg, arguments...)
	case "debug":
		p.Logger.Debug(msg, arguments...)
	case "info":
		p.Logger.Info(msg, arguments...)
	}
}

// Specifically logs error messages
func (p *Logger) LogError(msg string, err error, args ...any) {
	p.Logger.Error(msg, err, NewArgumentMapper(args).mapArguments().attributeBuilder()...)
}

func NewArgumentMapper(args []any) *Arguments {
	return &Arguments{args: args}
}

// Generates an argument map to be passed into the log even
// Must be even as all values will be mapped as key -> Value pairs
func (a *Arguments) mapArguments() *Arguments {
	var args []any
	if a, ok := a.args.([]any); ok {
		args = a
	}

	argMap := make(map[string]any)

	// Only 1 argument will create a key value pair of empty string
	// or if the singular argument is a map, the map will be passed on
	if len(args) == 1 {
		if isMap(args[0]) {
			return a.SetMap(args[0].(map[string]any))
		}

		addArgsToMap(args[0], "", argMap)

		return a.SetMap(argMap)
	}

	for i := 0; i < len(args)-1; i++ {
		addArgsToMap(args[i], args[i+1], argMap)
	}

	return a.SetMap(argMap)
}

func (a *Arguments) SetMap(value map[string]any) *Arguments {
	a.args = value

	return a
}

// Converts incoming attributes arguments args into slog Atttributes to be further mapped within record.setAttrsFromArgs()
func (a *Arguments) attributeBuilder() []any {
	var args map[string]any
	var ok bool

	if args, ok = a.args.(map[string]any); !ok {
		return nil
	}

	if len(args) == 0 {
		return nil
	}

	var attrs []any

	for k, v := range args {
		attrs = append(attrs, slog.Attr{Key: k, Value: slog.AnyValue(v)})
	}

	return attrs
}

// Handle alternate types passed in as keys
func addArgsToMap(key, value any, argMap map[string]any) {
	switch val := key.(type) {
	case string:
		argMap[val] = value
	case int:
		argMap[strconv.Itoa(val)] = value
	case bool:
		argMap[strconv.FormatBool(val)] = value
	}
}

func isMap(value any) bool {
	return reflect.ValueOf(value).Kind() == reflect.Map
}
