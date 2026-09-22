// Package value retains input values and converts them to display text held
// in caller-owned storage.
//
// Nil values, including typed nil references, and empty strings, slices, or
// arrays produce an empty string so format packages can apply their own
// placeholders. Byte slices are treated as text. Values implementing error use
// Error, including values that also implement fmt.Stringer. Other fmt.Stringer
// values use String. Remaining values use their fmt.Sprint representation.
// Applications that need a different representation can configure a transformer
// or convert the value before passing it to a Table or Stream.
package value

import "math"

type tagString *byte
type tagBytes *byte
type tagInt uint8
type tagInt8 uint8
type tagInt16 uint8
type tagInt32 uint8
type tagInt64 uint8
type tagUint uint8
type tagUint8 uint8
type tagUint16 uint8
type tagUint32 uint8
type tagUint64 uint8
type tagUintptr uint8
type tagFloat32 uint8
type tagFloat64 uint8
type tagBool uint8

// Value retains primitive values without boxing.
// Its private dynamic types distinguish original Go types while keeping the
// representation to three words.
// Pointer-shaped string and byte tags keep borrowed data reachable by the GC.
type Value struct {
	_   [0]func()
	num uint64
	any any
}

// String retains a string value without boxing.
func String(v string) Value {
	return packString(v)
}

// Bytes retains a []byte value without boxing.
func Bytes(v []byte) Value {
	return packBytes(v)
}

// Int retains an int value without boxing.
func Int(v int) Value {
	// #nosec G115 -- The signed bit pattern is restored when read or formatted.
	return Value{
		num: uint64(v),
		any: tagInt(0),
	}
}

// Int8 retains an int8 value without boxing.
func Int8(v int8) Value {
	// #nosec G115 -- The signed bit pattern is restored when read or formatted.
	return Value{
		num: uint64(v),
		any: tagInt8(0),
	}
}

// Int16 retains an int16 value without boxing.
func Int16(v int16) Value {
	// #nosec G115 -- The signed bit pattern is restored when read or formatted.
	return Value{
		num: uint64(v),
		any: tagInt16(0),
	}
}

// Int32 retains an int32 value without boxing.
func Int32(v int32) Value {
	// #nosec G115 -- The signed bit pattern is restored when read or formatted.
	return Value{
		num: uint64(v),
		any: tagInt32(0),
	}
}

// Int64 retains an int64 value without boxing.
func Int64(v int64) Value {
	// #nosec G115 -- The signed bit pattern is restored when read or formatted.
	return Value{
		num: uint64(v),
		any: tagInt64(0),
	}
}

// Uint retains a uint value without boxing.
func Uint(v uint) Value {
	return Value{
		num: uint64(v),
		any: tagUint(0),
	}
}

// Uint8 retains a uint8 value without boxing.
func Uint8(v uint8) Value {
	return Value{
		num: uint64(v),
		any: tagUint8(0),
	}
}

// Uint16 retains a uint16 value without boxing.
func Uint16(v uint16) Value {
	return Value{
		num: uint64(v),
		any: tagUint16(0),
	}
}

// Uint32 retains a uint32 value without boxing.
func Uint32(v uint32) Value {
	return Value{
		num: uint64(v),
		any: tagUint32(0),
	}
}

// Uint64 retains a uint64 value without boxing.
func Uint64(v uint64) Value {
	return Value{
		num: uint64(v),
		any: tagUint64(0),
	}
}

// Uintptr retains a uintptr value without boxing.
func Uintptr(v uintptr) Value {
	return Value{
		num: uint64(v),
		any: tagUintptr(0),
	}
}

// Float32 retains a float32 value without boxing.
func Float32(v float32) Value {
	return Value{
		num: uint64(math.Float32bits(v)),
		any: tagFloat32(0),
	}
}

// Float64 retains a float64 value without boxing.
func Float64(v float64) Value {
	return Value{
		num: math.Float64bits(v),
		any: tagFloat64(0),
	}
}

// Bool retains a boolean without boxing.
func Bool(v bool) Value {
	var num uint64
	if v {
		num = 1
	}
	return Value{
		num: num,
		any: tagBool(0),
	}
}

// Any retains an arbitrary value, including its original type.
func Any(v any) Value {
	return Value{
		any: v,
	}
}

// AsString returns the stored string without boxing.
// It panics if the value is not string.
func (o Value) AsString() string {
	if _, ok := o.any.(tagString); ok {
		return unpackString(o)
	}
	return o.any.(string)
}

// AsBytes returns the stored []byte without boxing.
// It panics if the value is not []byte.
func (o Value) AsBytes() []byte {
	if _, ok := o.any.(tagBytes); ok {
		return unpackBytes(o)
	}
	return o.any.([]byte)
}

// AsInt returns the stored int without boxing.
// It panics if the value is not int.
func (o Value) AsInt() int {
	if _, ok := o.any.(tagInt); ok {
		// #nosec G115 -- Restore the original type from its retained bits.
		return int(o.num)
	}
	return o.any.(int)
}

// AsInt8 returns the stored int8 without boxing.
// It panics if the value is not int8.
func (o Value) AsInt8() int8 {
	if _, ok := o.any.(tagInt8); ok {
		// #nosec G115 -- Restore the original type from its retained bits.
		return int8(o.num)
	}
	return o.any.(int8)
}

