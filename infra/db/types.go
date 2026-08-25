package db

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type JSONB[T any] struct {
	Data T
}

func NewJSONB[T any](data T) JSONB[T] {
	return JSONB[T]{Data: data}
}

func (j JSONB[T]) Value() (driver.Value, error) {
	return json.Marshal(j.Data)
}

func (j *JSONB[T]) Scan(value any) error {
	if value == nil {
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return fmt.Errorf("failed to scan JSONB: expected []byte or string, got %T", value)
	}

	if len(bytes) == 0 {
		return nil
	}

	return json.Unmarshal(bytes, &j.Data)
}

func (j JSONB[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(j.Data)
}

func (j *JSONB[T]) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &j.Data)
}
