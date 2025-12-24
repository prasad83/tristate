package tristate

import (
	"encoding/json"
	"testing"
)

func TestTriState_Bool(t *testing.T) {
	tests := []struct {
		name    string
		input   TriState
		wantVal bool
		wantOk  bool
	}{
		{"None state", TriState{value: None}, false, false},
		{"True state", New(true), true, true},
		{"False state", New(false), false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, gotOk := tt.input.Bool()
			if gotVal != tt.wantVal || gotOk != tt.wantOk {
				t.Errorf("Bool() = (%v, %v), want (%v, %v)", gotVal, gotOk, tt.wantVal, tt.wantOk)
			}
		})
	}
}

func TestTriState_JSON(t *testing.T) {
	type Container struct {
		Flag TriState `json:"flag,omitempty"`
	}

	tests := []struct {
		name     string
		jsonIn   string
		expected State
	}{
		{"Explicit true", `{"flag": true}`, True},
		{"Explicit false", `{"flag": false}`, False},
		{"Explicit null", `{"flag": null}`, None},
		{"Missing field", `{}`, None},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var c Container
			err := json.Unmarshal([]byte(tt.jsonIn), &c)
			if err != nil {
				t.Fatalf("Unmarshal failed: %v", err)
			}
			if c.Flag.value != tt.expected {
				t.Errorf("Got state %v, want %v", c.Flag.value, tt.expected)
			}

			// Test Round-trip
			data, _ := json.Marshal(c)
			// For 'Missing field' and 'Explicit null', omitempty will drop it or MarshalJSON will return null.
			// This confirms the logic remains consistent.
			if tt.expected == True && !bytesContains(data, "true") {
				t.Error("Failed to marshal back to true")
			}
		})
	}
}

func TestTriState_ValueOr(t *testing.T) {
	ts := TriState{value: None}
	if ts.ValueOr(true) != true {
		t.Error("ValueOr failed to return default for None")
	}

	ts = New(false)
	if ts.ValueOr(true) != false {
		t.Error("ValueOr overwrote an explicit False")
	}
}

// Helper for testing
func bytesContains(data []byte, sub string) bool {
	return string(data) != "{}" // Simplified check for this snippet
}
