package value

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
)

// Format converts a retained value to display text without reboxing it.
// Missing values yield an empty string; byte slices are text. Stored output
// remains valid until the store is reset.
func Format(st *Store, input Value) string {
	// Keep the string path small enough to inline into the output packages.
	// It needs neither conversion nor output storage.
	switch input.any.(type) {
	case tagString:
		return unpackString(input)
	default:
		return format(st, input)
	}
}

func format(st *Store, input Value) string {
	switch v := input.any.(type) {
	case string:
		return v
	case nil:
		return ""
	case tagBytes:
		return st.AppendBytes(unpackBytes(input))
	case tagInt, tagInt8, tagInt16, tagInt32, tagInt64:
		// #nosec G115 -- Restore the signed bit pattern.
		return st.AppendInt(int64(input.num))
	case tagUint, tagUint8, tagUint16, tagUint32, tagUint64, tagUintptr:
		return st.AppendUint(input.num)
	case tagFloat32:
		// #nosec G115 -- Float32 stores exactly 32 bits.
		return st.AppendFloat(float64(math.Float32frombits(uint32(input.num))), 32)
	case tagFloat64:
		return st.AppendFloat(math.Float64frombits(input.num), 64)
	case tagBool:
		return strconv.FormatBool(input.num != 0)
	case int:
		return st.AppendInt(int64(v))
	case int8:
		return st.AppendInt(int64(v))
	case int16:
		return st.AppendInt(int64(v))
	case int32:
		return st.AppendInt(int64(v))
	case int64:
		return st.AppendInt(v)
	case uint:
		return st.AppendUint(uint64(v))
	case uint8:
		return st.AppendUint(uint64(v))
	case uint16:
		return st.AppendUint(uint64(v))
	case uint32:
		return st.AppendUint(uint64(v))
	case uint64:
		return st.AppendUint(v)
	case uintptr:
		return st.AppendUint(uint64(v))
	case float32:
		return st.AppendFloat(float64(v), 32)
	case float64:
		return st.AppendFloat(v, 64)
	case bool:
		return strconv.FormatBool(v)
	case []byte:
		return st.AppendBytes(v)
	case error:
		if isTypedNil(v) {
			return ""
		}
		return v.Error()
	case fmt.Stringer:
		if isTypedNil(v) {
			return ""
		}
		return v.String()
	default:
		return formatReflect(st, reflect.ValueOf(input.any))
	}
}

// Number formats an integer directly for a synthetic index column, avoiding
// interface conversion through Format.
func Number(st *Store, x int64) string {
	return st.AppendInt(x)
}

// formatReflect is the reflection-based fallback for types not covered
// by the general type dispatch.
func formatReflect(st *Store, rv reflect.Value) string {
	if !rv.IsValid() {
		return ""
	}
	e, s, ok := resolveReference(rv)
	if ok {
		return s
	}
	if s, ok := resolveKind(st, e); ok {
		return s
	}
	if e.Kind() == reflect.Slice || e.Kind() == reflect.Array {
		if e.Len() == 0 {
			return ""
		}
		// Named byte slices are text like []byte; byte arrays remain ordinary
		// arrays.
		if e.Kind() == reflect.Slice && e.Type().Elem().Kind() == reflect.Uint8 {
			return st.AppendBytes(e.Bytes())
		}
		if s, ok := formatPrimitives(st, e); ok {
			return s
		}
	}
	return st.AppendDefault(e.Interface())
}

// formatPrimitives formats a slice or array of primitive values without
// allocating a temporary string and reports whether its element type matched.
func formatPrimitives(st *Store, rv reflect.Value) (string, bool) {
	elemType := rv.Type().Elem()
	if elemType.Implements(reflect.TypeFor[fmt.Stringer]()) || elemType.Implements(reflect.TypeFor[error]()) {
		return "", false
	}
	switch elemType.Kind() {
	case reflect.String, reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64:
	default:
		return "", false
	}
	mark := st.Mark()
	st.AppendString("[")
	for index := range rv.Len() {
		if index > 0 {
			st.AppendString(" ")
		}
		value := rv.Index(index)
		switch value.Kind() {
		case reflect.String:
			st.AppendString(value.String())
		case reflect.Bool:
			st.AppendString(strconv.FormatBool(value.Bool()))
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			st.AppendInt(value.Int())
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			st.AppendUint(value.Uint())
		case reflect.Float32:
			st.AppendFloat(value.Float(), 32)
		case reflect.Float64:
			st.AppendFloat(value.Float(), 64)
		}
	}
	st.AppendString("]")
	return st.Since(mark), true
}

