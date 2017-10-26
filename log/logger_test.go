package log

import (
	"testing"
)

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
		ms      Modifications //modifications
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
