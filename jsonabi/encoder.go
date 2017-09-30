package jsonabi

import (
	"encoding/hex"
	"encoding/json"
	"math/big"
	"strings"

	"github.com/pkg/errors"

	"github.com/hayeah/solar/abi"
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

	args := method.Inputs

	if len(args) != len(vals) {
		return nil, errors.Errorf("Expected %d arguments, got %d", len(method.Inputs), len(vals))
	}

	vals2, err := massageJSONValuesToABIValues(args, vals)
	if err != nil {
		return nil, err
	}

	return e.abi.Pack(methodName, vals2...)
}

func EncodeJSONValues(args abi.Arguments, jsonValues []byte) ([]byte, error) {
	var vals []interface{}
	err := json.Unmarshal(jsonValues, &vals)
	if err != nil {
		return nil, err
	}

	return EncodeValues(args, vals...)
}

func EncodeValues(args abi.Arguments, vals ...interface{}) ([]byte, error) {
	vals2, err := massageJSONValuesToABIValues(args, vals)
	if err != nil {
		return nil, err
	}
	return args.Pack(vals2)
}

func massageJSONValuesToABIValues(args abi.Arguments, vals []interface{}) ([]interface{}, error) {
	var vals2 []interface{}
	// massage the input arguments
	for i, arg := range args {
		val := vals[i]
		t := arg.Type

		if t.IsSlice {
			t.Elem
		}

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
		case abi.FixedBytesTy:
			hexstr, ok := val.(string)
			if !ok {
				return nil, errors.Errorf("Expected hex string: %#v", val)
			}

			bytes, err := hexstringToBytes(hexstr)
			if err != nil {
				return nil, err
			}

			// pretty.Println("fixed bytes type", t)

			fixedBytes, err := bytesToFixedBytesTy(bytes, t.SliceSize)
			if err != nil {
				return nil, err
			}
			vals2 = append(vals2, fixedBytes)
		case abi.BytesTy:
			switch val := val.(type) {
			case string:
				bytes, err := hexstringToBytes(val)
				if err != nil {
					return nil, err
				}
				vals2 = append(vals2, bytes)
			default:
				return nil, errors.Errorf("Expected hex string: %#v", val)
			}
		default:
			vals2 = append(vals2, val)
		}
	}

	return vals2, nil
}

func hexstringToBytes(s string) ([]byte, error) {
	if strings.HasPrefix(s, "0x") {
		s = s[2:]
	}

	bytes, err := hex.DecodeString(s)
	if err != nil {
		return nil, errors.Errorf("Expected hex string: %#v", s)
	}

	return bytes, nil
}

func bytesToFixedBytesTy(bytes []byte, size int) (interface{}, error) {
	// FIXME: Better way to convert slice to fixed byte array
	if len(bytes) > size {
		return nil, errors.Errorf("Expected %d bytes, got: %d", size, len(bytes))
	}

	nbytes := size
	if size > len(bytes) {
		nbytes = len(bytes)
	}

	switch size {
	case 32:
		var buf [32]byte
		copy(buf[:], bytes[0:nbytes])
		return buf, nil
	case 16:
		var buf [16]byte
		copy(buf[:], bytes[0:nbytes])
		return buf, nil
	case 8:
		var buf [8]byte
		copy(buf[:], bytes[0:nbytes])
		return buf, nil
	case 4:
		var buf [4]byte
		copy(buf[:], bytes[0:nbytes])
		return buf, nil
	default:
		return nil, errors.Errorf("Unsupported fixed bytes size %d", size)
	}
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
	// if size == 8 {
	// 	i := i.Int64()
	// 	return int8(i)
	// }

	// if size == 16 {
	// 	i := i.Int64()
	// 	return int16(i)
	// }

	// if size == 32 {
	// 	i := i.Int64()
	// 	return int32(i)
	// }

	// if size == 64 {
	// 	i := i.Int64()
	// 	return i
	// }

	// For any other sizes use big.Int
	return &i
}

func bigfloatToInty(f *big.Float, size int) (interface{}, error) {
	// TODO check truncation
	if !f.IsInt() {
		return nil, errors.Errorf("Expected integer got float: %s", f.String())
	}

	// go-ethereum's ABI encoding is picky about the exact type of a integer. Make it happy.
	// if size == 8 {
	// 	i, _ := f.Int64()
	// 	return int8(i), nil
	// }

	// if size == 16 {
	// 	i, _ := f.Int64()
	// 	return int16(i), nil
	// }

	// if size == 32 {
	// 	i, _ := f.Int64()
	// 	return int32(i), nil
	// }

	// if size == 64 {
	// 	i, _ := f.Int64()
	// 	return i, nil
	// }

	// For any other sizes use big.Int
	bigI, _ := f.Int(nil)
	return bigI, nil
}
