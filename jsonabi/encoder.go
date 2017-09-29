package jsonabi

import (
	"encoding/hex"
	"encoding/json"
	"math/big"
	"strings"

	"github.com/pkg/errors"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

type Encoder struct {
	abi *abi.ABI
}

func (e *Encoder) EncodeJSONValues(methodName string, jsonValues []byte) ([]byte, error) {
	var vals []interface{}
	err := json.Unmarshal(jsonValues, &vals)
	if err != nil {
		return nil, err
	}

	return e.encodeValues(methodName, vals...)
}

// todo JSON arrays
// todo Map as format specifier
func (e *Encoder) encodeValues(methodName string, vals ...interface{}) ([]byte, error) {
	method, ok := e.abi.Methods[methodName]
	if !ok {
		return nil, errors.Errorf("Cannot find method: %s", methodName)
	}

	if len(method.Inputs) != len(vals) {
		return nil, errors.Errorf("Expected %d arguments, got %d", len(method.Inputs), len(vals))
	}

	var vals2 []interface{}
	// massage the input arguments
	for i, arg := range method.Inputs {
		val := vals[i]
		t := arg.Type

		switch t.T {
		case abi.IntTy:
			// string or float
			switch val := val.(type) {
			case float64:
				bf := big.NewFloat(val)

				i, err := bigfloatToInty(bf, t.Size)
				if err != nil {
					return nil, err
				}
				vals2 = append(vals2, i)
			case string:
				i, err := stringToInty(val, t.Size)
				if err != nil {
					return nil, err
				}
				vals2 = append(vals2, i)
			default:
				return nil, errors.Errorf("Expected integer got: %v", val)
			}
		case abi.StringTy:
			switch val := val.(type) {
			case string:
				vals2 = append(vals2, val)
			default:
				return nil, errors.Errorf("Expected string got: %#v", val)
			}
		case abi.BytesTy:
			switch val := val.(type) {
			case string:
				if strings.HasPrefix(val, "0x") {
					val = val[2:]
				}

				bytes, err := hex.DecodeString(val)
				if err != nil {
					return nil, errors.Errorf("Expected hex string: %#v", val)
				}
				vals2 = append(vals2, bytes)
			default:
				return nil, errors.Errorf("Expected hex string: %#v", val)
			}
		default:
			vals2 = append(vals2, val)
		}
	}

	return e.abi.Pack(methodName, vals2...)
}

func stringToInty(s string, size int) (interface{}, error) {
	var i big.Int

	_, ok := i.SetString(s, 0)
	if !ok {
		return nil, errors.Errorf("Cannot convert to integer: %s", s)
	}

	return bigIntToInty(&i, size), nil
}

// TODO check truncation
func bigIntToInty(i *big.Int, size int) interface{} {
	// go-ethereum's ABI encoding is picky about the exact type of a integer. Make it happy.
	if size == 8 {
		i := i.Int64()
		return int8(i)
	}

	if size == 16 {
		i := i.Int64()
		return int16(i)
	}

	if size == 32 {
		i := i.Int64()
		return int32(i)
	}

	if size == 64 {
		i := i.Int64()
		return i
	}

	// For any other sizes use big.Int
	return &i
}

func bigfloatToInty(f *big.Float, size int) (interface{}, error) {
	// TODO check truncation
	if !f.IsInt() {
		return nil, errors.Errorf("Expected integer got float: %s", f.String())
	}

	// go-ethereum's ABI encoding is picky about the exact type of a integer. Make it happy.
	if size == 8 {
		i, _ := f.Int64()
		return int8(i), nil
	}

	if size == 16 {
		i, _ := f.Int64()
		return int16(i), nil
	}

	if size == 32 {
		i, _ := f.Int64()
		return int32(i), nil
	}

	if size == 64 {
		i, _ := f.Int64()
		return i, nil
	}

	// For any other sizes use big.Int
	bigI, _ := f.Int(nil)
	return bigI, nil
}
