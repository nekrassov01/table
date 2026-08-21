// Package value converts arbitrary Go values to display text held in
// caller-owned storage.
//
// Nil values, including typed nil references, and empty strings, slices, or
// arrays produce an empty string so format packages can apply their own
// placeholders. Byte slices are treated as text. Values implementing
// fmt.Stringer or error use those representations. Other slices, arrays, maps,
// and structs use short summaries instead of recursively formatting their
// contents. Applications that need a different representation can configure a
// transformer or convert the value before passing it to a Table or Stream.
package value

import (
	"fmt"
	"reflect"
	"strconv"
)

// Number formats an integer directly for a synthetic index column, avoiding
// interface conversion through Format.
func Number(st *Store, x int64) string {
	m := st.Mark()
	st.AppendInt(x)
	return st.Since(m)
}

// Format converts v to display text. Missing values include nil, a typed
// nil, and an empty string, slice, or array. They yield the empty string.
// A byte slice yields its text; every other slice or array yields a
// summary. Text appended to st is returned as a view that remains valid until
// the store is reset.
func Format(st *Store, v any) string {
	if v == nil {
		return ""
	}
	switch x := v.(type) {
	case string:
		return x
	case int:
		return appendInt(st, int64(x))
	case int8:
		return appendInt(st, int64(x))
	case int16:
		return appendInt(st, int64(x))
	case int32:
		return appendInt(st, int64(x))
	case int64:
		return appendInt(st, x)
	case uint:
		return appendUint(st, uint64(x))
	case uint8:
		return appendUint(st, uint64(x))
	case uint16:
		return appendUint(st, uint64(x))
	case uint32:
		return appendUint(st, uint64(x))
	case uint64:
		return appendUint(st, x)
	case uintptr:
		return appendUint(st, uint64(x))
	case float32:
		return appendFloat(st, float64(x), 32)
	case float64:
		return appendFloat(st, x, 64)
	case bool:
		return strconv.FormatBool(x)
	case []byte:
		return appendBytes(st, x)
	case fmt.Stringer:
		if isTypedNil(x) {
			return ""
		}
		return x.String()
	case error:
		if isTypedNil(x) {
			return ""
		}
		return x.Error()
	default:
		return formatReflect(st, v)
	}
}

// formatReflect is the reflection-based fallback for types not covered
// by the type switch in Format.
func formatReflect(st *Store, v any) string {
	rv := reflect.ValueOf(v)
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
	switch e.Kind() {
	case reflect.Struct:
		return summaryStruct.format(e.NumField())
	case reflect.Map:
		return summaryMap.format(e.Len())
	case reflect.Slice, reflect.Array:
		if e.Len() == 0 {
			return ""
		}
		// Named byte slices are text like []byte; byte arrays remain ordinary
		// arrays.
		if e.Kind() == reflect.Slice && e.Type().Elem().Kind() == reflect.Uint8 {
			return appendBytes(st, e.Bytes())
		}
		return summaryList.format(e.Len())
	default:
		return fmt.Sprint(e.Interface())
	}
}

// resolveReference dereferences interface and pointer chains and tries
// resolveStringer at each level. A nil reference resolves to an empty string.
func resolveReference(rv reflect.Value) (reflect.Value, string, bool) {
	for rv.Kind() == reflect.Interface || rv.Kind() == reflect.Pointer {
		if rv.IsNil() {
			return rv, "", true
		}
		if s, ok := resolveStringer(rv); ok {
			return rv, s, true
		}
		rv = rv.Elem()
	}
	switch rv.Kind() {
	case reflect.Map, reflect.Slice, reflect.Chan, reflect.Func:
		if rv.IsNil() {
			return rv, "", true
		}
	}
	if s, ok := resolveStringer(rv); ok {
		return rv, s, true
	}
	return rv, "", false
}

// resolveStringer formats rv through fmt.Stringer or error when possible.
func resolveStringer(rv reflect.Value) (string, bool) {
	if !rv.CanInterface() {
		return "", false
	}
	t := rv.Type()
	if t.Implements(reflect.TypeFor[fmt.Stringer]()) {
		return rv.Interface().(fmt.Stringer).String(), true
	}
	if t.Implements(reflect.TypeFor[error]()) {
		return rv.Interface().(error).Error(), true
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
		return appendInt(st, rv.Int()), true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return appendUint(st, rv.Uint()), true
	case reflect.Float32:
		return appendFloat(st, rv.Float(), 32), true
	case reflect.Float64:
		return appendFloat(st, rv.Float(), 64), true
	}
	return "", false
}

func appendBytes(st *Store, b []byte) string {
	m := st.Mark()
	st.AppendBytes(b)
	return st.Since(m)
}

func appendInt(st *Store, x int64) string {
	m := st.Mark()
	st.AppendInt(x)
	return st.Since(m)
}

func appendUint(st *Store, x uint64) string {
	m := st.Mark()
	st.AppendUint(x)
	return st.Since(m)
}

func appendFloat(st *Store, x float64, bitSize int) string {
	m := st.Mark()
	st.AppendFloat(x, bitSize)
	return st.Since(m)
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
