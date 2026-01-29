package internal_test

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/geraldo/bunny-sdk-go/internal"
)

// TestBunnyTimeUnmarshalJSON tests unmarshaling various time formats
func TestBunnyTimeUnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		checkFn func(*testing.T, internal.BunnyTime)
	}{
		{
			name:    "RFC3339 format",
			input:   `"2024-01-15T10:30:45Z"`,
			wantErr: false,
			checkFn: func(t *testing.T, bt internal.BunnyTime) {
				if bt.Year() != 2024 || bt.Month() != 1 || bt.Day() != 15 {
					t.Errorf("unexpected date: %v", bt.Time)
				}
			},
		},
		{
			name:    "RFC3339Nano format",
			input:   `"2024-01-15T10:30:45.123456789Z"`,
			wantErr: false,
			checkFn: func(t *testing.T, bt internal.BunnyTime) {
				if bt.Year() != 2024 {
					t.Errorf("unexpected year: %d", bt.Year())
				}
			},
		},
		{
			name:    "format without timezone",
			input:   `"2024-01-15T10:30:45"`,
			wantErr: false,
			checkFn: func(t *testing.T, bt internal.BunnyTime) {
				if bt.Year() != 2024 {
					t.Errorf("unexpected year: %d", bt.Year())
				}
			},
		},
		{
			name:    "format with milliseconds no timezone",
			input:   `"2024-01-15T10:30:45.999"`,
			wantErr: false,
			checkFn: func(t *testing.T, bt internal.BunnyTime) {
				if bt.Year() != 2024 {
					t.Errorf("unexpected year: %d", bt.Year())
				}
			},
		},
		{
			name:    "format with microseconds no timezone",
			input:   `"2024-01-15T10:30:45.999999"`,
			wantErr: false,
			checkFn: func(t *testing.T, bt internal.BunnyTime) {
				if bt.Year() != 2024 {
					t.Errorf("unexpected year: %d", bt.Year())
				}
			},
		},
		{
			name:    "format with nanoseconds no timezone",
			input:   `"2024-01-15T10:30:45.999999999"`,
			wantErr: false,
			checkFn: func(t *testing.T, bt internal.BunnyTime) {
				if bt.Year() != 2024 {
					t.Errorf("unexpected year: %d", bt.Year())
				}
			},
		},
		{
			name:    "empty string",
			input:   `""`,
			wantErr: false,
			checkFn: func(t *testing.T, bt internal.BunnyTime) {
				if !bt.IsZero() {
					t.Error("expected zero time for empty string")
				}
			},
		},
		{
			name:    "null value",
			input:   `null`,
			wantErr: false,
			checkFn: func(t *testing.T, bt internal.BunnyTime) {
				if !bt.IsZero() {
					t.Error("expected zero time for null")
				}
			},
		},
		{
			name:    "invalid format",
			input:   `"not a date"`,
			wantErr: true,
		},
		{
			name:    "invalid JSON",
			input:   `{not valid}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var bt internal.BunnyTime
			err := json.Unmarshal([]byte(tt.input), &bt)

			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.checkFn != nil {
				tt.checkFn(t, bt)
			}
		})
	}
}

// TestBunnyTimeMarshalJSON tests marshaling to JSON
func TestBunnyTimeMarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		time    internal.BunnyTime
		want    string
		wantErr bool
	}{
		{
			name: "valid time",
			time: internal.BunnyTime{Time: time.Date(2024, 1, 15, 10, 30, 45, 0, time.UTC)},
			want: `"2024-01-15T10:30:45Z"`,
		},
		{
			name: "zero time",
			time: internal.BunnyTime{},
			want: `null`,
		},
		{
			name: "time with nanoseconds",
			time: internal.BunnyTime{Time: time.Date(2024, 1, 15, 10, 30, 45, 123456789, time.UTC)},
			want: `"2024-01-15T10:30:45Z"`, // RFC3339 doesn't preserve nanoseconds when formatting
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.time)

			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if string(got) != tt.want {
					t.Errorf("MarshalJSON() = %s, want %s", got, tt.want)
				}
			}
		})
	}
}

// TestBunnyTimeRoundtrip tests marshaling and unmarshaling
func TestBunnyTimeRoundtrip(t *testing.T) {
	// Use time without nanoseconds since RFC3339 truncates them
	originalTime := time.Date(2024, 1, 15, 10, 30, 45, 0, time.UTC)
	bt := internal.BunnyTime{Time: originalTime}

	// Marshal
	data, err := json.Marshal(bt)
	if err != nil {
		t.Fatalf("Marshal() failed: %v", err)
	}

	// Unmarshal
	var bt2 internal.BunnyTime
	if err := json.Unmarshal(data, &bt2); err != nil {
		t.Fatalf("Unmarshal() failed: %v", err)
	}

	// Compare
	if !bt2.Equal(originalTime) {
		t.Errorf("roundtrip failed: got %v, want %v", bt2.Time, originalTime)
	}
}

// TestBunnyTimeInStruct tests BunnyTime in a struct
func TestBunnyTimeInStruct(t *testing.T) {
	type TestStruct struct {
		CreatedAt internal.BunnyTime `json:"createdAt"`
		UpdatedAt internal.BunnyTime `json:"updatedAt,omitempty"`
	}

	// Test unmarshaling
	jsonData := `{"createdAt":"2024-01-15T10:30:45","updatedAt":"2024-01-15T12:00:00.123"}`
	var ts TestStruct
	if err := json.Unmarshal([]byte(jsonData), &ts); err != nil {
		t.Fatalf("Unmarshal() failed: %v", err)
	}

	if ts.CreatedAt.Year() != 2024 {
		t.Errorf("unexpected CreatedAt year: %d", ts.CreatedAt.Year())
	}
	if ts.UpdatedAt.Year() != 2024 {
		t.Errorf("unexpected UpdatedAt year: %d", ts.UpdatedAt.Year())
	}

	// Test marshaling
	data, err := json.Marshal(ts)
	if err != nil {
		t.Fatalf("Marshal() failed: %v", err)
	}

	// Verify it's valid JSON
	if !strings.Contains(string(data), "createdAt") {
		t.Error("marshaled data doesn't contain createdAt field")
	}
}

// TestBunnyTimeZeroValue tests zero value handling
func TestBunnyTimeZeroValue(t *testing.T) {
	var bt internal.BunnyTime

	if !bt.IsZero() {
		t.Error("zero BunnyTime should report IsZero() as true")
	}

	// Marshal zero time
	data, err := json.Marshal(bt)
	if err != nil {
		t.Fatalf("Marshal() failed: %v", err)
	}

	if string(data) != "null" {
		t.Errorf("expected 'null', got %s", data)
	}
}

// TestBunnyTimeFormats tests all supported formats
func TestBunnyTimeFormats(t *testing.T) {
	formats := []string{
		`"2024-01-15T10:30:45Z"`,                     // RFC3339
		`"2024-01-15T10:30:45.123456789Z"`,           // RFC3339Nano
		`"2024-01-15T10:30:45"`,                      // No timezone
		`"2024-01-15T10:30:45.999"`,                  // Milliseconds
		`"2024-01-15T10:30:45.999999"`,               // Microseconds
		`"2024-01-15T10:30:45.999999999"`,            // Nanoseconds
	}

	for _, format := range formats {
		t.Run(format, func(t *testing.T) {
			var bt internal.BunnyTime
			err := json.Unmarshal([]byte(format), &bt)
			if err != nil {
				t.Errorf("failed to unmarshal %s: %v", format, err)
			}

			if bt.IsZero() {
				t.Errorf("time should not be zero for format %s", format)
			}

			// Verify year is 2024 for all formats
			if bt.Year() != 2024 {
				t.Errorf("expected year 2024, got %d for format %s", bt.Year(), format)
			}
		})
	}
}

// TestBunnyTimeEdgeCases tests edge cases
func TestBunnyTimeEdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "empty quotes",
			input:   `""`,
			wantErr: false,
		},
		{
			name:    "whitespace",
			input:   `"   "`,
			wantErr: true,
		},
		{
			name:    "partial date",
			input:   `"2024-01-15"`,
			wantErr: true,
		},
		{
			name:    "unix timestamp",
			input:   `1705315845`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var bt internal.BunnyTime
			err := json.Unmarshal([]byte(tt.input), &bt)

			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
