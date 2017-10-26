package log

import (
	"testing"
)

func TestModificationsString(t *testing.T) {
	tableTests := []struct {
		ms             Modifications // modification list
		expectedString string        // string representation
	}{
		{Modifications{1, 2, 31, 35}, "\x1b[1;2;31;35m"},
	}

	for _, tt := range tableTests {
		if msString := tt.ms.String(); msString != tt.expectedString {
			t.Errorf("Expected %q, got %q", tt.expectedString, msString)
		}
	}
}
