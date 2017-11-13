package solar

import (
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto/sha3"
	"github.com/hayeah/solar/abi"

	"encoding/hex"
)

type ABIDefinition struct {
	Name     string    `json:"name"`
	Type     string    `json:"type"`
	Payable  bool      `json:"payable"`
	Inputs   []ABIType `json:"inputs"`
	Outputs  []ABIType `json:"outputs"`
	Constant bool      `json:"constant"`
}

type ABIType struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type CompiledContract struct {
	Name string          `json:"name"`
	ABI  []ABIDefinition `json:"abi"`
	Bin  Bytes           `json:"bin"`
	// KECAAK256 of bytecode without auxdata
	BinKeccak256 Bytes `json:"binhash"`
}

func (c *CompiledContract) encodingABI() (*abi.ABI, error) {
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

type rawCompiledContract struct {
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

func (c *rawCompiledContract) BinHash256() []byte {
	bin := c.BinWithoutAuxData()
	h := sha3.NewKeccak256()
	h.Write(bin)
	binDigest := h.Sum(nil)
	return binDigest
}

func (c *rawCompiledContract) BinWithoutAuxData() []byte {
	// https://solidity.readthedocs.io/en/develop/miscellaneous.html#encoding-of-the-metadata-hash-in-the-bytecode
	// 0xa1 0x65 'b' 'z' 'z' 'r' '0' 0x58 0x20 <32 bytes swarm hash> 0x00 0x29
	// a1 65 62 7a 7a 72 30 58 20 [32 bytes] 0x00 0x29
	// 11 + 32 bytes

	return c.Bin[0 : len(c.Bin)-11-32-1]
}

func (c *rawCompiledContract) UnmarshalJSON(data []byte) error {
	type dataStruct struct {
		RawMetadata string `json:"metadata"`
		BinStr      string `json:"bin"`
	}

	var dest dataStruct

	err := json.Unmarshal(data, &dest)
	if err != nil {
		return err
	}

	// Recursively parse Metadata, which is a json string.
	err = json.Unmarshal([]byte(dest.RawMetadata), &c.Metadata)
	if err != nil {
		return err
	}

	bin, err := hex.DecodeString(dest.BinStr)
	if err != nil {
		return err
	}

	c.RawMetadata = dest.RawMetadata
	c.Bin = bin

	return nil
}

type Bytes []byte

func (b Bytes) String() string {
	return hex.EncodeToString(b)
}

func (b Bytes) MarshalJSON() ([]byte, error) {
	hexstr := fmt.Sprintf("\"%s\"", b.String())
	return []byte(hexstr), nil
}

func (b *Bytes) UnmarshalJSON(data []byte) error {
	// strip the quotes \"\"
	hexstr := data[1 : len(data)-1]
	dst := make([]byte, hex.DecodedLen(len(hexstr)))
	_, err := hex.Decode(dst, hexstr)
	if err != nil {
		return err
	}
	*b = dst

	return nil
}
