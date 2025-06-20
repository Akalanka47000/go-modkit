package enums

import (
	"encoding/json"
	"fmt"
)

type Type[T comparable] struct {
	Raw T
}

func (p Type[T]) Format(f fmt.State, c rune) {
	if f.Flag('+') {
		fmt.Fprintf(f, "%T(%v)", p.Raw, p.Raw)
	} else {
		fmt.Fprintf(f, "%v", p.Raw)
	}
}

func (p Type[T]) MarshalJSON() ([]byte, error) {
	bytes, err := json.Marshal(p.Raw)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal enum value: %w", err)
	}
	return bytes, nil
}

func (p *Type[T]) UnmarshalJSON(data []byte) error {
	var value T
	if err := json.Unmarshal(data, &value); err != nil {
		return fmt.Errorf("failed to unmarshal enum value: %w", err)
	}
	p.Raw = value
	return nil
}

// Type alises

type String = Type[string]

type Int = Type[int]

type Int32 = Type[int32]

type Int64 = Type[int64]

type Uint = Type[uint]

type Uint32 = Type[uint32]

type Uint64 = Type[uint64]

type Float32 = Type[float32]

type Float64 = Type[float64]

type Bool = Type[bool]

type Any = Type[any]
