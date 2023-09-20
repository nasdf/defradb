// Copyright 2022 Democratized Data Foundation
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

package client

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/fxamacker/cbor/v2"
)

// Value is an interface that points to a concrete Value implementation.
type Value interface {
	// Value returns the underlying primitive value.
	Value() any
	// Type returns the CRDT type for this value.
	Type() CType
	// IsDirty returns true if the value has been updated.
	IsDirty() bool
	// IsDeleted returns true if the value has been deleted.
	IsDeleted() bool
	// Delete marks the value as deleted and sets the value to nil.
	Delete()
	// Clean removes any dirty or updated marks from the value.
	Clean()
}

// WriteableValue defines a simple interface with a Bytes() method
// which is used to indicate if a Value is writeable type versus
// a composite type like a Sub-Document.
// Writeable types include simple Strings/Ints/Floats/Binary
// that can be loaded into a CRDT Register, Set, Counter, etc.
type WriteableValue interface {
	Value

	Bytes() ([]byte, error)
}

var (
	_ WriteableValue = (*PrimitiveValue)(nil)
	_ json.Marshaler = (*PrimitiveValue)(nil)
	_ cbor.Marshaler = (*PrimitiveValue)(nil)
)

// PrimitiveValue contains values of type String, Int, Float, Bool.
type PrimitiveValue struct {
	value   any
	ctype   CType
	dirty   bool
	deleted bool
}

func (p *PrimitiveValue) Value() any {
	return p.value
}

func (p *PrimitiveValue) Type() CType {
	return p.ctype
}

func (p *PrimitiveValue) IsDirty() bool {
	return p.dirty
}

func (p *PrimitiveValue) IsDeleted() bool {
	return p.deleted
}

func (p *PrimitiveValue) Delete() {
	p.value = nil
	p.deleted = true
}

func (p *PrimitiveValue) Clean() {
	p.dirty = false
	p.deleted = false
}

func (p *PrimitiveValue) Bytes() ([]byte, error) {
	return cbor.Marshal(p)
}

func (p *PrimitiveValue) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.value)
}

func (p *PrimitiveValue) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal(p.value)
}

var _ Value = (*ReferenceValue)(nil)

// ReferenceValue contains values of type Document.
type ReferenceValue struct {
	value   any
	ctype   CType
	dirty   bool
	deleted bool
}

func (r *ReferenceValue) Value() any {
	return r.value
}

func (r *ReferenceValue) Type() CType {
	return r.ctype
}

func (r *ReferenceValue) IsDirty() bool {
	return r.dirty
}

func (r *ReferenceValue) IsDeleted() bool {
	return r.deleted
}

func (r *ReferenceValue) Delete() {
	r.value = nil
	r.deleted = true
}

func (r *ReferenceValue) Clean() {
	r.dirty = false
	r.deleted = false
}

// ArrayToFieldKind converts a generic array to the expected
// FieldKind value type.
func ArrayToFieldKind(value []any, kind FieldKind, nillable bool) (any, error) {
	array := make([]any, len(value))
	for i, v := range value {
		var err error
		if !nillable || v != nil {
			array[i], err = ValueToFieldKind(v, kind)
		}
		if err != nil {
			return nil, err
		}
	}
	return array, nil
}

// ValueToFieldKind converts a generic value to the expected
// FieldKind value type.
func ValueToFieldKind(value any, kind FieldKind) (any, error) {
	switch kind {
	case FieldKind_BOOL:
		switch t := value.(type) {
		case bool:
			return t, nil
		case string:
			return strconv.ParseBool(t)
		}
	case FieldKind_BOOL_ARRAY, FieldKind_NILLABLE_BOOL_ARRAY:
		switch t := value.(type) {
		case []bool, nil:
			return t, nil
		case []any:
			return ArrayToFieldKind(t, FieldKind_BOOL, kind == FieldKind_NILLABLE_BOOL_ARRAY)
		}
	case FieldKind_INT:
		switch t := value.(type) {
		case uint, uint16, uint32, uint64, int, int16, int32, int64:
			return t, nil
		case float32:
			return int64(t), nil
		case float64:
			return int64(t), nil
		case json.Number:
			return t.Int64()
		case string:
			return strconv.ParseInt(t, 10, 64)
		}
	case FieldKind_INT_ARRAY, FieldKind_NILLABLE_INT_ARRAY:
		switch t := value.(type) {
		case []uint, []uint16, []uint32, []uint64, []int, []int16, []int32, []int64, nil:
			return t, nil
		case []any:
			return ArrayToFieldKind(t, FieldKind_INT, kind == FieldKind_NILLABLE_INT_ARRAY)
		}
	case FieldKind_FLOAT:
		switch t := value.(type) {
		case float64:
			return t, nil
		case float32:
			return float64(t), nil
		case uint:
			return float64(t), nil
		case uint8:
			return float64(t), nil
		case uint16:
			return float64(t), nil
		case uint32:
			return float64(t), nil
		case uint64:
			return float64(t), nil
		case int:
			return float64(t), nil
		case int8:
			return float64(t), nil
		case int16:
			return float64(t), nil
		case int32:
			return float64(t), nil
		case int64:
			return float64(t), nil
		case json.Number:
			return t.Float64()
		case string:
			return strconv.ParseFloat(t, 64)
		}
	case FieldKind_FLOAT_ARRAY, FieldKind_NILLABLE_FLOAT_ARRAY:
		switch t := value.(type) {
		case []float32, []float64, nil:
			return t, nil
		case []any:
			return ArrayToFieldKind(t, FieldKind_FLOAT, kind == FieldKind_NILLABLE_FLOAT_ARRAY)
		}
	case FieldKind_DATETIME:
		switch t := value.(type) {
		case string:
			val, err := time.Parse(time.RFC3339, t)
			if err != nil {
				return nil, err
			}
			return val.Format(time.RFC3339), nil
		case time.Time:
			return t.Format(time.RFC3339), nil
		}
	case FieldKind_STRING:
		switch t := value.(type) {
		case string:
			return t, nil
		case []byte:
			return string(t), nil
		}
	case FieldKind_STRING_ARRAY, FieldKind_NILLABLE_STRING_ARRAY:
		switch t := value.(type) {
		case []string, nil:
			return t, nil
		case []any:
			return ArrayToFieldKind(t, FieldKind_STRING, kind == FieldKind_NILLABLE_STRING_ARRAY)
		}
	}
	return nil, fmt.Errorf("invalid value: %v", value)
}
