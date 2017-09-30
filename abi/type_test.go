// Copyright 2016 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package abi

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ethereum/go-ethereum/common"
)

// typeWithoutStringer is a alias for the Type type which simply doesn't implement
// the stringer interface to allow printing type details in the tests below.
type typeWithoutStringer Type

type encodeTestCase struct {
	value    interface{}
	typeName string
	result   string
	errorMsg string
}

func runEncodeTestCase(is *assert.Assertions, tt *encodeTestCase) {
	ty, err := NewType(tt.typeName)
	is.NoError(err)

	data, err := ty.Pack(tt.value)

	if tt.result != "" {
		is.NoError(err, "packing")
		is.Equal(tt.result, hex.EncodeToString(data))
	} else if tt.errorMsg != "" {
		is.Error(err)
		is.True(strings.HasPrefix(err.Error(), tt.errorMsg),
			fmt.Sprintf("Error message should have prefix: %#v\n\tGot: %#v", tt.errorMsg, err.Error()))
	}
}

func TestEncodeIntTy(t *testing.T) {
	is := assert.New(t)

	tests := []encodeTestCase{
		{float64(1), "int8", "0000000000000000000000000000000000000000000000000000000000000001", ""},
		{float64(2), "int8", "0000000000000000000000000000000000000000000000000000000000000002", ""},
		{float64(-1), "int8", "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", ""},
		{float64(1), "int16", "0000000000000000000000000000000000000000000000000000000000000001", ""},
		{float64(1), "int24", "0000000000000000000000000000000000000000000000000000000000000001", ""},
		{float64(1), "int32", "0000000000000000000000000000000000000000000000000000000000000001", ""},
		{float64(1), "int64", "0000000000000000000000000000000000000000000000000000000000000001", ""},
		{float64(1), "int128", "0000000000000000000000000000000000000000000000000000000000000001", ""},
		{float64(1), "int256", "0000000000000000000000000000000000000000000000000000000000000001", ""},

		{int64(1), "int8",
			"0000000000000000000000000000000000000000000000000000000000000001", ""},

		{float64(2.1), "int8", "", "Expected int8 got"},
		{nil, "int8", "", "Expected int8 got"},
	}

	for _, tt := range tests {
		runEncodeTestCase(is, &tt)
	}
}

func TestEncodeUintTy(t *testing.T) {
	is := assert.New(t)

	tests := []encodeTestCase{
		{float64(1), "uint8", "0000000000000000000000000000000000000000000000000000000000000001", ""},
		{float64(1), "uint16", "0000000000000000000000000000000000000000000000000000000000000001", ""},
		{float64(1), "uint24", "0000000000000000000000000000000000000000000000000000000000000001", ""},
		{float64(1), "uint32", "0000000000000000000000000000000000000000000000000000000000000001", ""},
		{float64(1), "uint64", "0000000000000000000000000000000000000000000000000000000000000001", ""},
		{float64(1), "uint128", "0000000000000000000000000000000000000000000000000000000000000001", ""},
		{float64(1), "uint256", "0000000000000000000000000000000000000000000000000000000000000001", ""},

		{int64(1), "uint8",
			"0000000000000000000000000000000000000000000000000000000000000001", ""},

		{float64(-1), "uint8", "", "Expected uint8 got"},
	}

	for _, tt := range tests {
		runEncodeTestCase(is, &tt)
	}
}

func TestEncodeBool(t *testing.T) {
	is := assert.New(t)

	tests := []encodeTestCase{
		{true, "bool", "0000000000000000000000000000000000000000000000000000000000000001", ""},
		{false, "bool", "0000000000000000000000000000000000000000000000000000000000000000", ""},
		{1, "bool", "", "Expected bool got"},
	}

	for _, tt := range tests {
		runEncodeTestCase(is, &tt)
	}
}

