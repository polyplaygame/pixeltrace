package tools

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTimezoneAndOffset(t *testing.T) {
	tests := []struct {
		locationName   string
		expectedTZ     string
		expectedOffset float64
		expectError    bool
	}{
		{"Asia/Kolkata", "IST", 5.5, false},
		{"America/New_York", "EST", -5, false},
		{"Europe/London", "GMT", 0, false},
		{"Invalid/Location", "", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.locationName, func(t *testing.T) {
			tz, offset, err := TimezoneAndHourOffset(tt.locationName)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedTZ, tz)
				assert.Equal(t, tt.expectedOffset, offset)
			}
		})
	}
}
