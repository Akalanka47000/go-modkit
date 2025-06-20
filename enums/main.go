// Package enums provides a simpler way to define and use type safe enums in Go.
// This is by no means a standard implementation but more of an opinionated way to make things easier.
package enums

import (
	"encoding/json"
	"fmt"
	"reflect"
	"slices"
	"strings"
)

type enum[T any, V comparable] struct {
	Values      T
	values      []V
	reflectType reflect.Type
}

// Create a new enum instance from a given struct. You must specifically pass in the type of the enum values (V) apart from your base type (T).
//
// Optionally, you can pass in a boolean to indicate whether the enum values should be treated as lowercase strings if they are automatically derived from struct field names.
// This is only applicable for string enums.
func New[V comparable, T any](value T, lowercase ...bool) enum[T, V] {
	var e enum[T, V]
	var zero V
	var target map[string]V

	bytes, err := json.Marshal(value)
	if err != nil {
		panic(fmt.Errorf("failed to marshal enum values: %w", err))
	}

	err = json.Unmarshal(bytes, &target)
	if err != nil {
		panic(fmt.Errorf("failed to unmarshal enum values into map: %w", err))
	}

	for k, v := range target {
		if v == zero {
			rv := reflect.ValueOf(&v).Elem()
			if rv.Kind() == reflect.Struct && rv.NumField() > 0 {
				enumTypeHolder := rv.Field(0)
				if enumTypeHolder.Kind() == reflect.Struct && enumTypeHolder.NumField() > 0 {
					rawValue := enumTypeHolder.Field(0)
					if rawValue.CanSet() && rawValue.Kind() == reflect.String {
						if len(lowercase) > 0 && lowercase[0] {
							delete(target, k)
							k = strings.ToLower(k)
						}
						rawValue.SetString(k)
						v = rv.Interface().(V)
						target[k] = v
					}
				}
			}
		}
		e.values = append(e.values, target[k])
	}

	bytes, err = json.Marshal(target)
	if err != nil {
		panic(fmt.Errorf("failed to marshal enum values after processing: %w", err))
	}

	err = json.Unmarshal(bytes, &e.Values)
	if err != nil {
		panic(fmt.Errorf("failed to unmarshal enum values into values: %w", err))
	}

	e.reflectType = reflect.TypeOf(e.Values)

	return e
}

// Returns all supported values of the enum as a slice
func (e enum[T, V]) ListValues() []V {
	return e.values
}

// Returns a boolean indicating whether the provided value is a valid enum value
func (e enum[T, V]) IsValid(value V) bool {
	return slices.Contains(e.values, value)
}

// Returns an error if the provided value is not a valid enum value or otherwise returns nil
func (e enum[T, V]) Validate(value V) error {
	if !e.IsValid(value) {
		return fmt.Errorf("invalid value for type %s: %v. Valid values include: %v", strings.TrimSuffix(e.reflectType.Name(), "s"), value, e.values)
	}
	return nil
}
