package log

import (
	"bytes"
	golog "log"
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

func TestError(t *testing.T) {
	defer bufOut.Reset() // reset stdout buffer after test
	defer bufErr.Reset() // reset stderr buffer after test
	// write err
	testLogger.Error("errorcode:100500")
	// check io.Writer of Logger
	// error goes to .err writer
	expectedOutput := "\x1b[1;31;5m[ERROR]\x1b[0m\terrorcode:100500\n"
	if bufErrContent := bufErr.String(); bufErrContent != expectedOutput {
		t.Errorf("Expected %q, got %q", expectedOutput, bufErrContent)
	}
	// .out writer must be still empty
	if bufOutContent := bufOut.String(); bufOutContent != "" {
		t.Errorf("Expected %q, got %q", "", bufOutContent)
	}
}