func TestEncodeSlices(t *testing.T) {
	is := assert.New(t)

	tests := []encodeTestCase{
		{nil,
			"uint8[]",
			"" +
				"0000000000000000000000000000000000000000000000000000000000000000",
			""},
		{[]interface{}{},
			"uint8[]",
			"" +
				"0000000000000000000000000000000000000000000000000000000000000000",
			""},
		{[]interface{}{float64(1), float64(2)},
			"uint8[]",
			"" +
				"0000000000000000000000000000000000000000000000000000000000000002" +
				"0000000000000000000000000000000000000000000000000000000000000001" +
				"0000000000000000000000000000000000000000000000000000000000000002",
			""},

		{[]interface{}{float64(1), float64(-1)},
			"int32[]",
			"" +
				"0000000000000000000000000000000000000000000000000000000000000002" +
				"0000000000000000000000000000000000000000000000000000000000000001" +
				"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
			""},

		{[]interface{}{float64(1), float64(2.1)},
			"uint8[]",
			"",
			"uint8[] at 1: Expected uint8 got"},

		{[]interface{}{float64(1), float64(-1)},
			"uint8[]",
			"",
			"uint8[] at 1: Expected uint8 got"},
	}

	for _, tt := range tests {
		runEncodeTestCase(is, &tt)
	}
}

func TestEncodeArrays(t *testing.T) {
	is := assert.New(t)

	tests := []encodeTestCase{
		{[]interface{}{float64(1), float64(2), float64(3)},
			"uint8[3]",
			"" +
				"0000000000000000000000000000000000000000000000000000000000000001" +
				"0000000000000000000000000000000000000000000000000000000000000002" +
				"0000000000000000000000000000000000000000000000000000000000000003",
			""},
	}

	for _, tt := range tests {
		runEncodeTestCase(is, &tt)
	}
}

func TestEncodeString(t *testing.T) {
	is := assert.New(t)

	tests := []encodeTestCase{
		{"abcd",
			"string",
			"" +
				"0000000000000000000000000000000000000000000000000000000000000004" +
				"6162636400000000000000000000000000000000000000000000000000000000",
			""},
		{"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaabbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
			"string",
			"" +
				"0000000000000000000000000000000000000000000000000000000000000040" +
				"6161616161616161616161616161616161616161616161616161616161616161" +
				"6262626262626262626262626262626262626262626262626262626262626262",
			""},
		{1,
			"string",
			"",
			"Expected string got"},
	}

	for _, tt := range tests {
		runEncodeTestCase(is, &tt)
	}
}

