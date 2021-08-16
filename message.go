package poc

import "math/big"

type Commit struct {
	Pid 	string	`json:"pid"`
	Proof   string	`json:"proof"`
	K 		uint8	`json:"k"`
	Difficulty 	*big.Int	`json:"difficulty"`
	Number 		uint64		`json:"number"`
	Timestamp 	int64		`json:"timestamp"`
}

type MinerInfo struct {
	Number      	uint64 		`json:"number"`
	Challenge		string 		`json:"challenge"`
	LastBlockTime 	uint64		`json:"lastBlockTime"`
	Difficulty		*big.Int	`json:"difficulty"`
}
