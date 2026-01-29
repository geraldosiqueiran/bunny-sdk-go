// Package internal provides internal utilities for the Bunny SDK.
package internal

import (
	"strings"
	"time"
)

// BunnyTime is a custom time type that handles multiple date formats
// returned by Bunny.net APIs.
type BunnyTime struct {
	time.Time
}

// Bunny.net returns dates in various formats without consistent timezone info
var bunnyTimeFormats = []string{
	time.RFC3339,
	time.RFC3339Nano,
	"2006-01-02T15:04:05",
	"2006-01-02T15:04:05.999",
	"2006-01-02T15:04:05.999999",
	"2006-01-02T15:04:05.999999999",
}

// UnmarshalJSON implements json.Unmarshaler for BunnyTime.
func (bt *BunnyTime) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), `"`)
	if s == "" || s == "null" {
		return nil
	}

	var err error
	for _, format := range bunnyTimeFormats {
		bt.Time, err = time.Parse(format, s)
		if err == nil {
			return nil
		}
	}

	return err
}

// MarshalJSON implements json.Marshaler for BunnyTime.
func (bt BunnyTime) MarshalJSON() ([]byte, error) {
	if bt.IsZero() {
		return []byte("null"), nil
	}
	return []byte(`"` + bt.Format(time.RFC3339) + `"`), nil
}
