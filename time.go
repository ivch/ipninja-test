package main

import (
	"database/sql/driver"
	"strings"
	"time"
)

// Time type for json time parsing
type Time struct {
	time.Time
}

// UnmarshalJSON unmarshal my time from json, layout: 1504 (Hi)
func (ct *Time) UnmarshalJSON(b []byte) (err error) {
	ct.Time, err = time.Parse(`"2006-01-02 15:04:05"`, string(b))

	if err != nil && strings.Contains(err.Error(), "month out of range") {
		err = nil
		ct.Time = time.Unix(0, 0)
		return
	}

	return
}

// MarshalJSON marshal my time to json, layout: 1504 (Hi)
func (ct Time) MarshalJSON() ([]byte, error) {
	return []byte(ct.Format("\"2006-01-02 15:04:05\"")), nil
}

// Value get value of known type
func (ct Time) Value() (driver.Value, error) {
	return ct.Time, nil
}

// Scan convert form known type
func (ct *Time) Scan(src interface{}) error {
	ct.Time = src.(time.Time)
	return nil
}