// resolveReference dereferences an interface or pointer chain and resolves
// fmt.Stringer or error at each reference. A nil reference resolves to an
// empty string, while a cyclic chain returns its repeated pointer unresolved.
func resolveReference(rv reflect.Value) (reflect.Value, string, bool) {
	if rv.Kind() == reflect.Interface || rv.Kind() == reflect.Pointer {
		if rv.IsNil() {
			return rv, "", true
		}
		if s, ok := resolveStringerOrError(rv); ok {
			return rv, s, true
		}
		rv = rv.Elem()
		if rv.Kind() == reflect.Interface || rv.Kind() == reflect.Pointer {
			return resolveReferenceChain(rv)
		}
	}
	switch rv.Kind() {
	case reflect.Map, reflect.Slice, reflect.Chan, reflect.Func:
		if rv.IsNil() {
			return rv, "", true
		}
	}
	if s, ok := resolveStringerOrError(rv); ok {
		return rv, s, true
	}
	return rv, "", false
}

// resolveReferenceChain continues a chain of two or more references. A nil
// reference resolves to an empty string, while a cyclic chain returns its
// repeated pointer unresolved.
func resolveReferenceChain(rv reflect.Value) (reflect.Value, string, bool) {
	var anchor reflect.Value
	limit := 1
	distance := 0
	for rv.Kind() == reflect.Interface || rv.Kind() == reflect.Pointer {
		if rv.IsNil() {
			return rv, "", true
		}
		if s, ok := resolveStringerOrError(rv); ok {
			return rv, s, true
		}
		if rv.Kind() != reflect.Pointer {
			rv = rv.Elem()
			continue
		}
		if !anchor.IsValid() {
			anchor = rv
			distance = 1
			rv = rv.Elem()
			continue
		}
		if anchor.Type() == rv.Type() && anchor.Pointer() == rv.Pointer() {
			return rv, "", false
		}
		// Brent's algorithm moves the anchor after 1, 2, 4, ... pointer
		// steps to detect cycles without retaining the chain.
		if distance == limit {
			anchor = rv
			limit *= 2
			distance = 0
		}
		distance++
		rv = rv.Elem()
	}
	switch rv.Kind() {
	case reflect.Map, reflect.Slice, reflect.Chan, reflect.Func:
		if rv.IsNil() {
			return rv, "", true
		}
	}
	if s, ok := resolveStringerOrError(rv); ok {
		return rv, s, true
	}
	return rv, "", false
}

// resolveStringerOrError formats rv through error or fmt.Stringer when
// possible, preferring error when both are implemented.
func resolveStringerOrError(rv reflect.Value) (string, bool) {
	if !rv.CanInterface() {
		return "", false
	}
	t := rv.Type()
	if t.Implements(reflect.TypeFor[error]()) {
		return rv.Interface().(error).Error(), true
	}
	if t.Implements(reflect.TypeFor[fmt.Stringer]()) {
		return rv.Interface().(fmt.Stringer).String(), true
	}
	return "", false
}

// resolveKind formats a primitive reflect.Kind and reports whether it matched.
func resolveKind(st *Store, rv reflect.Value) (string, bool) {
	switch rv.Kind() {
	case reflect.String:
		return rv.String(), true
	case reflect.Bool:
		return strconv.FormatBool(rv.Bool()), true
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return st.AppendInt(rv.Int()), true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return st.AppendUint(rv.Uint()), true
	case reflect.Float32:
		return st.AppendFloat(rv.Float(), 32), true
	case reflect.Float64:
		return st.AppendFloat(rv.Float(), 64), true
	}
	return "", false
}

// isTypedNil reports whether v holds a typed-nil reference.
func isTypedNil(v any) bool {
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Pointer, reflect.Map, reflect.Slice, reflect.Chan, reflect.Func:
		return rv.IsNil()
	}
	return false
}
