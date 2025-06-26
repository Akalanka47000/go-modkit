// Package enums provides a simpler way to define and use type safe enums in Go.
// This is by no means a standard implementation but more of an opinionated way to make things easier.
package enums

import (
	"fmt"
	"reflect"
	"slices"
	"strings"
	"unsafe"
)

type Enum[V comparable] struct {
	Values      []V // List of all value values for this enum
	reflectType reflect.Type
}

type String = Enum[string]
type Int = Enum[int]
type Int32 = Enum[int32]
type Int64 = Enum[int64]
type Float32 = Enum[float32]
type Float64 = Enum[float64]
type Bool = Enum[bool]

// Create a new enum instance from a given struct.
//
// Optionally, you can pass in a boolean to indicate whether the enum values should be treated as lowercase strings if they are automatically derived from struct field names.
// This is only applicable for string enums.
func New[T any](value T, lowercase ...bool) T {
	rv := reflect.ValueOf(&value)
	re := rv.Elem()

	for i := range re.NumField() {
		field := re.Field(i)
		if strings.HasPrefix(field.Type().Name(), "Enum[") {
			continue
		}
		if field.Kind() == reflect.Ptr {
			panic(fmt.Errorf("enum fields cannot be pointers. Please use a value type instead of a pointer for field %s", re.Type().Field(i).Name))
		}
		if field.IsZero() {
			if field.CanSet() && field.Kind() == reflect.String {
				key := re.Type().Field(i).Name
				if len(lowercase) > 0 && lowercase[0] {
					key = strings.ToLower(key)
				}
				field.SetString(key)
			}
		}
		re.FieldByName("Values").Set(reflect.Append(re.FieldByName("Values"), field))
	}

	reflectTypeField := re.FieldByName("reflectType")

	reflect.NewAt(reflectTypeField.Type(), unsafe.Pointer(reflectTypeField.UnsafeAddr())).
		Elem().
		Set(reflect.ValueOf(reflect.TypeOf(value)))

	return value
}

// Returns a boolean indicating whether the provided value is a valid enum value
func (e Enum[V]) IsValid(value V) bool {
	return slices.Contains(e.Values, value)
}

// Returns an error if the provided value is not a valid enum value or otherwise returns nil
func (e Enum[V]) Validate(value V) error {
	if !e.IsValid(value) {
		return fmt.Errorf("invalid value for type %s: %v. Valid values include: %v", strings.TrimSuffix(e.reflectType.Name(), "s"), value, e.Values)
	}
	return nil
}
