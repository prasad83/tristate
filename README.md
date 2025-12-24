# Go Tri-State

A type-safe, idiomatic Go package for handling variables that can be **True**, **False**, or **None** (unspecified).

## Why not `*bool`?

Using pointers to booleans (`*bool`) is a common anti-pattern in Go that leads to several issues:

* **Ambiguity:** Does `nil` mean "not set," "inherit," or "ignore"?
* **Safety:** Requires constant nil-checks to avoid runtime panics.
* **Readability:** `if flag != nil && *flag` is significantly less clear than `if flag.IsTrue()`.
* **JSON Complexity:** Harder to control the difference between a missing key and a `null` value.

This package provides a `TriState` type that treats "None" as a first-class citizen.

## Installation

```bash
go get github.com/prasad83/tristate

```

## Core Concepts

The `TriState` type is built on an enum-backed struct where the **zero-value is always `None**`.

### 1. State Definitions

* `None`: The default state. Represents "not provided," "unset," or "use default."
* `True`: Explicitly set to true.
* `False`: Explicitly set to false.

### 2. The Comma-OK Idiom

To maintain safety, converting a `TriState` back to a standard `bool` follows Go’s `value, ok` pattern. This prevents developers from accidentally treating a `None` value as `false`.

---

## Usage

### Basic Initialization

```go
import "github.com/prasad83/tristate"

var flag tristate.TriState // Defaults to None

flag = tristate.New(true)  // Set to True
flag.IsNone()              // false
flag.IsTrue()              // true

```

### Safe Boolean Access

```go
// Use ValueOr to provide a fallback logic
isActive := flag.ValueOr(true) 

// Or use the explicit check
if v, ok := flag.Bool(); ok {
    fmt.Printf("Explicitly set to: %v", v)
} else {
    fmt.Println("Value was not provided")
}

```

### API & JSON Integration

The `TriState` type is specifically designed for API boundaries. It handles `null` and missing keys gracefully.

```go
type FeatureConfig struct {
    BetaEnabled tristate.TriState `json:"beta_enabled,omitempty"`
}

// JSON: {"beta_enabled": true}  -> State: True
// JSON: {"beta_enabled": false} -> State: False
// JSON: {"beta_enabled": null}  -> State: None
// JSON: {}                      -> State: None

```

---

## Technical Design Details

| Operation | Internal State | Method `IsTrue()` | Method `IsNone()` |
| --- | --- | --- | --- |
| `var t TriState` | `0 (None)` | `false` | `true` |
| `New(true)` | `2 (True)` | `true` | `false` |
| `New(false)` | `1 (False)` | `false` | `false` |

### Key Benefits

* **Zero Allocations:** The struct wraps a `uint8`, making it extremely lightweight.
* **No Panics:** Eliminated pointer dereferencing risks.
* **Readable API:** Focuses on intent (e.g., `IsFalse()`) rather than implementation.
