package log

import (
	"bytes"
	golog "log"
	"reflect"
	"testing"
)

var bufOut, bufErr bytes.Buffer
var testLogger = &Logger{
	out:   golog.New(&bufOut, "", 0), // nil flag for empty datetime prefix
	err:   golog.New(&bufErr, "", 0),
	debug: false,
}

func TestModificationsString(t *testing.T) {
	tableTests := []struct {
		ms             Modifications // modification list
		expectedString string        // string representation
	}{
		{Modifications{1, 2, 31, 35}, "\x1b[1;2;31;35m"},                             // some hardcoded
		{Modifications{Default}, "\x1b[0m"},                                          // closing wrapper
		{Modifications{Bold, SemiBright, RedChar, RedBackground}, "\x1b[1;2;31;41m"}, // some set via variables
	}

	for _, tt := range tableTests {
		if msString := tt.ms.String(); msString != tt.expectedString {
			t.Errorf("Expected %q, got %q", tt.expectedString, msString)
		}
	}
}

// look to join this and next func
func TestWrap(t *testing.T) {
	tableTests := []struct {
		s       string        // string to wrap
		ms      Modifications // modifications
		wrapped string        // wrapped
	}{
		{"foo", nil, "foo"}, // no modifactions
		{"bar", Modifications{Bold, AquamarineChar, Blink, YellowBackground}, "\x1b[1;36;5;43mbar\x1b[0m"}, // some set of mods
	}

	for _, tt := range tableTests {
		if wrapped := Wrap(tt.s, tt.ms...); wrapped != tt.wrapped {
			t.Errorf("Expected %q, got %q", tt.wrapped, wrapped)
		}
	}
}

func TestWrapShortcuts(t *testing.T) {
	tableTests := []struct {
		shortcut string // shortcut for wrapping
		expected string // wrapped
	}{
		{debugPrefix, "\x1b[1;37m[DEBUG]\x1b[0m\t"},
		{errorPrefix, "\x1b[1;31;5m[ERROR]\x1b[0m\t"},
		{infoPrefix, "\x1b[1;32m[INFO]\x1b[0m\t"},
		{warningPrefix, "\x1b[1;33;5m[WARN]\x1b[0m\t"},
	}

	for _, tt := range tableTests {
		if tt.shortcut != tt.expected {
			t.Errorf("Expected %q, got %q", tt.expected, tt.shortcut)
		}
	}
}

func TestSetDebug(t *testing.T) {
	for _, debug := range []bool{true, false} {
		if testLogger.SetDebug(debug); testLogger.debug != debug {
			t.Errorf("Expected \"%t\", got \"%t\"", debug, testLogger.debug)
		}
	}
}

func TestLoggerPrint(t *testing.T) {
	type tableTest struct {
		methodName     string        // logger method
		args           []interface{} // args to logger method
		debug          bool          // logger printer debug mode
		expectedBuffer *bytes.Buffer // logger's buffer to print
		expectedOutput string        // expected buffer's output
	}
	tableTests := []tableTest{
		{"Error", []interface{}{"errorcode:100500"}, false, &bufErr, "\x1b[1;31;5m[ERROR]\x1b[0m\terrorcode:100500\n"},
		{"Errorf", []interface{}{"FOO%s:%dBAR", "errorcode", 100500}, false, &bufErr, "\x1b[1;31;5m[ERROR]\x1b[0m\tFOOerrorcode:100500BAR\n"},

		{"Info", []interface{}{"infocode:100500"}, false, &bufOut, "\x1b[1;32m[INFO]\x1b[0m\tinfocode:100500\n"},
		{"Infof", []interface{}{"FOO%s:%dBAR", "infocode", 100500}, false, &bufOut, "\x1b[1;32m[INFO]\x1b[0m\tFOOinfocode:100500BAR\n"},

		{"Warn", []interface{}{"warncode:100500"}, false, &bufErr, "\x1b[1;33;5m[WARN]\x1b[0m\twarncode:100500\n"},
		{"Warnf", []interface{}{"FOO%s:%dBAR", "warncode", 100500}, false, &bufErr, "\x1b[1;33;5m[WARN]\x1b[0m\tFOOwarncode:100500BAR\n"},

		{"Debug", []interface{}{"debugcode:100500"}, false, &bufErr, ""},
		{"Debugf", []interface{}{"FOO%s:%dBAR", "debugcode", 100500}, false, &bufErr, ""},

		{"Debug", []interface{}{"debugcode:100500"}, true, &bufErr, "\x1b[1;37m[DEBUG]\x1b[0m\tdebugcode:100500\n"},
		{"Debugf", []interface{}{"FOO%s:%dBAR", "debugcode", 100500}, true, &bufErr, "\x1b[1;37m[DEBUG]\x1b[0m\tFOOdebugcode:100500BAR\n"},
	}

	for _, tt := range tableTests {
		// set mode
		testLogger.SetDebug(tt.debug)
		// prepare and convert to reflect
		method := reflect.ValueOf(testLogger).MethodByName(tt.methodName)
		args := []reflect.Value{}
		for _, a := range tt.args {
			args = append(args, reflect.ValueOf(a))
		}
		// call Logger print method
		method.Call(args)
		// Check buffer output
		if bufContent := tt.expectedBuffer.String(); bufContent != tt.expectedOutput {
			t.Errorf("Expected %q, got %q", tt.expectedOutput, bufContent)
		}
		// reset buffer to future tests
		tt.expectedBuffer.Reset()
	}
}
