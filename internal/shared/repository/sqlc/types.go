package sqlc

import (
	"encoding/json"
	"time"
)

// NullableTime wraps time.Time and handles nullability
type NullableTime struct {
	Time  time.Time
	Valid bool
}

// MarshalJSON ensures nullability during JSON serialization
func (nt NullableTime) MarshalJSON() ([]byte, error) {
	if nt.Valid {
		return json.Marshal(nt.Time)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON handles nullability during JSON deserialization
func (nt *NullableTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		nt.Valid = false
		return nil
	}
	err := json.Unmarshal(data, &nt.Time)
	if err != nil {
		return err
	}
	nt.Valid = true
	return nil
}
