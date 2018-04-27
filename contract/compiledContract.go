package contract

import (
	"encoding/hex"
	"encoding/json"

	"github.com/ethereum/go-ethereum/crypto/sha3"
	"github.com/pkg/errors"
	"github.com/qtumproject/solar/abi"
)

type ABIDefinition struct {
	Name     string    `json:"name"`
	Type     string    `json:"type"`
	Payable  bool      `json:"payable"`
	Inputs   []ABIType `json:"inputs"`
	Outputs  []ABIType `json:"outputs"`
	Constant bool      `json:"constant"`

	// Event
	Anonymous bool `json:"anonymous"`
}

type ABIType struct {
	Name string `json:"name"`
	Type string `json:"type"`

	// Event
	Indexed bool `json:"indexed"`
}

type CompiledContract struct {
	// Where the contract is defined
	Source string          `json:"source"`
	Name   string          `json:"name"`
	ABI    []ABIDefinition `json:"abi"`
	Bin    Bytes           `json:"bin"`
	// KECAAK256 of bytecode without auxdata
	BinKeccak256 Bytes `json:"binhash"`
}

func (c *CompiledContract) EncodingABI() (*abi.ABI, error) {
	jsonABI, err := json.Marshal(c.ABI)
	if err != nil {
		return nil, err
	}

	var encodingABI abi.ABI
	err = json.Unmarshal(jsonABI, &encodingABI)
	if err != nil {
		return nil, err
	}

	return &encodingABI, nil
}

func (c *CompiledContract) ToBytes(jsonParams []byte) (Bytes, error) {
	calldata := c.Bin

	abi, err := c.EncodingABI()
	if err != nil {
		return nil, errors.Wrap(err, "abi")
	}

	constructor := abi.Constructor

	if len(constructor.Inputs) == 0 && len(jsonParams) != 0 {
		return nil, errors.New("does not expect constructor params")
	}

	if len(constructor.Inputs) != 0 {
		var params []interface{}
		err = json.Unmarshal(jsonParams, &params)
		if err != nil {
			return nil, errors.Errorf("expected constructor params in JSON, got: %#v", string(jsonParams))
		}

		packedParams, err := abi.Constructor.Pack(params...)
		if err != nil {
			return nil, errors.Wrap(err, "constructor")
		}

		calldata = append(calldata, packedParams...)
	}

	return calldata, nil
}

type RawCompiledContract struct {
	RawMetadata string `json:"metadata"`
	Bin         []byte
	Metadata    struct {
		Output struct {
			Version  int64
			Language string
			ABI      []ABIDefinition `json:"abi"`
		}
	}
}

func (c *RawCompiledContract) BinHash256() []byte {
	bin := c.BinWithoutAuxData()
	h := sha3.NewKeccak256()
	h.Write(bin)
	binDigest := h.Sum(nil)
	return binDigest
}

func (c *RawCompiledContract) BinWithoutAuxData() []byte {
	// https://solidity.readthedocs.io/en/develop/miscellaneous.html#encoding-of-the-metadata-hash-in-the-bytecode
	// 0xa1 0x65 'b' 'z' 'z' 'r' '0' 0x58 0x20 <32 bytes swarm hash> 0x00 0x29
	// a1 65 62 7a 7a 72 30 58 20 [32 bytes] 0x00 0x29
	// 11 + 32 bytes

	return c.Bin[0 : len(c.Bin)-11-32-1]
}

func (c *RawCompiledContract) UnmarshalJSON(data []byte) error {
	type dataStruct struct {
		RawMetadata string `json:"metadata"`
		BinStr      string `json:"bin"`
	}

	var dest dataStruct

	err := json.Unmarshal(data, &dest)
	if err != nil {
		return errors.Wrap(err, "parse contract raw metadata")
	}

	if dest.RawMetadata == "" && dest.BinStr == "" {
		// this contract is interface only
		return nil
	}

	// Recursively parse Metadata, which is a json string.
	err = json.Unmarshal([]byte(dest.RawMetadata), &c.Metadata)
	if err != nil {
		return errors.Wrap(err, "parse contract metadata")
	}

	bin, err := hex.DecodeString(dest.BinStr)
	if err != nil {
		return errors.Wrap(err, "decode byte string")
	}

	c.RawMetadata = dest.RawMetadata
	c.Bin = bin

	return nil
}