// AsInt16 returns the stored int16 without boxing.
// It panics if the value is not int16.
func (o Value) AsInt16() int16 {
	if _, ok := o.any.(tagInt16); ok {
		// #nosec G115 -- Restore the original type from its retained bits.
		return int16(o.num)
	}
	return o.any.(int16)
}

// AsInt32 returns the stored int32 without boxing.
// It panics if the value is not int32.
func (o Value) AsInt32() int32 {
	if _, ok := o.any.(tagInt32); ok {
		// #nosec G115 -- Restore the original type from its retained bits.
		return int32(o.num)
	}
	return o.any.(int32)
}

// AsInt64 returns the stored int64 without boxing.
// It panics if the value is not int64.
func (o Value) AsInt64() int64 {
	if _, ok := o.any.(tagInt64); ok {
		// #nosec G115 -- Restore the original type from its retained bits.
		return int64(o.num)
	}
	return o.any.(int64)
}

// AsUint returns the stored uint without boxing.
// It panics if the value is not uint.
func (o Value) AsUint() uint {
	if _, ok := o.any.(tagUint); ok {
		// #nosec G115 -- Restore the original type from its retained bits.
		return uint(o.num)
	}
	return o.any.(uint)
}

// AsUint8 returns the stored uint8 without boxing.
// It panics if the value is not uint8.
func (o Value) AsUint8() uint8 {
	if _, ok := o.any.(tagUint8); ok {
		// #nosec G115 -- Restore the original type from its retained bits.
		return uint8(o.num)
	}
	return o.any.(uint8)
}

// AsUint16 returns the stored uint16 without boxing.
// It panics if the value is not uint16.
func (o Value) AsUint16() uint16 {
	if _, ok := o.any.(tagUint16); ok {
		// #nosec G115 -- Restore the original type from its retained bits.
		return uint16(o.num)
	}
	return o.any.(uint16)
}

// AsUint32 returns the stored uint32 without boxing.
// It panics if the value is not uint32.
func (o Value) AsUint32() uint32 {
	if _, ok := o.any.(tagUint32); ok {
		// #nosec G115 -- Restore the original type from its retained bits.
		return uint32(o.num)
	}
	return o.any.(uint32)
}

// AsUint64 returns the stored uint64 without boxing.
// It panics if the value is not uint64.
func (o Value) AsUint64() uint64 {
	if _, ok := o.any.(tagUint64); ok {
		return uint64(o.num)
	}
	return o.any.(uint64)
}

// AsUintptr returns the stored uintptr without boxing.
// It panics if the value is not uintptr.
func (o Value) AsUintptr() uintptr {
	if _, ok := o.any.(tagUintptr); ok {
		// #nosec G115 -- Restore the original type from its retained bits.
		return uintptr(o.num)
	}
	return o.any.(uintptr)
}

// AsFloat32 returns the stored float32 without boxing.
// It panics if the value is not float32.
func (o Value) AsFloat32() float32 {
	if _, ok := o.any.(tagFloat32); ok {
		// #nosec G115 -- Float32 stores exactly 32 bits.
		return math.Float32frombits(uint32(o.num))
	}
	return o.any.(float32)
}

// AsFloat64 returns the stored float64 without boxing.
// It panics if the value is not float64.
func (o Value) AsFloat64() float64 {
	if _, ok := o.any.(tagFloat64); ok {
		return math.Float64frombits(o.num)
	}
	return o.any.(float64)
}

// AsBool returns the stored bool without boxing.
// It panics if the value is not bool.
func (o Value) AsBool() bool {
	if _, ok := o.any.(tagBool); ok {
		return o.num != 0
	}
	return o.any.(bool)
}

// AsAny restores the original type and value. Converting the result to an
// interface can allocate.
func (o Value) AsAny() any {
	switch o.any.(type) {
	case tagString:
		return unpackString(o)
	case tagBytes:
		return unpackBytes(o)
	case tagInt:
		// #nosec G115 -- Restore the original type from its retained bits.
		return int(o.num)
	case tagInt8:
		// #nosec G115 -- Restore the original type from its retained bits.
		return int8(o.num)
	case tagInt16:
		// #nosec G115 -- Restore the original type from its retained bits.
		return int16(o.num)
	case tagInt32:
		// #nosec G115 -- Restore the original type from its retained bits.
		return int32(o.num)
	case tagInt64:
		// #nosec G115 -- Restore the original type from its retained bits.
		return int64(o.num)
	case tagUint:
		// #nosec G115 -- Restore the original type from its retained bits.
		return uint(o.num)
	case tagUint8:
		// #nosec G115 -- Restore the original type from its retained bits.
		return uint8(o.num)
	case tagUint16:
		// #nosec G115 -- Restore the original type from its retained bits.
		return uint16(o.num)
	case tagUint32:
		// #nosec G115 -- Restore the original type from its retained bits.
		return uint32(o.num)
	case tagUint64:
		// #nosec G115 -- Restore the original type from its retained bits.
		return uint64(o.num)
	case tagUintptr:
		// #nosec G115 -- Restore the original type from its retained bits.
		return uintptr(o.num)
	case tagFloat32:
		// #nosec G115 -- Float32 stores exactly 32 bits.
		return math.Float32frombits(uint32(o.num))
	case tagFloat64:
		return math.Float64frombits(o.num)
	case tagBool:
		return o.num != 0
	default:
		return o.any
	}
}
