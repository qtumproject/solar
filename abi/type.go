// Copyright 2015 The go-ethereum Authors
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
	"regexp"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

const (
	IntTy byte = iota
	UintTy
	BoolTy
	StringTy
	SliceTy
	AddressTy
	FixedBytesTy
	BytesTy
	HashTy
	FixedPointTy
	FunctionTy
)

// Type is the reflection of the supported argument type
type Type struct {
	IsSlice, IsArray bool
	SliceSize        int

	Elem *Type

	Kind reflect.Kind
	Type reflect.Type
	Size int
	T    byte // Our own type checking

	stringKind string // holds the unparsed string for deriving signatures
}

var (
	// fullTypeRegex parses the abi types
	//
	// Types can be in the format of:
	//
	// 	Input  = Type [ "[" [ Number ] "]" ] Name .
	// 	Type   = [ "u" ] "int" [ Number ] [ x ] [ Number ].
	//
	// Examples:
	//
	//      string     int       uint       fixed
	//      string32   int8      uint8      uint[]
	//      address    int256    uint256    fixed128x128[2]
	fullTypeRegex = regexp.MustCompile(`([a-zA-Z0-9]+)(\[([0-9]*)\])?`)
	// typeRegex parses the abi sub types
	typeRegex = regexp.MustCompile("([a-zA-Z]+)(([0-9]+)(x([0-9]+))?)?")
)

// NewType creates a new reflection type of abi type given in t.
func NewType(t string) (typ Type, err error) {
	res := fullTypeRegex.FindAllStringSubmatch(t, -1)[0]
	// check if type is slice and parse type.
	switch {
	case res[3] != "":
		// err is ignored. Already checked for number through the regexp
		typ.SliceSize, _ = strconv.Atoi(res[3])
		typ.IsArray = true
	case res[2] != "":
		typ.IsSlice, typ.SliceSize = true, -1
	case res[0] == "":
		return Type{}, fmt.Errorf("abi: type parse error: %s", t)
	}
	if typ.IsArray || typ.IsSlice {
		sliceType, err := NewType(res[1])
		if err != nil {
			return Type{}, err
		}
		typ.Elem = &sliceType
		typ.stringKind = sliceType.stringKind + t[len(res[1]):]
		// Although we know that this is an array, we cannot return
		// as we don't know the type of the element, however, if it
		// is still an array, then don't determine the type.
		if typ.Elem.IsArray || typ.Elem.IsSlice {
			return typ, nil
		}
	}

	// parse the type and size of the abi-type.
	parsedType := typeRegex.FindAllStringSubmatch(res[1], -1)[0]
	// varSize is the size of the variable
	var varSize int
	if len(parsedType[3]) > 0 {
		var err error
		varSize, err = strconv.Atoi(parsedType[2])
		if err != nil {
			return Type{}, fmt.Errorf("abi: error parsing variable size: %v", err)
		}
	}
	// varType is the parsed abi type
	varType := parsedType[1]
	// substitute canonical integer
	if varSize == 0 && (varType == "int" || varType == "uint") {
		varSize = 256
		t += "256"
	}

	// only set stringKind if not array or slice, as for those,
	// the correct string type has been set
	if !(typ.IsArray || typ.IsSlice) {
		typ.stringKind = t
	}

	switch varType {
	case "int":
		typ.Kind, typ.Type = reflectIntKindAndType(false, varSize)
		typ.Size = varSize
		typ.T = IntTy
	case "uint":
		typ.Kind, typ.Type = reflectIntKindAndType(true, varSize)
		typ.Size = varSize
		typ.T = UintTy
	case "bool":
		typ.Kind = reflect.Bool
		typ.T = BoolTy
	case "address":
		typ.Kind = reflect.Array
		typ.Type = address_t
		typ.Size = 20
		typ.T = AddressTy
	case "string":
		typ.Kind = reflect.String
		typ.Size = -1
		typ.T = StringTy
	case "bytes":
		sliceType, _ := NewType("uint8")
		typ.Elem = &sliceType
		if varSize == 0 {
			typ.IsSlice = true
			typ.T = BytesTy
			typ.SliceSize = -1
		} else {
			typ.IsArray = true
			typ.T = FixedBytesTy
			typ.SliceSize = varSize
		}
	case "function":
		sliceType, _ := NewType("uint8")
		typ.Elem = &sliceType
		typ.IsArray = true
		typ.T = FunctionTy
		typ.SliceSize = 24
	default:
		return Type{}, fmt.Errorf("unsupported arg type: %s", t)
	}

	return
}

