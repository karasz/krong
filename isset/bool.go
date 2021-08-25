package isset

import (
	"fmt"
)

type Bool struct {
	value bool
	valid bool
}

// IsValid returns whether a value has been set
func (b *Bool) IsValid() bool {
	return b.valid
}

// Set a value.
func (b *Bool) Set(nb bool) {
	b.valid = true
	b.value = nb
}

// Unset the variable, like setting it to nil
func (b *Bool) Unset() {
	b.valid = false
}

// Get returns the contained value.
func (b *Bool) Get() (bool, error) {
	if b.valid {
		return b.value, nil
	}
	return false, fmt.Errorf("runtime error: attempt to get value of Bool which is set to nil")
}

func (b *Bool) IsSetBool2Int() int {
	if b.valid {
		if b.value {
			return 1
		}
		return 0
	}
	return -1
}
func (b *Bool) Reset() {
	b.valid = false
	b.value = false
}
func (b *Bool) IsSetInt2Bool(i int) Bool {
	switch i {
	case 1:
		b.Set(true)
	case 0:
		b.Set(false)
	case -1:
		b.valid = false
		b.value = false
	}
	return *b
}
