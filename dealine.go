package poc

import (
	"crypto/sha256"
	"math/big"
)

func CalculateDeadline(pid [32]byte,challenge []byte,difficulty *big.Int) *big.Int {
	final := sha256.Sum256(append(pid[:],challenge...))
	f := []byte{final[7],final[6],final[5],final[4],final[3],final[2],final[1],final[0]}
	fin := big.NewInt(0).SetBytes(f)
	fin.Div(fin,big.NewInt(1000))
	fb := big.NewInt(0).Div(difficulty,fin)
	return fb
}

func BytesTo32(b []byte) [32]byte {
	var h [32]byte
	copy(h[:],b)
	return h
}