// String implements Stringer
func (t Type) String() (out string) {
	return t.stringKind
}

func (t Type) Pack(v interface{}) ([]byte, error) {
	// pretty.Println("type", t)

	if (t.IsSlice || t.IsArray) && t.T != BytesTy && t.T != FixedBytesTy && t.T != FunctionTy {
		return t.encodeSlice(v)
	}

	if v == nil {
		return nil, errors.Errorf("Expected %s got nil", t.String())
	}

	switch t.T {
	case IntTy:
		return t.encodeIntTy(v)
	case UintTy:
		return t.encodeUintTy(v)
	case StringTy:
		return t.encodeString(v)
	case BytesTy:
		return t.encodeBytes(v)
	case AddressTy:
		return t.encodeAddress(v)
	case FixedBytesTy:
		return t.encodeFixedBytes(v)
	case BoolTy:
		return t.encodeBool(v)
	}

	return nil, nil
}

func (t Type) hexStringToBytes(v string) ([]byte, error) {
	if strings.HasPrefix(v, "0x") {
		v = v[2:]
	}

	bytes, err := hex.DecodeString(v)
	if err != nil {
		return nil, errors.Errorf("Expected %s in hex got: %v", t.String(), v)
	}

	return bytes, nil
}

func (t Type) encodeBool(v interface{}) ([]byte, error) {
	switch v := v.(type) {
	case bool:
		if v {
			return t.encodeUintTy(int64(1))
		}

		return t.encodeUintTy(int64(0))
	default:
		return nil, errors.Errorf("Expected %s got: %v", t.String(), v)
	}
}

func (t Type) encodeString(v interface{}) ([]byte, error) {
	switch v := v.(type) {
	case string:
		return packBytesSlice([]byte(v), len(v)), nil
	default:
		return nil, errors.Errorf("Expected %s got: %v", t.String(), v)
	}
}

func (t Type) encodeBytes(v interface{}) ([]byte, error) {
	switch v := v.(type) {
	case string:
		bytes, err := t.hexStringToBytes(v)
		if err != nil {
			return nil, err
		}
		return packBytesSlice(bytes, len(bytes)), nil
	case []byte:
		return packBytesSlice(v, len(v)), nil
	default:
		return nil, errors.Errorf("Expected %s in hex got: %v", t.String(), v)
	}
}

func (t Type) encodeAddress(v interface{}) ([]byte, error) {
	switch v := v.(type) {
	case string:
		bytes, err := t.hexStringToBytes(v)
		if err != nil {
			return nil, err
		}

		if len(bytes) != 20 {
			return nil, errors.Errorf("Expected %s to have 20 bytes got: %v", t.String(), len(bytes))
		}

		return common.LeftPadBytes(bytes, 32), nil
	case []byte:
		bytes := v

		if len(bytes) == 20 {
			return nil, errors.Errorf("Expected %s to have 20 bytes got: %v", t.String(), len(bytes))
		}

		return common.LeftPadBytes(bytes, 32), nil
	default:
		return nil, errors.Errorf("Expected %s in hex got: %v", t.String(), v)
	}
}

func (t Type) encodeFixedBytes(v interface{}) ([]byte, error) {
	switch v := v.(type) {
	case string:
		bytes, err := t.hexStringToBytes(v)
		if err != nil {
			return nil, err
		}

		if len(bytes) > t.SliceSize {
			return nil, errors.Errorf("Expected %s to have %d bytes got: %d", t.String(), t.SliceSize, len(bytes))
		}

		return common.LeftPadBytes(bytes, 32), nil
	case []byte:
		bytes := v
		if len(bytes) > t.SliceSize {
			return nil, errors.Errorf("Expected %s to have %d bytes got: %d", t.String(), t.SliceSize, len(bytes))
		}
		return common.LeftPadBytes(bytes, 32), nil
	default:
		return nil, errors.Errorf("Expected %s in hex got: %v", t.String(), v)
	}
}

