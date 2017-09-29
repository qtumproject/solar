package jsonabi

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strings"
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

func testMethodABI(methodName, typ string) *abi.ABI {
	const abiTemplateString = `
	[{
		"name": "{{.MethodName}}",
		"type": "function",
		"payable": false,
		"inputs": [
			{
				"name": "_a",
				"type": "{{.Type}}"
			}
		],
		"outputs": [],
		"constant": false
	}]
	`

	abiTemplate := template.Must(template.New("testABI").Parse(abiTemplateString))

	type ctx struct {
		MethodName string
		Type       string
	}

	var abiJSON bytes.Buffer
	err := abiTemplate.Execute(&abiJSON, &ctx{
		MethodName: methodName,
		Type:       typ,
	})

	if err != nil {
		panic(err)
	}

	testABI, err := abi.JSON(&abiJSON)
	if err != nil {
		panic(err)
	}

	return &testABI
}

func TestEncodeInty(t *testing.T) {
	is := assert.New(t)

	examples := []struct {
		argType   string
		val       string
		errString string
	}{
		{"int8", `[1]`, ""},
		{"int16", `[1]`, ""},
		{"int24", `[1]`, ""},
		{"int32", `[1]`, ""},
		{"int40", `[1]`, ""},
		{"int64", `[1]`, ""},
		{"int72", `[1]`, ""},
		{"int128", `[1]`, ""},
		{"int256", `[1]`, ""},

		{"int256", `[1.0]`, ""},

		{"int32", `["0xffff"]`, ""},
		{"int64", `["0xffff"]`, ""},
		{"int72", `["0xffff"]`, ""},
		{"int160", `["0xffff"]`, ""},
		{"int256", `["0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"]`, ""},

		{"int256", `[1.21]`, "Expected integer got float"},
		{"int32", `[1.21]`, "Expected integer got float"},
		{"int32", `[[1,2]]`, "Expected integer got"},
		{"int32", `[null]`, "Expected integer got"},
	}

	for _, e := range examples {
		methodName := "testMethod"
		testABI := testMethodABI(methodName, e.argType)

		enc := Encoder{
			abi: testABI,
		}

		_, err := enc.EncodeJSONValues(methodName, []byte(e.val))

		if e.errString != "" {
			is.True(strings.HasPrefix(err.Error(), e.errString),
				fmt.Sprintf("Error message should have prefix: %#v\n\tGot: %#v", e.errString, err.Error()))
		} else {
			is.NoError(err)
		}
	}
}

func TestEncodeStringy(t *testing.T) {
	is := assert.New(t)

	examples := []struct {
		argType   string
		val       string
		errString string
	}{
		{"string", `["abcd"]`, ""},
		{"string", `[1]`, "Expected string got"},
	}

	for _, e := range examples {
		methodName := "testMethod"
		testABI := testMethodABI(methodName, e.argType)

		enc := Encoder{
			abi: testABI,
		}

		_, err := enc.EncodeJSONValues(methodName, []byte(e.val))

		if e.errString != "" {
			is.True(strings.HasPrefix(err.Error(), e.errString),
				fmt.Sprintf("Error message should have prefix: %#v\n\tGot: %#v", e.errString, err.Error()))
		} else {
			is.NoError(err)
		}
	}
}

func TestEncodeBytes(t *testing.T) {
	is := assert.New(t)

	examples := []struct {
		argType   string
		val       string
		errString string
	}{
		{"bytes", `["ffaa00ee"]`, ""},
		{"bytes", `["0xffaa00ee"]`, ""},
		{"bytes", `["0xInvalidHexString"]`, "Expected hex string"},
		{"bytes", `[1]`, "Expected hex string"},
	}

	for _, e := range examples {
		methodName := "testMethod"
		testABI := testMethodABI(methodName, e.argType)

		enc := Encoder{
			abi: testABI,
		}

		data, err := enc.EncodeJSONValues(methodName, []byte(e.val))
		fmt.Println("data", hex.EncodeToString(data))

		if e.errString != "" {
			is.True(strings.HasPrefix(err.Error(), e.errString),
				fmt.Sprintf("Error message should have prefix: %#v\n\tGot: %#v", e.errString, err.Error()))
		} else {
			is.NoError(err)
		}
	}
}
