package b58addr

import "testing"

import "github.com/stretchr/testify/assert"

func TestToHexAddressString(t *testing.T) {
	is := assert.New(t)
	hexstr := ToHexString("qQGqkA16ZY6bCYy7Qjr77eU4BPsdadibCG")
	is.Equal(hexstr, "49a80104c0d27a9ba29678d07e87a57151107613")
}
