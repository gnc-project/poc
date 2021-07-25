package poc

import (
	"crypto/sha256"
	"math/big"
)

func CalculateDeadline(challenge [32]byte,quality []byte,difficulty *big.Int) *big.Int {
	final := sha256.Sum256(append(challenge[:],quality...))
	f := []byte{final[3],final[2],final[1],final[0]}
	fin := big.NewInt(0).SetBytes(f)
	fb := big.NewInt(0).Div(difficulty,fin)
	return fb
}

func BytesTo32(b []byte) [32]byte {
	var h [32]byte
	copy(h[:],b)
	return h
}