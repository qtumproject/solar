package abi

import (
	"encoding/hex"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type encodeValuesTestCase struct {
	values    []interface{}
	typeNames []string
	result    string
	errorMsg  string
}

func runEncodeValuesTestCase(is *assert.Assertions, tt *encodeValuesTestCase) {
	var args Arguments
	for i, name := range tt.typeNames {
		ty, err := NewType(name)
		is.NoError(err)
		arg := Argument{
			Name: fmt.Sprintf("arg%d", i),
			Type: ty,
		}

		args = append(args, arg)
	}

	data, err := args.Pack(tt.values)

	if tt.result != "" {
		is.NoError(err, "packing")
		is.Equal(tt.result, hex.EncodeToString(data))
	} else if tt.errorMsg != "" {
		is.Error(err)
		is.True(strings.HasPrefix(err.Error(), tt.errorMsg),
			fmt.Sprintf("Error message should have prefix: %#v\n\tGot: %#v", tt.errorMsg, err.Error()))
	}
}

func TestEncodeValues(t *testing.T) {
	is := assert.New(t)

	tests := []encodeValuesTestCase{
		{
			[]interface{}{1, 2, -3},
			[]string{"int8", "uint16", "int32"},
			"" +
				"0000000000000000000000000000000000000000000000000000000000000001" +
				"0000000000000000000000000000000000000000000000000000000000000002" +
				"fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffd",
			"",
		},

		{
			[]interface{}{1, 2, -3},
			[]string{"int8", "int16", "uint32"},
			"",
			"[2]: Expected uint32 got: -3",
		},

		{
			[]interface{}{"abc", 2, "0xffaaffaa"},
			[]string{"string", "uint16", "bytes"},
			"" +
				"0000000000000000000000000000000000000000000000000000000000000060" +
				"0000000000000000000000000000000000000000000000000000000000000002" +
				"00000000000000000000000000000000000000000000000000000000000000a0" +
				// arg0: string
				"0000000000000000000000000000000000000000000000000000000000000003" +
				"6162630000000000000000000000000000000000000000000000000000000000" +
				// arg2: bytes
				"0000000000000000000000000000000000000000000000000000000000000004" +
				"ffaaffaa00000000000000000000000000000000000000000000000000000000",
			"",
		},
	}

	for _, tt := range tests {
		runEncodeValuesTestCase(is, &tt)
	}
}
