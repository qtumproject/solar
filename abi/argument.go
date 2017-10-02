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
	"encoding/json"
	"fmt"
	"reflect"
)

// Argument holds the name of the argument and the corresponding type.
// Types are used when packing and testing arguments.
type Argument struct {
	Name    string
	Type    Type
	Indexed bool // indexed is only used by events
}

func (a *Argument) UnmarshalJSON(data []byte) error {
	var extarg struct {
		Name    string
		Type    string
		Indexed bool
	}
	err := json.Unmarshal(data, &extarg)
	if err != nil {
		return fmt.Errorf("argument json err: %v", err)
	}

	a.Type, err = NewType(extarg.Type)
	if err != nil {
		return err
	}
	a.Name = extarg.Name
	a.Indexed = extarg.Indexed

	return nil
}

type Arguments []Argument

func (types Arguments) Pack(vals []interface{}) ([]byte, error) {
	// Make sure arguments match up and pack them
	if len(types) != len(vals) {
		return nil, fmt.Errorf("argument count mismatch: %d for %d", len(vals), len(types))
	}
	// variable input is the output appended at the end of packed
	// output. This is used for strings and bytes types input.
	var variableInput []byte

	var ret []byte
	for i, val := range vals {
		input := types[i]
		// pack the input
		packed, err := input.Type.Pack(val)
		if err != nil {
			return nil, fmt.Errorf("argv[%d] '%s': %v", i, input.Name, err)
		}

		// check for a slice type (string, bytes, slice)
		if input.Type.requiresLengthPrefix() {
			// calculate the offset
			offset := len(types)*32 + len(variableInput)
			// set the offset
			ret = append(ret, packNum(reflect.ValueOf(offset))...)
			// Append the packed output to the variable input. The variable input
			// will be appended at the end of the input.
			variableInput = append(variableInput, packed...)
		} else {
			// append the packed value to the input
			ret = append(ret, packed...)
		}
	}
	// append the variable input at the end of the packed input
	ret = append(ret, variableInput...)

	return ret, nil
}
