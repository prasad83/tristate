// Package tristate provides a type-safe implementation of tri-state logic.
package tristate

import (
	"bytes"
	"fmt"
)

// State represents the underlying value of the TriState.
type State uint8

const (
	None  State = iota // 0: Not provided / Inherit
	False              // 1: Explicitly False
	True               // 2: Explicitly True
)

// TriState wraps the State to provide a clean API and JSON support.
type TriState struct {
	value State
}

// --- Factory Methods ---

func New(v bool) TriState {
	if v {
		return TriState{value: True}
	}
	return TriState{value: False}
}

// --- Accessors ---

func (t TriState) IsNone() bool  { return t.value == None }
func (t TriState) IsTrue() bool  { return t.value == True }
func (t TriState) IsFalse() bool { return t.value == False }

// Bool returns the boolean value and a 'valid' bit.
// If state is None, it returns (false, false).
func (t TriState) Bool() (val bool, ok bool) {
	switch t.value {
	case True:
		return true, true
	case False:
		return false, true
	default:
		return false, false
	}
}

// ValueOr returns the boolean value if set, or the provided default if None.
func (t TriState) ValueOr(defaultVal bool) bool {
	if v, ok := t.Bool(); ok {
		return v
	}
	return defaultVal
}

// --- JSON Marshaling ---

// MarshalJSON converts the TriState to true, false, or null.
func (t TriState) MarshalJSON() ([]byte, error) {
	switch t.value {
	case True:
		return []byte("true"), nil
	case False:
		return []byte("false"), nil
	default:
		return []byte("null"), nil
	}
}

// UnmarshalJSON handles incoming true, false, and null values.
func (t *TriState) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, []byte("null")) {
		t.value = None
		return nil
	}
	if bytes.Equal(data, []byte("true")) {
		t.value = True
		return nil
	}
	if bytes.Equal(data, []byte("false")) {
		t.value = False
		return nil
	}
	return fmt.Errorf("invalid tristate value: %s", string(data))
}
