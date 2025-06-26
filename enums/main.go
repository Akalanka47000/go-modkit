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
	values      []V
	reflectType reflect.Type
}

// Create a new enum instance from a given struct.
//
// Optionally, you can pass in a boolean to indicate whether the enum values should be treated as lowercase strings if they are automatically derived from struct field names.
// This is only applicable for string enums.
func New[T any](value T, opts ...optionBuilder) T {
	rv := reflect.ValueOf(&value)
	re := rv.Elem()

	var options options

	for _, opt := range opts {
		opt(&options)
	}

	for i := 0; i < re.NumField(); i++ {
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
				if options.lowercase {
					key = strings.ToLower(key)
				}
				if options.uppercase {
					key = strings.ToUpper(key)
				}
				field.SetString(key)
			}
		}
		valuesField := re.FieldByName("values")
		unsafeValuesField := reflect.NewAt(valuesField.Type(), unsafe.Pointer(valuesField.UnsafeAddr())).Elem()
		unsafeValuesField.Set(reflect.Append(unsafeValuesField, field))
	}

	reflectTypeField := re.FieldByName("reflectType")

	reflect.NewAt(reflectTypeField.Type(), unsafe.Pointer(reflectTypeField.UnsafeAddr())).
		Elem().
		Set(reflect.ValueOf(reflect.TypeOf(value)))

	return value
}

// Returns a slice of all valid enum values.
func (e Enum[V]) Values() []V {
	return e.values
}

// Returns a boolean indicating whether the provided value is a valid enum value
func (e Enum[V]) IsValid(value V) bool {
	return slices.Contains(e.values, value)
}

// Returns an error if the provided value is not a valid enum value or otherwise returns nil
func (e Enum[V]) Validate(value V) error {
	if !e.IsValid(value) {
		return fmt.Errorf("invalid value for type %s: %v. Valid values include: %v", strings.TrimSuffix(e.reflectType.Name(), "s"), value, e.values)
	}
	return nil
}
