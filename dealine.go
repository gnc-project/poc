package poc

import (
	"crypto/sha256"
	"math/big"
)

func CalcuteDeadline(pid []byte,challenge []byte,difficulty *big.Int) *big.Int {
	final := sha256.Sum256(append(pid,challenge...))
	f := []byte{final[3],final[2],final[1],final[0]}
	fin := big.NewInt(0).SetBytes(f)
	fb := big.NewInt(0).Div(difficulty,fin)
	return fb
}