func TestEncodeBytes(t *testing.T) {
	is := assert.New(t)

	tests := []encodeTestCase{
		{"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaabbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
			"bytes",
			"" +
				"0000000000000000000000000000000000000000000000000000000000000040" +
				"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" +
				"bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
			""},
		{"0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaabbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
			"bytes",
			"" +
				"0000000000000000000000000000000000000000000000000000000000000040" +
				"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" +
				"bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
			""},
		{1,
			"bytes",
			"",
			"Expected bytes in hex got"},
	}

	for _, tt := range tests {
		runEncodeTestCase(is, &tt)
	}
}
func TestEncodeAddress(t *testing.T) {
	is := assert.New(t)

	tests := []encodeTestCase{
		{"0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			"address",
			"000000000000000000000000aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			""},
		{"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			"address",
			"000000000000000000000000aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			""},
		{"aa",
			"address",
			"",
			"Expected address to have 20 bytes got: 1"},
		{1,
			"address",
			"",
			"Expected address in hex got"},
	}

	for _, tt := range tests {
		runEncodeTestCase(is, &tt)
	}
}

func TestEncodeFixedBytes(t *testing.T) {
	is := assert.New(t)

	tests := []encodeTestCase{
		{"0xaa",
			"bytes4",
			"00000000000000000000000000000000000000000000000000000000000000aa",
			""},
		{"0xaabbccdd",
			"bytes4",
			"00000000000000000000000000000000000000000000000000000000aabbccdd",
			""},

		{"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			"bytes32",
			"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			""},

		{"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			"bytes4",
			"",
			"Expected bytes4 to have 4 bytes got: 32"},
		{1,
			"bytes4",
			"",
			"Expected bytes4 in hex got"},
	}

	for _, tt := range tests {
		runEncodeTestCase(is, &tt)
	}
}

// Tests that all allowed types get recognized by the type parser.
func TestTypeRegexp(t *testing.T) {
	tests := []struct {
		blob string
		kind Type
	}{
		{"bool", Type{Kind: reflect.Bool, T: BoolTy, stringKind: "bool"}},
		{"bool[]", Type{IsSlice: true, SliceSize: -1, Kind: reflect.Bool, T: BoolTy, Elem: &Type{Kind: reflect.Bool, T: BoolTy, stringKind: "bool"}, stringKind: "bool[]"}},
		{"bool[2]", Type{IsArray: true, SliceSize: 2, Kind: reflect.Bool, T: BoolTy, Elem: &Type{Kind: reflect.Bool, T: BoolTy, stringKind: "bool"}, stringKind: "bool[2]"}},
		{"int8", Type{Kind: reflect.Int8, Type: int8_t, Size: 8, T: IntTy, stringKind: "int8"}},
		{"int16", Type{Kind: reflect.Int16, Type: int16_t, Size: 16, T: IntTy, stringKind: "int16"}},
		{"int32", Type{Kind: reflect.Int32, Type: int32_t, Size: 32, T: IntTy, stringKind: "int32"}},
		{"int64", Type{Kind: reflect.Int64, Type: int64_t, Size: 64, T: IntTy, stringKind: "int64"}},
		{"int256", Type{Kind: reflect.Ptr, Type: big_t, Size: 256, T: IntTy, stringKind: "int256"}},
		{"int8[]", Type{IsSlice: true, SliceSize: -1, Kind: reflect.Int8, Type: int8_t, Size: 8, T: IntTy, Elem: &Type{Kind: reflect.Int8, Type: int8_t, Size: 8, T: IntTy, stringKind: "int8"}, stringKind: "int8[]"}},
		{"int8[2]", Type{IsArray: true, SliceSize: 2, Kind: reflect.Int8, Type: int8_t, Size: 8, T: IntTy, Elem: &Type{Kind: reflect.Int8, Type: int8_t, Size: 8, T: IntTy, stringKind: "int8"}, stringKind: "int8[2]"}},
		{"int16[]", Type{IsSlice: true, SliceSize: -1, Kind: reflect.Int16, Type: int16_t, Size: 16, T: IntTy, Elem: &Type{Kind: reflect.Int16, Type: int16_t, Size: 16, T: IntTy, stringKind: "int16"}, stringKind: "int16[]"}},
		{"int16[2]", Type{IsArray: true, SliceSize: 2, Kind: reflect.Int16, Type: int16_t, Size: 16, T: IntTy, Elem: &Type{Kind: reflect.Int16, Type: int16_t, Size: 16, T: IntTy, stringKind: "int16"}, stringKind: "int16[2]"}},
		{"int32[]", Type{IsSlice: true, SliceSize: -1, Kind: reflect.Int32, Type: int32_t, Size: 32, T: IntTy, Elem: &Type{Kind: reflect.Int32, Type: int32_t, Size: 32, T: IntTy, stringKind: "int32"}, stringKind: "int32[]"}},
		{"int32[2]", Type{IsArray: true, SliceSize: 2, Kind: reflect.Int32, Type: int32_t, Size: 32, T: IntTy, Elem: &Type{Kind: reflect.Int32, Type: int32_t, Size: 32, T: IntTy, stringKind: "int32"}, stringKind: "int32[2]"}},
		{"int64[]", Type{IsSlice: true, SliceSize: -1, Kind: reflect.Int64, Type: int64_t, Size: 64, T: IntTy, Elem: &Type{Kind: reflect.Int64, Type: int64_t, Size: 64, T: IntTy, stringKind: "int64"}, stringKind: "int64[]"}},
		{"int64[2]", Type{IsArray: true, SliceSize: 2, Kind: reflect.Int64, Type: int64_t, Size: 64, T: IntTy, Elem: &Type{Kind: reflect.Int64, Type: int64_t, Size: 64, T: IntTy, stringKind: "int64"}, stringKind: "int64[2]"}},
		{"int256[]", Type{IsSlice: true, SliceSize: -1, Kind: reflect.Ptr, Type: big_t, Size: 256, T: IntTy, Elem: &Type{Kind: reflect.Ptr, Type: big_t, Size: 256, T: IntTy, stringKind: "int256"}, stringKind: "int256[]"}},
		{"int256[2]", Type{IsArray: true, SliceSize: 2, Kind: reflect.Ptr, Type: big_t, Size: 256, T: IntTy, Elem: &Type{Kind: reflect.Ptr, Type: big_t, Size: 256, T: IntTy, stringKind: "int256"}, stringKind: "int256[2]"}},
		{"uint8", Type{Kind: reflect.Uint8, Type: uint8_t, Size: 8, T: UintTy, stringKind: "uint8"}},
		{"uint16", Type{Kind: reflect.Uint16, Type: uint16_t, Size: 16, T: UintTy, stringKind: "uint16"}},
		{"uint32", Type{Kind: reflect.Uint32, Type: uint32_t, Size: 32, T: UintTy, stringKind: "uint32"}},
		{"uint64", Type{Kind: reflect.Uint64, Type: uint64_t, Size: 64, T: UintTy, stringKind: "uint64"}},
		{"uint256", Type{Kind: reflect.Ptr, Type: big_t, Size: 256, T: UintTy, stringKind: "uint256"}},
		{"uint8[]", Type{IsSlice: true, SliceSize: -1, Kind: reflect.Uint8, Type: uint8_t, Size: 8, T: UintTy, Elem: &Type{Kind: reflect.Uint8, Type: uint8_t, Size: 8, T: UintTy, stringKind: "uint8"}, stringKind: "uint8[]"}},
		{"uint8[2]", Type{IsArray: true, SliceSize: 2, Kind: reflect.Uint8, Type: uint8_t, Size: 8, T: UintTy, Elem: &Type{Kind: reflect.Uint8, Type: uint8_t, Size: 8, T: UintTy, stringKind: "uint8"}, stringKind: "uint8[2]"}},
		{"uint16[]", Type{IsSlice: true, SliceSize: -1, Kind: reflect.Uint16, Type: uint16_t, Size: 16, T: UintTy, Elem: &Type{Kind: reflect.Uint16, Type: uint16_t, Size: 16, T: UintTy, stringKind: "uint16"}, stringKind: "uint16[]"}},
		{"uint16[2]", Type{IsArray: true, SliceSize: 2, Kind: reflect.Uint16, Type: uint16_t, Size: 16, T: UintTy, Elem: &Type{Kind: reflect.Uint16, Type: uint16_t, Size: 16, T: UintTy, stringKind: "uint16"}, stringKind: "uint16[2]"}},
		{"uint32[]", Type{IsSlice: true, SliceSize: -1, Kind: reflect.Uint32, Type: uint32_t, Size: 32, T: UintTy, Elem: &Type{Kind: reflect.Uint32, Type: uint32_t, Size: 32, T: UintTy, stringKind: "uint32"}, stringKind: "uint32[]"}},
		{"uint32[2]", Type{IsArray: true, SliceSize: 2, Kind: reflect.Uint32, Type: uint32_t, Size: 32, T: UintTy, Elem: &Type{Kind: reflect.Uint32, Type: uint32_t, Size: 32, T: UintTy, stringKind: "uint32"}, stringKind: "uint32[2]"}},
		{"uint64[]", Type{IsSlice: true, SliceSize: -1, Kind: reflect.Uint64, Type: uint64_t, Size: 64, T: UintTy, Elem: &Type{Kind: reflect.Uint64, Type: uint64_t, Size: 64, T: UintTy, stringKind: "uint64"}, stringKind: "uint64[]"}},
		{"uint64[2]", Type{IsArray: true, SliceSize: 2, Kind: reflect.Uint64, Type: uint64_t, Size: 64, T: UintTy, Elem: &Type{Kind: reflect.Uint64, Type: uint64_t, Size: 64, T: UintTy, stringKind: "uint64"}, stringKind: "uint64[2]"}},
		{"uint256[]", Type{IsSlice: true, SliceSize: -1, Kind: reflect.Ptr, Type: big_t, Size: 256, T: UintTy, Elem: &Type{Kind: reflect.Ptr, Type: big_t, Size: 256, T: UintTy, stringKind: "uint256"}, stringKind: "uint256[]"}},
		{"uint256[2]", Type{IsArray: true, SliceSize: 2, Kind: reflect.Ptr, Type: big_t, Size: 256, T: UintTy, Elem: &Type{Kind: reflect.Ptr, Type: big_t, Size: 256, T: UintTy, stringKind: "uint256"}, stringKind: "uint256[2]"}},
		{"bytes32", Type{IsArray: true, SliceSize: 32, Elem: &Type{Kind: reflect.Uint8, Type: uint8_t, Size: 8, T: UintTy, stringKind: "uint8"}, T: FixedBytesTy, stringKind: "bytes32"}},
		{"bytes[]", Type{IsSlice: true, SliceSize: -1, Elem: &Type{IsSlice: true, SliceSize: -1, Elem: &Type{Kind: reflect.Uint8, Type: uint8_t, Size: 8, T: UintTy, stringKind: "uint8"}, T: BytesTy, stringKind: "bytes"}, stringKind: "bytes[]"}},
		{"bytes[2]", Type{IsArray: true, SliceSize: 2, Elem: &Type{IsSlice: true, SliceSize: -1, Elem: &Type{Kind: reflect.Uint8, Type: uint8_t, Size: 8, T: UintTy, stringKind: "uint8"}, T: BytesTy, stringKind: "bytes"}, stringKind: "bytes[2]"}},
		{"bytes32[]", Type{IsSlice: true, SliceSize: -1, Elem: &Type{IsArray: true, SliceSize: 32, Elem: &Type{Kind: reflect.Uint8, Type: uint8_t, Size: 8, T: UintTy, stringKind: "uint8"}, T: FixedBytesTy, stringKind: "bytes32"}, stringKind: "bytes32[]"}},
		{"bytes32[2]", Type{IsArray: true, SliceSize: 2, Elem: &Type{IsArray: true, SliceSize: 32, Elem: &Type{Kind: reflect.Uint8, Type: uint8_t, Size: 8, T: UintTy, stringKind: "uint8"}, T: FixedBytesTy, stringKind: "bytes32"}, stringKind: "bytes32[2]"}},
		{"string", Type{Kind: reflect.String, Size: -1, T: StringTy, stringKind: "string"}},
		{"string[]", Type{IsSlice: true, SliceSize: -1, Kind: reflect.String, T: StringTy, Size: -1, Elem: &Type{Kind: reflect.String, T: StringTy, Size: -1, stringKind: "string"}, stringKind: "string[]"}},
		{"string[2]", Type{IsArray: true, SliceSize: 2, Kind: reflect.String, T: StringTy, Size: -1, Elem: &Type{Kind: reflect.String, T: StringTy, Size: -1, stringKind: "string"}, stringKind: "string[2]"}},
		{"address", Type{Kind: reflect.Array, Type: address_t, Size: 20, T: AddressTy, stringKind: "address"}},
		{"address[]", Type{IsSlice: true, SliceSize: -1, Kind: reflect.Array, Type: address_t, T: AddressTy, Size: 20, Elem: &Type{Kind: reflect.Array, Type: address_t, Size: 20, T: AddressTy, stringKind: "address"}, stringKind: "address[]"}},
		{"address[2]", Type{IsArray: true, SliceSize: 2, Kind: reflect.Array, Type: address_t, T: AddressTy, Size: 20, Elem: &Type{Kind: reflect.Array, Type: address_t, Size: 20, T: AddressTy, stringKind: "address"}, stringKind: "address[2]"}},

		// TODO when fixed types are implemented properly
		// {"fixed", Type{}},
		// {"fixed128x128", Type{}},
		// {"fixed[]", Type{}},
		// {"fixed[2]", Type{}},
		// {"fixed128x128[]", Type{}},
		// {"fixed128x128[2]", Type{}},
	}
	for i, tt := range tests {
		typ, err := NewType(tt.blob)
		if err != nil {
			t.Errorf("type %d: failed to parse type string: %v", i, err)
		}
		if !reflect.DeepEqual(typ, tt.kind) {
			t.Errorf("type %d: parsed type mismatch:\n  have %+v\n  want %+v", i, typeWithoutStringer(typ), typeWithoutStringer(tt.kind))
		}
	}
}

func TestTypeCheck(t *testing.T) {
	for i, test := range []struct {
		typ   string
		input interface{}
		err   string
	}{
		{"uint", big.NewInt(1), ""},
		{"int", big.NewInt(1), ""},
		{"uint30", big.NewInt(1), ""},
		{"uint30", uint8(1), "abi: cannot use uint8 as type ptr as argument"},
		{"uint16", uint16(1), ""},
		{"uint16", uint8(1), "abi: cannot use uint8 as type uint16 as argument"},
		{"uint16[]", []uint16{1, 2, 3}, ""},
		{"uint16[]", [3]uint16{1, 2, 3}, ""},
		{"uint16[]", []uint32{1, 2, 3}, "abi: cannot use []uint32 as type []uint16 as argument"},
		{"uint16[3]", [3]uint32{1, 2, 3}, "abi: cannot use [3]uint32 as type [3]uint16 as argument"},
		{"uint16[3]", [4]uint16{1, 2, 3}, "abi: cannot use [4]uint16 as type [3]uint16 as argument"},
		{"uint16[3]", []uint16{1, 2, 3}, ""},
		{"uint16[3]", []uint16{1, 2, 3, 4}, "abi: cannot use [4]uint16 as type [3]uint16 as argument"},
		{"address[]", []common.Address{{1}}, ""},
		{"address[1]", []common.Address{{1}}, ""},
		{"address[1]", [1]common.Address{{1}}, ""},
		{"address[2]", [1]common.Address{{1}}, "abi: cannot use [1]array as type [2]array as argument"},
		{"bytes32", [32]byte{}, ""},
		{"bytes32", [33]byte{}, "abi: cannot use [33]uint8 as type [32]uint8 as argument"},
		{"bytes32", common.Hash{1}, ""},
		{"bytes31", [31]byte{}, ""},
		{"bytes31", [32]byte{}, "abi: cannot use [32]uint8 as type [31]uint8 as argument"},
		{"bytes", []byte{0, 1}, ""},
		{"bytes", [2]byte{0, 1}, ""},
		{"bytes", common.Hash{1}, ""},
		{"string", "hello world", ""},
		{"bytes32[]", [][32]byte{{}}, ""},
		{"function", [24]byte{}, ""},
	} {
		typ, err := NewType(test.typ)
		if err != nil {
			t.Fatal("unexpected parse error:", err)
		}

		err = typeCheck(typ, reflect.ValueOf(test.input))
		if err != nil && len(test.err) == 0 {
			t.Errorf("%d failed. Expected no err but got: %v", i, err)
			continue
		}
		if err == nil && len(test.err) != 0 {
			t.Errorf("%d failed. Expected err: %v but got none", i, test.err)
			continue
		}

		if err != nil && len(test.err) != 0 && err.Error() != test.err {
			t.Errorf("%d failed. Expected err: '%v' got err: '%v'", i, test.err, err)
		}
	}
}
