package jsonabi

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

func testMethodABI(typeName string) abi.Arguments {
	const abiTemplateString = `
	[{
		"name": "arg1",
		"type": "{{.Type}}"
	}]
	`

	abiTemplate := template.Must(template.New("testABI").Parse(abiTemplateString))

	type ctx struct {
		Type string
	}

	var abiJSON bytes.Buffer
	err := abiTemplate.Execute(&abiJSON, &ctx{
		Type: typeName,
	})

	if err != nil {
		panic(err)
	}

	var args abi.Arguments
	dec := json.NewDecoder(&abiJSON)

	err = dec.Decode(&args)
	if err != nil {
		panic(err)
	}

	return args
}

type encodingTestCase struct {
	argType  string
	input    string
	hasError bool
	result   string
}

func runEncodingOneTest(is *assert.Assertions, e encodingTestCase) {
	args := testMethodABI(e.argType)

	data, err := EncodeJSONValues(args, []byte(e.input))
	fmt.Println("data", hex.EncodeToString(data))

	if e.hasError {
		is.True(strings.HasPrefix(err.Error(), e.result),
			fmt.Sprintf("Error message should have prefix: %#v\n\tGot: %#v", e.result, err.Error()))
	} else {
		is.NoError(err)
		is.Equal(e.result, hex.EncodeToString(data))
	}
}

func TestEncodeInty(t *testing.T) {
	is := assert.New(t)

	examples := []encodingTestCase{
		{"int8", `[1]`, false, "0000000000000000000000000000000000000000000000000000000000000001"},
		{"int8", `[2]`, false, "0000000000000000000000000000000000000000000000000000000000000002"},
		{"int8", `[-1]`, false, "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"},
		{"int16", `[1]`, false, "0000000000000000000000000000000000000000000000000000000000000001"},
		{"int24", `[1]`, false, "0000000000000000000000000000000000000000000000000000000000000001"},
		{"int32", `[1]`, false, "0000000000000000000000000000000000000000000000000000000000000001"},
		{"int40", `[1]`, false, "0000000000000000000000000000000000000000000000000000000000000001"},
		{"int64", `[1]`, false, "0000000000000000000000000000000000000000000000000000000000000001"},
		{"int72", `[1]`, false, "0000000000000000000000000000000000000000000000000000000000000001"},
		{"int128", `[1]`, false, "0000000000000000000000000000000000000000000000000000000000000001"},
		{"int256", `[1]`, false, "0000000000000000000000000000000000000000000000000000000000000001"},
		{"int256", `[1.0]`, false, "0000000000000000000000000000000000000000000000000000000000000001"},

		{"int32", `["0xffff"]`, false, "000000000000000000000000000000000000000000000000000000000000ffff"},
		{"int64", `["0xffff"]`, false, "000000000000000000000000000000000000000000000000000000000000ffff"},
		{"int256",
			`["0xfafa00000000000000000000000000abcd00000000000000000000000000fafa"]`,
			false,
			"fafa00000000000000000000000000abcd00000000000000000000000000fafa"},

		{"int256", `[1.21]`, true, "Expected integer got float"},
		{"int32", `[1.21]`, true, "Expected integer got float"},
		{"int32", `[null]`, true, "Expected integer got"},
	}

	for _, e := range examples {
		runEncodingOneTest(is, e)
	}
}

func TestEncodeStringy(t *testing.T) {
	is := assert.New(t)

	examples := []encodingTestCase{
		{"string", `["abcdabcdabcdabcdabcdabcdabcdabcdabcdabcd"]`, false,
			"" +
				"0000000000000000000000000000000000000000000000000000000000000020" +
				"0000000000000000000000000000000000000000000000000000000000000028" +
				"6162636461626364616263646162636461626364616263646162636461626364" +
				"6162636461626364000000000000000000000000000000000000000000000000",
		},
		{"string", `[1]`, true, "Expected string got"},
	}

	for _, e := range examples {
		runEncodingOneTest(is, e)
	}
}

func TestEncodeBytes(t *testing.T) {
	is := assert.New(t)

	examples := []encodingTestCase{
		{"bytes", `["ffaa00ee"]`, false,
			"" +
				"0000000000000000000000000000000000000000000000000000000000000020" +
				"0000000000000000000000000000000000000000000000000000000000000004" +
				"ffaa00ee00000000000000000000000000000000000000000000000000000000",
		},
		{"bytes", `["0xffaa00ee"]`, false,
			"" +
				"0000000000000000000000000000000000000000000000000000000000000020" +
				"0000000000000000000000000000000000000000000000000000000000000004" +
				"ffaa00ee00000000000000000000000000000000000000000000000000000000",
		},

		{"bytes", `["0xInvalidHexString"]`, true,
			"Expected hex string"},
		{"bytes", `[1]`, true,
			"Expected hex string"},
	}

	for _, e := range examples {
		runEncodingOneTest(is, e)
	}
}
