package difficulty

import (
	"github.com/gnc-project/poc"
	"math/big"
	"time"
)

func CalcNextRequiredDifficulty(lastHeaderTime time.Time,difficulty *big.Int, newBlockTime time.Time) *big.Int {

	// Currently we use the retarget formula in ethereum homestead like,
	// which is:
	// diff = parent_diff + parent_diff // 2048 * max(1 - (block_slot - parent_slot) // 10, -199)
	diffUnit := new(big.Int).SetUint64(2048)

	var max = func(x, y int64) int64 {
		if x > y {
			return x
		} else {
			return y
		}
	}

	// Calc ParentDiff
	ParentDiff := difficulty

	// Calc max(1 - (block_slot - parent_slot) // 10, -199)
	MaxOne := new(big.Int).SetInt64(max(1-(newBlockTime.Unix()/poc.PoCSlot-lastHeaderTime.Unix()/int64(poc.PoCSlot))/int64(poc.ToleranceSlot), -199))

	// Calc adjusted part, parent_diff // 2048 * MaxOne
	Adjusted := new(big.Int).Mul(new(big.Int).Div(ParentDiff, diffUnit), MaxOne)

	// Calc Difficulty.
	Diff := new(big.Int).Set(ParentDiff).Add(ParentDiff, Adjusted)

	return Diff
}

