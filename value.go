package table

import "github.com/nekrassov01/table/internal/value"

// Value is a body value. Its zero value is empty. Primitive constructors avoid
// interface boxing; Any retains arbitrary Go values. Values are not comparable.
type Value = value.Value

// String stores a string value, preserving its type.
func String(v string) Value {
	return value.String(v)
}

// Bytes borrows v without copying. Do not modify v until Render returns.
// AsAny returns the original contents and length with capacity equal to length.
func Bytes(v []byte) Value {
	return value.Bytes(v)
}

// Int stores an int value, preserving its type.
func Int(v int) Value {
	return value.Int(v)
}

// Int8 stores an int8 value, preserving its type.
func Int8(v int8) Value {
	return value.Int8(v)
}

// Int16 stores an int16 value, preserving its type.
func Int16(v int16) Value {
	return value.Int16(v)
}

// Int32 stores an int32 value, preserving its type.
func Int32(v int32) Value {
	return value.Int32(v)
}

// Int64 stores an int64 value, preserving its type.
func Int64(v int64) Value {
	return value.Int64(v)
}

// Uint stores a uint value, preserving its type.
func Uint(v uint) Value {
	return value.Uint(v)
}

// Uint8 stores a uint8 value, preserving its type.
func Uint8(v uint8) Value {
	return value.Uint8(v)
}

// Uint16 stores a uint16 value, preserving its type.
func Uint16(v uint16) Value {
	return value.Uint16(v)
}

// Uint32 stores a uint32 value, preserving its type.
func Uint32(v uint32) Value {
	return value.Uint32(v)
}

// Uint64 stores a uint64 value, preserving its type.
func Uint64(v uint64) Value {
	return value.Uint64(v)
}

// Uintptr stores a uintptr value, preserving its type.
func Uintptr(v uintptr) Value {
	return value.Uintptr(v)
}

// Float32 stores a float32 value, preserving its type.
func Float32(v float32) Value {
	return value.Float32(v)
}

// Float64 stores a float64 value, preserving its type.
func Float64(v float64) Value {
	return value.Float64(v)
}

// Bool stores a bool value, preserving its type.
func Bool(v bool) Value {
	return value.Bool(v)
}

// Any stores an arbitrary value with its original type. Boxing may allocate.
func Any(v any) Value {
	return value.Any(v)
}
