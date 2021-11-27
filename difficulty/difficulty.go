package difficulty

import (
	"math/big"
	"time"
)

func CalcNextRequiredDifficulty(lastHeaderTime time.Time,difficulty *big.Int, newBlockTime time.Time) *big.Int {

	// Calc ParentDiff
	ParentDiff := difficulty

	// Calc adjusted part, parent_diff // 2048 * MaxOne
	Adjusted := new(big.Int).Mul(ParentDiff,big.NewInt(20))

	if newBlockTime.Unix() - lastHeaderTime.Unix() < 18 {
		return new(big.Int).Add(ParentDiff,Adjusted)
	}
	if newBlockTime.Unix() - lastHeaderTime.Unix() > 18 {
		return new(big.Int).Sub(ParentDiff,Adjusted)
	}

	return ParentDiff
}

