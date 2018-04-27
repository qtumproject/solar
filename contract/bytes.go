package contract

import (
	"encoding/hex"
	"fmt"
)

var bytesFormatWithPrefix = false

func SetFormatBytesWithPrefix(v bool) {
	// possible race condition... but usually done in setup, so whatever...
	bytesFormatWithPrefix = v
}

type Bytes []byte

func (b Bytes) String() string {
	if bytesFormatWithPrefix {
		return "0x" + hex.EncodeToString(b)
	} else {
		return hex.EncodeToString(b)
	}
}

func (b Bytes) MarshalJSON() ([]byte, error) {
	hexstr := fmt.Sprintf(`"%s"`, b.String())
	return []byte(hexstr), nil
}

func (b *Bytes) UnmarshalJSON(data []byte) error {
	// strip the quotes \"\"
	data = data[1 : len(data)-1]

	// strip '0x' prefix
	if data[0] == '0' && (data[1] == 'x' || data[1] == 'X') {
		data = data[2:]
	}

	dst := make([]byte, hex.DecodedLen(len(data)))
	_, err := hex.Decode(dst, data)
	if err != nil {
		return err
	}
	*b = dst

	return nil
}
