package isset

import (
	"fmt"
	"time"
)

type Time struct {
	time.Time
	valid bool
}

// IsValid returns whether a value has been set
func (t *Time) IsValid() bool {
	return t.valid
}

// Set a value.
func (t *Time) Set(tm time.Time) {
	t.valid = true
	t.Time = tm
}

// Unset the variable, like setting it to nil
func (t *Time) Unset() {
	t.valid = false
}

// Get returns the contained value.
func (t *Time) Get() (time.Time, error) {
	if t.valid {
		return t.Time, nil
	}
	return time.Time{}, fmt.Errorf("runtime error: attempt to get value of Null time which is set to nil")
}

// Now returns a NTime having the Time now
func Now() Time {
	return Time{time.Now(), true}
}
