package b58addr

import base58 "github.com/jbenet/go-base58"
import "encoding/hex"

/*
An arbitrarily sized payload.

* A set of 58 alphanumeric symbols consisting of easily distinguished uppercase and lowercase letters (0OIl are not used)
* One byte of version/application information. Bitcoin addresses use 0x00 for this byte (future ones may use 0x05).
* Four bytes (32 bits) of SHA256-based error checking code. This code can be used to automatically detect and possibly correct typographical errors.
* An extra step for preservation of leading zeroes in the data.

data := version + payload
checksum := take4(sha256(sha256(data)))
addr := data + checksum
*/

// qcli gethexaddress qQGqkA16ZY6bCYy7Qjr77eU4BPsdadibCG
// 49a80104c0d27a9ba29678d07e87a57151107613
func ToHexString(data string) string {
	// reverse
	buf := base58.Decode(data)

	// [version (1 byte)][address (20 bytes)][digest (4 bytes)]
	hexstr := hex.EncodeToString(buf[1:21])
	return hexstr
}

func reverse(numbers []byte) {
	for i, j := 0, len(numbers)-1; i < j; i, j = i+1, j-1 {
		numbers[i], numbers[j] = numbers[j], numbers[i]
	}
}