func (t Type) encodeSlice(v interface{}) ([]byte, error) {
	var vv []interface{}
	if v != nil {
		var ok bool
		vv, ok = v.([]interface{})
		if !ok {
			return nil, errors.Errorf("Expected %s got: %v", t.String(), v)
		}
	}

	var packed []byte
	for i := 0; i < len(vv); i++ {
		data, err := t.Elem.Pack(vv[i])
		if err != nil {
			return nil, errors.Wrapf(err, "%s at %d", t.String(), i)
		}
		packed = append(packed, data...)
	}

	if t.IsSlice {
		return packBytesSlice(packed, len(vv)), nil
	}

	return packed, nil
}

// TODO: handle truncation
func (t Type) encodeIntTy(v interface{}) ([]byte, error) {
	fmt.Println("encode int", v)
	switch v := v.(type) {
	case int, int8, int16, int32, int64:
		n, err := strconv.Atoi(fmt.Sprintf("%d", v))
		if err != nil {
			return nil, errors.Errorf("Expected %s got: %v", t.String(), v)
		}
		i := big.NewInt(int64(n))
		return U256(i), nil
	case *big.Int:
		return U256(v), nil
	case float64:
		f := big.NewFloat(v)
		if !f.IsInt() {
			return nil, errors.Errorf("Expected %s got: %v", t.String(), v)
		}

		i, _ := f.Int(nil)

		return U256(i), nil
	case string:
		var i big.Int
		_, ok := i.SetString(v, 0)
		if !ok {
			return nil, errors.Errorf("Expected big number string for %s got: %v", t.String(), v)
		}

		return U256(&i), nil
	default:
		return nil, errors.Errorf("Expected %s got: %v", t.String(), v)
	}
}

func (t Type) encodeUintTy(v interface{}) ([]byte, error) {
	switch v := v.(type) {
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64:
		n, err := strconv.Atoi(fmt.Sprintf("%d", v))
		if err != nil {
			return nil, errors.Errorf("Expected %s got: %v", t.String(), v)
		}

		if n < 0 {
			return nil, errors.Errorf("Expected %s got: %v", t.String(), v)
		}

		i := big.NewInt(int64(n))
		return U256(i), nil
	case *big.Int:
		if v.Sign() == -1 {
			return nil, errors.Errorf("Expected %s got: %v", t.String(), v)
		}
		return U256(v), nil
	case float64:
		f := big.NewFloat(v)
		if !f.IsInt() || f.Sign() == -1 {
			return nil, errors.Errorf("Expected %s got: %v", t.String(), f.String())
		}

		i, _ := f.Int(nil)

		return U256(i), nil
	case string:
		var i big.Int
		_, ok := i.SetString(v, 0)
		if !ok {
			return nil, errors.Errorf("Expected big number string for %s got: %v", t.String(), v)
		}

		return U256(&i), nil
	default:
		return nil, errors.Errorf("Expected %s got: %v", t.String(), v)
	}
}

func (t Type) pack(v reflect.Value) ([]byte, error) {
	// dereference pointer first if it's a pointer
	v = indirect(v)

	if err := typeCheck(t, v); err != nil {
		return nil, err
	}

	if (t.IsSlice || t.IsArray) && t.T != BytesTy && t.T != FixedBytesTy && t.T != FunctionTy {
		var packed []byte

		for i := 0; i < v.Len(); i++ {
			val, err := t.Elem.pack(v.Index(i))
			if err != nil {
				return nil, err
			}
			packed = append(packed, val...)
		}
		if t.IsSlice {
			return packBytesSlice(packed, v.Len()), nil
		} else if t.IsArray {
			return packed, nil
		}
	}

	return packElement(t, v), nil
}

// requireLengthPrefix returns whether the type requires any sort of length
// prefixing.
func (t Type) requiresLengthPrefix() bool {
	return t.T != FixedBytesTy && (t.T == StringTy || t.T == BytesTy || t.IsSlice)
}
