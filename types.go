package poc

import "encoding/hex"

func MustDecode(input string) []byte {
	dec,err := hex.DecodeString(input)
	if err != nil {
		panic(err)
	}
	return dec
}
