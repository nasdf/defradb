// Copyright 2023 Democratized Data Foundation
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
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type ValueToFieldKindTestCase struct {
	name   string
	kind   FieldKind
	values []any
	expect any
}

var valueToFieldKindTestCases = []ValueToFieldKindTestCase{
	{
		name:   "DocKey",
		kind:   FieldKind_DocKey,
		values: []any{"bae-2edb7fdd-cad7-5ad4-9c7d-6920245a96ed"},
		expect: DocKey{
			version: 1,
			uuid:    uuid.UUID{0x2e, 0xdb, 0x7f, 0xdd, 0xca, 0xd7, 0x5a, 0xd4, 0x9c, 0x7d, 0x69, 0x20, 0x24, 0x5a, 0x96, 0xed},
		},
	},
	{
		name:   "Bool",
		kind:   FieldKind_BOOL,
		values: []any{"TRUE", "1"},
		expect: true,
	},
	{
		name:   "[Bool]",
		kind:   FieldKind_BOOL_ARRAY,
		values: []any{[]any{true, "0"}},
		expect: []any{true, false},
	},
	{
		name:   "[Bool] = nil",
		kind:   FieldKind_BOOL_ARRAY,
		expect: nil,
	},
	{
		name:   "Int",
		kind:   FieldKind_INT,
		values: []any{float32(21), float64(21), json.Number("21"), "21"},
		expect: int64(21),
	},
	{
		name:   "[Int]",
		kind:   FieldKind_INT_ARRAY,
		values: []any{[]any{float32(1), float64(2), json.Number("3"), "4"}},
		expect: []any{int64(1), int64(2), int64(3), int64(4)},
	},
	{
		name:   "[Int] = nil",
		kind:   FieldKind_INT_ARRAY,
		expect: nil,
	},
	{
		name:   "Float",
		kind:   FieldKind_FLOAT,
		values: []any{float32(3), uint(3), uint8(3), uint16(3), uint32(3), uint64(3), int(3), int8(3), int16(3), int32(3), int64(3), json.Number("3"), "3"},
		expect: float64(3),
	},
	{
		name:   "[Float]",
		kind:   FieldKind_FLOAT_ARRAY,
		values: []any{[]any{float32(1), uint(2), uint8(3), uint16(4), uint32(5), uint64(6), int(7), int8(8), int16(9), int32(10), int64(11), json.Number("12"), "13"}},
		expect: []any{float64(1), float64(2), float64(3), float64(4), float64(5), float64(6), float64(7), float64(8), float64(9), float64(10), float64(11), float64(12), float64(13)},
	},
	{
		name:   "[Float] = nil",
		kind:   FieldKind_FLOAT_ARRAY,
		expect: nil,
	},
	{
		name:   "String",
		kind:   FieldKind_STRING,
		values: []any{[]byte("test")},
		expect: "test",
	},
	{
		name:   "[String]",
		kind:   FieldKind_STRING_ARRAY,
		values: []any{[]any{[]byte("1"), "2"}},
		expect: []any{"1", "2"},
	},
	{
		name:   "[String] = nil",
		kind:   FieldKind_STRING_ARRAY,
		expect: nil,
	},
	{
		name:   "[Bool!]",
		kind:   FieldKind_NILLABLE_BOOL_ARRAY,
		values: []any{[]any{true, "0", nil}},
		expect: []any{true, false, nil},
	},
	{
		name:   "[Bool!] = nil",
		kind:   FieldKind_NILLABLE_BOOL_ARRAY,
		expect: nil,
	},
	{
		name:   "[Int!]",
		kind:   FieldKind_NILLABLE_INT_ARRAY,
		values: []any{[]any{float32(1), float64(2), json.Number("3"), "4", nil}},
		expect: []any{int64(1), int64(2), int64(3), int64(4), nil},
	},
	{
		name:   "[Int!] = nil",
		kind:   FieldKind_NILLABLE_INT_ARRAY,
		expect: nil,
	},
	{
		name:   "[Float!]",
		kind:   FieldKind_NILLABLE_FLOAT_ARRAY,
		values: []any{[]any{float32(1), uint(2), uint8(3), uint16(4), uint32(5), uint64(6), int(7), int8(8), int16(9), int32(10), int64(11), json.Number("12"), "13", nil}},
		expect: []any{float64(1), float64(2), float64(3), float64(4), float64(5), float64(6), float64(7), float64(8), float64(9), float64(10), float64(11), float64(12), float64(13), nil},
	},
	{
		name:   "[Float!] = nil",
		kind:   FieldKind_NILLABLE_FLOAT_ARRAY,
		expect: nil,
	},
	{
		name:   "[String!]",
		kind:   FieldKind_NILLABLE_STRING_ARRAY,
		values: []any{[]any{[]byte("1"), "2", nil}},
		expect: []any{"1", "2", nil},
	},
	{
		name:   "[String!] = nil",
		kind:   FieldKind_NILLABLE_STRING_ARRAY,
		expect: nil,
	},
}

func TestValueToFieldKind(t *testing.T) {
	for _, c := range valueToFieldKindTestCases {
		t.Run(c.name, func(t *testing.T) {
			// expected value should return same type
			actual, err := ValueToFieldKind(c.expect, c.kind)
			require.NoError(t, err)
			assert.EqualValues(t, c.expect, actual)

			for _, v := range c.values {
				actual, err := ValueToFieldKind(v, c.kind)
				require.NoError(t, err)
				assert.EqualValues(t, c.expect, actual)
			}
		})
	}
}
