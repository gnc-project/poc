package poc

import "math/big"

type MinerInfo struct {
	Difficulty  *big.Int `json:"difficulty"`
	Number      *big.Int `json:"number"`
	Challenge	[]byte `json:"challenge"`
}
